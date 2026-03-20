# CosVPN Dashboard — Веб-панель управления VPN-сервером

> Дата: 2026-03-20
> Статус: Утверждён
> Автор: CosinnDev

---

## 1. Цель

Встроенный веб-дашборд для управления CosVPN-сервером. Один бинарник `cosvpn-go` = VPN-протокол + веб-панель. Доступ по HTTPS через браузер. Только для администратора.

## 2. Архитектура

```
[Браузер] → HTTPS :8443 → [cosvpn-go: встроенный веб-сервер]
                                    ↓
                           [wg show, конфиги, ключи — локально]
                                    ↓
                           [cosvpn-go: VPN-туннель :443]
```

- Веб-сервер встроен в `cosvpn-go` через Go `net/http`
- Фронтенд — чистый HTML/CSS/JS, встроен через Go `embed`
- Бэкенд обращается к VPN напрямую (exec `wg`, чтение/запись конфигов, генерация ключей)
- Один бинарник, ноль внешних зависимостей

## 3. Страницы и функционал

### 3.1. Логин

- Поле пароля + кнопка "Войти"
- Пароль из env `COSVPN_ADMIN_PASSWORD`
- При успехе — JWT в httpOnly cookie (expires 24h)
- Rate limiting: 5 попыток / минута

### 3.2. Dashboard (главная после логина)

Карточки со статистикой:
- **Сервер:** аптайм, CPU, RAM, диск
- **VPN:** статус (up/down), порт, публичный IP
- **Клиенты:** всего / онлайн сейчас
- **Обфускация:** текущий режим (auto/direct/tls)
- **Трафик:** общий ↑↓ за сегодня

### 3.3. Клиенты

Таблица со всеми клиентами:

| Колонка | Данные |
|---------|--------|
| Имя | имя клиента |
| IP | 10.0.0.X |
| Статус | online (зелёный) / offline (серый) |
| Последний хендшейк | время |
| Трафик ↑ | отправлено |
| Трафик ↓ | получено |
| Действия | QR / Скачать .conf / Удалить |

Кнопка "Добавить клиента":
- Модальное окно: поле "Имя клиента"
- При создании: генерация ключей, добавление пира, создание .conf
- Показать QR-код + кнопка "Скачать конфиг"

### 3.4. Настройки

Форма с параметрами:
- **ObfuscationKey** — показать текущий (маскированный), кнопка "Перегенерировать"
- **ObfuscationMode** — select: auto / direct / tls
- **ListenPort** — текущий порт VPN (443)
- **DNS** — DNS-серверы для клиентов (1.1.1.1, 8.8.8.8)
- **MTU** — текущий MTU (1420)
- **Subnet** — подсеть VPN (10.0.0.0/24)
- Кнопка "Применить" → обновляет конфиг, перезапускает VPN

### 3.5. Логи

- Таблица последних 100 событий
- Типы: подключение, отключение, ошибка, смена настроек
- Колонки: время, тип, клиент, детали
- Автообновление каждые 10 секунд

## 4. Безопасность

- **HTTPS** с самоподписанным сертификатом (переиспользуем TLS-сертификат обфускации)
- **Пароль** из переменной окружения `COSVPN_ADMIN_PASSWORD`
- **JWT** токен: HS256, secret из `COSVPN_ADMIN_PASSWORD`, expires 24h
- **httpOnly cookie** — JavaScript не имеет доступа к токену
- **Rate limiting** на `/api/login`: 5 попыток / минута, блокировка на 5 минут
- **Порт 8443** — отдельный от VPN (443), можно ограничить по IP через iptables

## 5. API

### Авторизация

```
POST /api/login
Body: {"password": "..."}
Response: Set-Cookie: token=<jwt>; HttpOnly; Secure; SameSite=Strict
```

Все остальные endpoints требуют валидный JWT cookie.

### Статус

```
GET /api/status
Response: {
  "server": {"uptime": "5d 3h", "cpu": 12, "ram": 45, "disk": 18},
  "vpn": {"status": "up", "port": 443, "publicIP": "72.56.26.254", "interface": "wg0"},
  "clients": {"total": 5, "online": 2},
  "obfuscation": {"mode": "auto", "keySet": true},
  "traffic": {"up": "1.2 GB", "down": "5.4 GB"}
}
```

### Клиенты

```
GET /api/clients
Response: [
  {
    "name": "test-client1",
    "ip": "10.0.0.2",
    "online": false,
    "lastHandshake": "2026-03-20T10:30:00Z",
    "transferUp": 1732,
    "transferDown": 9196,
    "publicKey": "r7aM7..."
  }
]

POST /api/clients
Body: {"name": "my-phone"}
Response: {
  "name": "my-phone",
  "ip": "10.0.0.6",
  "config": "[Interface]\nPrivateKey=...\n...",
  "qrDataURL": "data:image/png;base64,..."
}

DELETE /api/clients/:name
Response: {"ok": true}

GET /api/clients/:name/qr
Response: image/png (QR-код)

GET /api/clients/:name/conf
Response: application/octet-stream (файл .conf)
```

### Настройки

```
GET /api/settings
Response: {
  "obfuscationMode": "auto",
  "obfuscationKeySet": true,
  "listenPort": 443,
  "dns": ["1.1.1.1", "8.8.8.8"],
  "mtu": 1420,
  "subnet": "10.0.0.0/24"
}

PUT /api/settings
Body: {"obfuscationMode": "tls", "dns": ["1.1.1.1"], "mtu": 1420}
Response: {"ok": true, "restarted": true}
```

### Логи

```
GET /api/logs?limit=100
Response: [
  {"time": "2026-03-20T10:30:00Z", "type": "connect", "client": "test-client1", "details": "handshake completed"},
  {"time": "2026-03-20T10:25:00Z", "type": "settings", "client": "", "details": "obfuscation mode changed to tls"}
]
```

## 6. Файловая структура

```
CosVPN-Go/
├── admin/
│   ├── server.go         — HTTP-сервер, роутинг, CORS, static file serving
│   ├── handlers.go       — обработчики API endpoints
│   ├── auth.go           — логин, JWT генерация/валидация, rate limiting
│   ├── wgctl.go          — обёртка: wg show, генерация ключей, управление конфигом
│   ├── logger.go         — in-memory лог событий (ring buffer на 100 записей)
│   └── static/           — go:embed
│       ├── index.html    — SPA: dashboard + clients + settings + logs
│       ├── login.html    — страница логина
│       ├── app.js        — логика UI (fetch API, рендер таблиц, модалки)
│       └── style.css     — стили (тёмная тема, карточки, таблицы)
```

## 7. Интеграция в main.go

```go
// В main.go или при инициализации Device:
adminPassword := os.Getenv("COSVPN_ADMIN_PASSWORD")
if adminPassword != "" {
    go admin.StartServer(":8443", adminPassword, "/etc/wireguard")
}
```

Если `COSVPN_ADMIN_PASSWORD` не задан — дашборд не запускается. Безопасно по умолчанию.

## 8. Визуальный стиль

- Тёмная тема (фон #0a0a0f, карточки #1a1a2e)
- Акцентный цвет: cyan (#06b6d4) — как в иконках CosVPN
- Моноширинный шрифт для IP, ключей, логов
- Минималистичный UI без лишних украшений
- Адаптив: работает на мобильном (для экстренного доступа)

## 9. Деплой

1. Пересобрать `cosvpn-go` с новым пакетом `admin/`
2. Задать `COSVPN_ADMIN_PASSWORD` на сервере
3. Открыть порт 8443 (или выбранный)
4. Зайти на `https://72.56.26.254:8443`
