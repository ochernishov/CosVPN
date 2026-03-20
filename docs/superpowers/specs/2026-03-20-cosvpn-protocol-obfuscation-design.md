# CosVPN Protocol — Обфускация и собственный протокол

> Дата: 2026-03-20
> Статус: Утверждён
> Автор: CosinnDev

---

## 1. Цель

Превратить CosVPN из форка WireGuard в **собственный VPN-продукт** с уникальным протоколом, который:
- Не детектируется DPI (Deep Packet Inspection)
- Не блокируется провайдерами (в т.ч. российскими)
- Выглядит как обычный HTTPS-трафик
- Работает на всех платформах (Android, macOS, iOS, сервер)

## 2. Проблема

Стандартный WireGuard легко блокируется DPI по трём признакам:
1. **Message Type** — первые 4 байта содержат значения 1-4 в Little-Endian
2. **Фиксированные размеры пакетов** — handshake initiation всегда 148 байт, response 92 байта, cookie reply 64 байта
3. **UDP на нестандартном порту** — типичные порты 51820, легко блокируются

## 3. Решения (ключевые)

- **Собственный протокол** — CosVPN Protocol, несовместим со стандартным WireGuard
- **Убрать ВСЕ упоминания WireGuard** из кода, UI, конфигов, логов
- **Два слоя обфускации**: обфускация пакетов + TLS-обёртка
- **Автовыбор режима** с возможностью ручного переключения

## 4. Архитектура

### 4.1. Общая схема

```
[Клиент CosVPN]                              [Сервер CosVPN]
     |                                              |
     |  IP-пакет из TUN                             |
     v                                              |
[CosVPN-Go: шифрование ChaCha20-Poly1305]          |
     |                                              |
     v                                              |
[obfs: обфускация пакета]                           |
  - XOR magic bytes с obfs_key                      |
  - Random padding (размер от 0 до 64 байт)         |
  - Junk packet injection (вероятность 30%)          |
     |                                              |
     v                                              |
[transport: выбор режима]                           |
  - direct: UDP пакет как есть                      |
  - tls: обёртка в TLS 1.3 record                   |
  - auto: попробовать direct, fallback на tls       |
     |                                              |
     v                                              v
[UDP/TLS ────────── Интернет ──────────> receive + deobfs]
     443                                           443
```

### 4.2. Слой 1 — Обфускация пакетов (пакет `obfs/`)

#### 4.2.1. XOR Header Obfuscation

Стандартный WireGuard отправляет Message Type как первые 4 байта (uint32 LE):
- `01 00 00 00` = Handshake Initiation
- `02 00 00 00` = Handshake Response
- `03 00 00 00` = Cookie Reply
- `04 00 00 00` = Transport Data

CosVPN шифрует первые 16 байт заголовка XOR с ключом обфускации:

```go
func ObfuscateHeader(packet []byte, key [16]byte) {
    for i := 0; i < 16 && i < len(packet); i++ {
        packet[i] ^= key[i]
    }
}
```

`obfs_key` — 128-бит ключ, общий между клиентом и сервером. Генерируется один раз при настройке.

#### 4.2.2. Random Padding

После шифрования CosVPN добавляет случайный padding:

```go
type ObfuscatedPacket struct {
    Header    [16]byte     // XOR-обфусцированный заголовок
    Payload   []byte       // зашифрованные данные (ChaCha20-Poly1305)
    Padding   []byte       // случайные байты (0-64)
    PadLength byte         // последний байт = длина padding
}
```

- Padding 0-64 байт (crypto/rand)
- Последний байт пакета = длина padding (чтобы receiver мог отрезать)
- Результат: размеры пакетов всегда разные, DPI не может ловить по фиксированным размерам

#### 4.2.3. Junk Packet Injection

С вероятностью 30% перед реальным пакетом отправляется junk-пакет:

```go
func MaybeSendJunk(conn Bind, endpoint Endpoint, key [16]byte) {
    if rand.Intn(100) < 30 {
        junkSize := 64 + rand.Intn(1200) // 64-1264 байт
        junk := make([]byte, junkSize)
        rand.Read(junk)
        // Первый байт маркер: 0xFF (после XOR с key)
        junk[0] = 0xFF ^ key[0]
        conn.Send([][]byte{junk}, endpoint)
    }
}
```

Receiver узнаёт junk по маркеру `0xFF` (после de-XOR) и отбрасывает.

### 4.3. Слой 2 — TLS-обёртка (пакет `obfs/tls.go`)

Когда режим = `tls` или `auto` (fallback):

1. Клиент устанавливает TLS 1.3 соединение с сервером на порту 443
2. Стандартный TLS handshake (выглядит как обычный HTTPS)
3. Внутри TLS-канала передаются обфусцированные CosVPN-пакеты
4. Снаружи DPI видит обычный TLS-трафик на порту 443

```go
type TLSTransport struct {
    conn     *tls.Conn
    bind     Bind          // оригинальный UDP bind
    config   *tls.Config
}

func (t *TLSTransport) Send(bufs [][]byte, ep Endpoint) error {
    for _, buf := range bufs {
        // Префикс: 2 байта длины пакета
        length := make([]byte, 2)
        binary.BigEndian.PutUint16(length, uint16(len(buf)))
        t.conn.Write(length)
        t.conn.Write(buf)
    }
    return nil
}
```

Сервер слушает на 443:
- Если приходит TLS ClientHello → принимает TLS-соединение, внутри CosVPN
- Если приходит UDP-пакет → обрабатывает как direct-режим (обфусцированный CosVPN)

### 4.4. Слой 3 — Автовыбор режима (пакет `transport/`)

```go
type AutoTransport struct {
    direct    *DirectTransport   // UDP
    tls       *TLSTransport      // TLS
    mode      string             // "auto" / "direct" / "tls"
    current   Transport          // текущий активный
}

func (a *AutoTransport) Connect(endpoint Endpoint) error {
    if a.mode == "tls" {
        return a.tls.Connect(endpoint)
    }
    if a.mode == "direct" {
        return a.direct.Connect(endpoint)
    }
    // auto: пробуем direct 3 секунды
    err := a.direct.ConnectWithTimeout(endpoint, 3*time.Second)
    if err == nil {
        a.current = a.direct
        return nil
    }
    // fallback на TLS
    a.current = a.tls
    return a.tls.Connect(endpoint)
}
```

## 5. Модифицируемые файлы

### Существующие (модификация)

| Файл | Изменения |
|------|-----------|
| `device/noise-protocol.go` | Убрать WireGuard constants (WGIdentifier, WGLabelMAC1), заменить на CosVPN. Вызов obfs перед marshal/unmarshal |
| `device/send.go` | Вызов `obfs.Obfuscate()` перед отправкой, junk injection |
| `device/receive.go` | Вызов `obfs.Deobfuscate()` при получении, отбрасывание junk |
| `conn/bind_std.go` | Подключение TLSTransport как альтернативного bind |
| `device/uapi.go` | Новые параметры: `obfuscation_key`, `obfuscation_mode` |
| `device/device.go` | Хранение obfs config в Device struct |
| `main.go` | Переименование: cosvpn-go, логи CosVPN |

### Новые файлы

| Файл | Назначение |
|------|-----------|
| `obfs/obfs.go` | Ядро обфускации: XOR, padding, junk detection |
| `obfs/tls.go` | TLS-обёртка: клиент и серверная часть |
| `obfs/config.go` | Структура конфигурации обфускации |
| `transport/auto.go` | Автовыбор режима (direct/tls/auto) |
| `transport/transport.go` | Интерфейс Transport |

### Ребрендинг (все файлы)

- Все строки "WireGuard" → "CosVPN"
- `WGIdentifier = "WireGuard v1 zx2c4 Jason@zx2c4.com"` → `CosVPNIdentifier = "CosVPN v1 CosinnDev"`
- `WGLabelMAC1 = "mac1----"` → `CosVPNLabelMAC1 = "cvpn-mac"` (8 байт)
- `WGLabelCookie = "cookie--"` → `CosVPNLabelCookie = "cvpn-cok"` (8 байт)
- Бинарник: `wireguard-go` → `cosvpn-go`
- Логи, ошибки, комментарии — всё CosVPN

## 6. Формат конфигурации CosVPN

```ini
[Interface]
PrivateKey = <base64>
ListenPort = 443
ObfuscationKey = <base64-128bit>
ObfuscationMode = auto

[Peer]
PublicKey = <base64>
PresharedKey = <base64>
AllowedIPs = 10.0.0.0/24
Endpoint = 72.56.26.254:443
PersistentKeepalive = 25
```

Новые параметры:
- `ObfuscationKey` — 128-бит ключ для XOR-обфускации заголовков (base64)
- `ObfuscationMode` — режим: `auto` (по умолчанию), `direct`, `tls`

## 7. Серверный деплой

На существующем сервере 72.56.26.254:
1. Компилируем `cosvpn-go` (`GOOS=linux GOARCH=amd64 go build -o cosvpn-go`)
2. Копируем бинарник на сервер
3. Заменяем wg-quick на cosvpn-quick (обёртка)
4. Генерируем ObfuscationKey: `head -c 16 /dev/urandom | base64`
5. Обновляем конфиг, перезапускаем

## 8. Клиентская интеграция

CosVPN-Go компилируется как библиотека:
- **Android**: `gomobile bind` → `.aar` файл, импорт в CosVPN-Android
- **macOS/iOS**: `gomobile bind` → `.xcframework`, импорт в CosVPN-Apple

Обфускация встроена в библиотеку — клиентские приложения просто передают конфиг с `ObfuscationKey` и `ObfuscationMode`.

## 9. Безопасность

- **Криптография не меняется** — ChaCha20-Poly1305, Curve25519, BLAKE2s остаются
- **Noise IKpsk2 handshake** — остаётся, но с изменёнными идентификаторами
- **XOR-обфускация** — не для безопасности, а для обхода DPI. Безопасность обеспечивается нижележащим шифрованием
- **TLS 1.3** — дополнительный слой шифрования в tls-режиме
- **ObfuscationKey** — должен храниться как секрет, наравне с PrivateKey

## 10. Этапы реализации

1. **Ребрендинг CosVPN-Go** — убрать все WireGuard, переименовать
2. **Пакет obfs/** — XOR, padding, junk
3. **Модификация send.go/receive.go** — вызов обфускации
4. **Пакет transport/** — автовыбор режима
5. **TLS-обёртка** — obfs/tls.go
6. **Серверный бинарник** — сборка, деплой на 72.56.26.254
7. **Тестирование** — подключение через обфусцированный протокол
8. **Интеграция в клиенты** — пересборка CosVPN-Android, CosVPN-Apple
