# CosVPN Protocol Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Превратить CosVPN-Go в собственный VPN-протокол с обфускацией, не детектируемый DPI.

**Architecture:** Модифицируем wireguard-go: (1) ребрендинг всех идентификаторов, (2) новый пакет `obfs/` для XOR-обфускации заголовков + padding + junk injection, (3) новый пакет `obfs/tls.go` для TLS-обёртки, (4) `transport/` для автовыбора режима. Один код на сервере и во всех клиентах.

**Tech Stack:** Go 1.23+, `crypto/tls`, `crypto/rand`, `encoding/binary`, `golang.org/x/crypto`

**Spec:** `docs/superpowers/specs/2026-03-20-cosvpn-protocol-obfuscation-design.md`

**Working directory:** `/Users/cos/CosinnDev/React/СosVPN/CosVPN-Go/`

---

## File Map

### New files
| File | Responsibility |
|------|---------------|
| `obfs/obfs.go` | XOR header obfuscation, padding, junk detection |
| `obfs/obfs_test.go` | Tests for obfuscation layer |
| `obfs/tls.go` | TLS 1.3 transport wrapper (client + server) |
| `obfs/tls_test.go` | Tests for TLS wrapper |
| `obfs/config.go` | ObfsConfig struct, parsing |
| `transport/transport.go` | Transport interface, auto/direct/tls selection |
| `transport/transport_test.go` | Tests for transport selection |

### Modified files
| File | Changes |
|------|---------|
| `go.mod` | Module path: `golang.zx2c4.com/wireguard` → `github.com/ochernishov/cosvpn` |
| `main.go` | Rename binary, env vars, log messages |
| `device/noise-protocol.go:49-54` | WGIdentifier → CosVPNIdentifier, labels |
| `device/noise-protocol.go:56-77` | Message sizes (add max padding) |
| `device/send.go` | Call `obfs.Obfuscate()` before UDP send |
| `device/receive.go` | Call `obfs.Deobfuscate()` on receive, drop junk |
| `device/device.go` | Add ObfsConfig to Device struct |
| `device/uapi.go` | Parse `obfuscation_key`, `obfuscation_mode` |
| All `*.go` files | Copyright headers, WireGuard → CosVPN in strings |

---

### Task 1: Ребрендинг — идентификаторы протокола

**Files:**
- Modify: `device/noise-protocol.go:49-54`
- Modify: `main.go`
- Modify: `go.mod`

- [ ] **Step 1: Изменить протокольные константы**

В файле `device/noise-protocol.go` строки 49-54:
```go
// Было:
const (
	NoiseConstruction = "Noise_IKpsk2_25519_ChaChaPoly_BLAKE2s"
	WGIdentifier      = "WireGuard v1 zx2c4 Jason@zx2c4.com"
	WGLabelMAC1       = "mac1----"
	WGLabelCookie     = "cookie--"
)

// Стало:
const (
	NoiseConstruction   = "Noise_IKpsk2_25519_ChaChaPoly_BLAKE2s"
	CosVPNIdentifier    = "CosVPN v1 CosinnDev"
	CosVPNLabelMAC1     = "cvpn-mac"
	CosVPNLabelCookie   = "cvpn-cok"
)
```

- [ ] **Step 2: Обновить все ссылки на старые константы**

Найти и заменить во всех `.go` файлах:
- `WGIdentifier` → `CosVPNIdentifier`
- `WGLabelMAC1` → `CosVPNLabelMAC1`
- `WGLabelCookie` → `CosVPNLabelCookie`

Run: `grep -rn "WGIdentifier\|WGLabelMAC1\|WGLabelCookie" device/`

- [ ] **Step 3: Обновить go.mod**

```go
// Было:
module golang.zx2c4.com/wireguard

// Стало:
module github.com/ochernishov/cosvpn
```

Обновить все import paths во всех файлах:
```bash
find . -name "*.go" -exec sed -i '' 's|golang.zx2c4.com/wireguard|github.com/ochernishov/cosvpn|g' {} +
```

- [ ] **Step 4: Обновить main.go**

Заменить:
- `ENV_WG_TUN_FD` → `ENV_COSVPN_TUN_FD`
- `ENV_WG_UAPI_FD` → `ENV_COSVPN_UAPI_FD`
- `ENV_WG_PROCESS_FOREGROUND` → `ENV_COSVPN_PROCESS_FOREGROUND`
- Все строки "wireguard" в логах → "cosvpn"
- Usage: `wireguard-go` → `cosvpn-go`

- [ ] **Step 5: Обновить copyright headers**

Во всех `.go` файлах заменить:
```
Copyright (C) 2017-2025 WireGuard LLC. All Rights Reserved.
```
на:
```
Copyright (C) 2026 CosinnDev. Based on WireGuard by Jason A. Donenfeld.
```

- [ ] **Step 6: Проверить компиляцию**

Run: `cd CosVPN-Go && go build -o cosvpn-go .`
Expected: бинарник `cosvpn-go` собрался без ошибок

- [ ] **Step 7: Запустить существующие тесты**

Run: `go test ./...`
Expected: все тесты проходят (тесты noise/cookie могут упасть из-за смены констант — это ожидаемо, поправить)

- [ ] **Step 8: Commit**

```bash
git add -A CosVPN-Go/
git commit -m "feat(cosvpn-go): rebrand protocol identifiers — WireGuard → CosVPN"
```

---

### Task 2: Пакет obfs — ядро обфускации

**Files:**
- Create: `obfs/obfs.go`
- Create: `obfs/obfs_test.go`
- Create: `obfs/config.go`

- [ ] **Step 1: Создать obfs/config.go**

```go
package obfs

// ObfsConfig хранит параметры обфускации
type ObfsConfig struct {
	Key  [16]byte // 128-бит ключ для XOR заголовков
	Mode string   // "auto", "direct", "tls"
}

// DefaultConfig возвращает конфиг с отключённой обфускацией
func DefaultConfig() ObfsConfig {
	return ObfsConfig{Mode: "direct"}
}

// IsEnabled возвращает true если обфускация включена (ключ не нулевой)
func (c *ObfsConfig) IsEnabled() bool {
	var zero [16]byte
	return c.Key != zero
}
```

- [ ] **Step 2: Написать тест для obfs/obfs.go**

Создать `obfs/obfs_test.go`:
```go
package obfs

import (
	"crypto/rand"
	"encoding/binary"
	"testing"
)

func TestObfuscateDeobfuscateRoundtrip(t *testing.T) {
	var key [16]byte
	rand.Read(key[:])

	// Создаём пакет, имитирующий WG Initiation (type=1, 148 байт)
	packet := make([]byte, 148)
	binary.LittleEndian.PutUint32(packet[0:4], 1)
	rand.Read(packet[4:])
	original := make([]byte, len(packet))
	copy(original, packet)

	// Obfuscate
	obfuscated, err := Obfuscate(packet, key)
	if err != nil {
		t.Fatal(err)
	}

	// Размер должен отличаться (padding добавлен)
	if len(obfuscated) == len(original) {
		t.Error("obfuscated packet should have different size due to padding")
	}

	// Первые 4 байта не должны быть 01 00 00 00
	msgType := binary.LittleEndian.Uint32(obfuscated[0:4])
	if msgType == 1 {
		t.Error("message type should be obfuscated")
	}

	// Deobfuscate
	restored, err := Deobfuscate(obfuscated, key)
	if err != nil {
		t.Fatal(err)
	}

	// Должен совпасть с оригиналом
	if len(restored) != len(original) {
		t.Errorf("restored size %d != original %d", len(restored), len(original))
	}
	for i := range original {
		if restored[i] != original[i] {
			t.Errorf("byte %d differs: got %x, want %x", i, restored[i], original[i])
			break
		}
	}
}

func TestJunkPacketDetection(t *testing.T) {
	var key [16]byte
	rand.Read(key[:])

	junk := MakeJunkPacket(key)
	if !IsJunkPacket(junk, key) {
		t.Error("should detect junk packet")
	}

	// Обычный пакет не должен быть junk
	normal := make([]byte, 148)
	binary.LittleEndian.PutUint32(normal[0:4], 1)
	if IsJunkPacket(normal, key) {
		t.Error("normal packet should not be detected as junk")
	}
}

func TestObfuscateDisabledWithZeroKey(t *testing.T) {
	var key [16]byte // нулевой ключ = обфускация выключена

	packet := make([]byte, 100)
	rand.Read(packet[:])
	original := make([]byte, len(packet))
	copy(original, packet)

	result, err := Obfuscate(packet, key)
	if err != nil {
		t.Fatal(err)
	}

	// С нулевым ключом пакет не должен меняться
	if len(result) != len(original) {
		t.Error("with zero key, packet should pass through unchanged")
	}
}
```

- [ ] **Step 3: Запустить тест, убедиться что FAIL**

Run: `cd CosVPN-Go && go test ./obfs/ -v`
Expected: FAIL — `Obfuscate`, `Deobfuscate`, `MakeJunkPacket`, `IsJunkPacket` не определены

- [ ] **Step 4: Реализовать obfs/obfs.go**

```go
package obfs

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	mrand "math/rand"
)

const (
	headerSize     = 16          // обфусцируемая часть заголовка
	maxPadding     = 64          // максимальный padding
	junkMarker     = byte(0xFF)  // маркер junk-пакета (до XOR)
	minJunkSize    = 64
	maxJunkSize    = 1264
)

// Obfuscate применяет обфускацию к пакету:
// 1. XOR первых 16 байт заголовка с ключом
// 2. Добавляет случайный padding (0-64 байт)
// 3. Последний байт = длина padding
// С нулевым ключом возвращает пакет без изменений.
func Obfuscate(packet []byte, key [16]byte) ([]byte, error) {
	if !isKeySet(key) {
		return packet, nil
	}
	if len(packet) < 4 {
		return nil, errors.New("packet too short")
	}

	// XOR header
	xorHeader(packet, key)

	// Random padding
	padLen := mrand.Intn(maxPadding + 1)
	padding := make([]byte, padLen+1) // +1 для байта длины
	rand.Read(padding[:padLen])
	padding[padLen] = byte(padLen) // последний байт = длина padding

	result := make([]byte, len(packet)+len(padding))
	copy(result, packet)
	copy(result[len(packet):], padding)

	return result, nil
}

// Deobfuscate снимает обфускацию:
// 1. Читает последний байт = длина padding, отрезает
// 2. XOR первых 16 байт для восстановления заголовка
func Deobfuscate(packet []byte, key [16]byte) ([]byte, error) {
	if !isKeySet(key) {
		return packet, nil
	}
	if len(packet) < 5 {
		return nil, errors.New("packet too short for deobfuscation")
	}

	// Читаем padding length из последнего байта
	padLen := int(packet[len(packet)-1])
	if padLen+1 > len(packet)-4 {
		return nil, errors.New("invalid padding length")
	}

	// Отрезаем padding + байт длины
	restored := packet[:len(packet)-padLen-1]

	// XOR header (обратная операция)
	xorHeader(restored, key)

	return restored, nil
}

// MakeJunkPacket создаёт junk-пакет для сбивания DPI-паттернов
func MakeJunkPacket(key [16]byte) []byte {
	size := minJunkSize + mrand.Intn(maxJunkSize-minJunkSize)
	junk := make([]byte, size)
	rand.Read(junk)
	// Маркер: первый байт после XOR = 0xFF
	junk[0] = junkMarker ^ key[0]
	return junk
}

// IsJunkPacket проверяет, является ли пакет junk
func IsJunkPacket(packet []byte, key [16]byte) bool {
	if len(packet) < 1 {
		return false
	}
	return (packet[0] ^ key[0]) == junkMarker
}

// ShouldSendJunk возвращает true с вероятностью ~30%
func ShouldSendJunk() bool {
	return mrand.Intn(100) < 30
}

func xorHeader(packet []byte, key [16]byte) {
	n := headerSize
	if len(packet) < n {
		n = len(packet)
	}
	for i := 0; i < n; i++ {
		packet[i] ^= key[i]
	}
}

func isKeySet(key [16]byte) bool {
	var zero [16]byte
	return key != zero
}
```

- [ ] **Step 5: Запустить тесты**

Run: `go test ./obfs/ -v`
Expected: PASS — все 3 теста проходят

- [ ] **Step 6: Commit**

```bash
git add obfs/
git commit -m "feat(obfs): add packet obfuscation — XOR header, padding, junk injection"
```

---

### Task 3: Интеграция обфускации в send/receive

**Files:**
- Modify: `device/device.go`
- Modify: `device/send.go`
- Modify: `device/receive.go`
- Modify: `device/uapi.go`

- [ ] **Step 1: Добавить ObfsConfig в Device struct**

В `device/device.go` добавить поле в struct `Device`:
```go
import "github.com/ochernishov/cosvpn/obfs"

type Device struct {
	// ... существующие поля ...
	obfsConfig obfs.ObfsConfig
}
```

Добавить методы:
```go
func (device *Device) SetObfsConfig(config obfs.ObfsConfig) {
	device.obfsConfig = config
}

func (device *Device) GetObfsConfig() obfs.ObfsConfig {
	return device.obfsConfig
}
```

- [ ] **Step 2: Модифицировать send.go — обфускация перед отправкой**

В файле `device/send.go`, в функции `peer.SendBuffers()` (или в месте финальной отправки через `device.net.bind.Send()`), перед вызовом `Send`:

```go
// Junk injection
if obfs.ShouldSendJunk() && device.obfsConfig.IsEnabled() {
	junk := obfs.MakeJunkPacket(device.obfsConfig.Key)
	device.net.bind.Send([][]byte{junk}, peer.endpoint)
}

// Obfuscate each buffer
if device.obfsConfig.IsEnabled() {
	for i, buf := range bufs {
		obfuscated, err := obfs.Obfuscate(buf, device.obfsConfig.Key)
		if err == nil {
			bufs[i] = obfuscated
		}
	}
}
```

- [ ] **Step 3: Модифицировать receive.go — деобфускация при получении**

В `device/receive.go`, в `RoutineReceiveIncoming()`, после получения пакета и перед определением типа сообщения:

```go
// Deobfuscate
if device.obfsConfig.IsEnabled() {
	// Drop junk packets
	if obfs.IsJunkPacket(packet, device.obfsConfig.Key) {
		continue
	}
	deobfuscated, err := obfs.Deobfuscate(packet, device.obfsConfig.Key)
	if err != nil {
		continue
	}
	packet = deobfuscated
	size = len(packet)
}
```

- [ ] **Step 4: Добавить параметры в UAPI (uapi.go)**

В `device/uapi.go`, в функции `IpcSetOperation()`, добавить обработку новых ключей:

```go
case "obfuscation_key":
	keyBytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil || len(keyBytes) != 16 {
		errMsg = "invalid obfuscation_key"
	} else {
		var key [16]byte
		copy(key[:], keyBytes)
		device.obfsConfig.Key = key
	}
case "obfuscation_mode":
	if value != "auto" && value != "direct" && value != "tls" {
		errMsg = "invalid obfuscation_mode"
	} else {
		device.obfsConfig.Mode = value
	}
```

В `IpcGetOperation()` добавить вывод:
```go
sendf("obfuscation_key=%s", base64.StdEncoding.EncodeToString(device.obfsConfig.Key[:]))
sendf("obfuscation_mode=%s", device.obfsConfig.Mode)
```

- [ ] **Step 5: Собрать и проверить**

Run: `go build -o cosvpn-go . && go test ./... 2>&1 | tail -20`
Expected: компиляция и тесты проходят

- [ ] **Step 6: Commit**

```bash
git add device/ obfs/
git commit -m "feat(cosvpn-go): integrate obfuscation into send/receive pipeline"
```

---

### Task 4: TLS-обёртка

**Files:**
- Create: `obfs/tls.go`
- Create: `obfs/tls_test.go`

- [ ] **Step 1: Написать тест для TLS transport**

Создать `obfs/tls_test.go`:
```go
package obfs

import (
	"testing"
	"time"
)

func TestTLSTransportRoundtrip(t *testing.T) {
	// Запускаем TLS-сервер
	server, err := NewTLSListener("127.0.0.1:0") // random port
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	addr := server.Addr()

	// Канал для полученных данных
	received := make(chan []byte, 1)
	go func() {
		buf, _ := server.ReadPacket()
		received <- buf
	}()

	// Подключаемся клиентом
	client, err := NewTLSClient(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Отправляем пакет
	payload := []byte("hello cosvpn")
	err = client.WritePacket(payload)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем получение
	select {
	case data := <-received:
		if string(data) != string(payload) {
			t.Errorf("got %q, want %q", data, payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for packet")
	}
}
```

- [ ] **Step 2: Запустить тест — FAIL**

Run: `go test ./obfs/ -run TestTLS -v`
Expected: FAIL

- [ ] **Step 3: Реализовать obfs/tls.go**

```go
package obfs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
	"time"
)

// TLSListener принимает TLS-соединения и извлекает CosVPN-пакеты
type TLSListener struct {
	listener net.Listener
	conn     net.Conn
	mu       sync.Mutex
}

// TLSClient заворачивает CosVPN-пакеты в TLS
type TLSClient struct {
	conn net.Conn
	mu   sync.Mutex
}

// NewTLSListener создаёт TLS-сервер с самоподписанным сертификатом
func NewTLSListener(addr string) (*TLSListener, error) {
	cert, err := generateSelfSignedCert()
	if err != nil {
		return nil, fmt.Errorf("generate cert: %w", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}

	listener, err := tls.Listen("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &TLSListener{listener: listener}, nil
}

func (l *TLSListener) Addr() string {
	return l.listener.Addr().String()
}

func (l *TLSListener) AcceptAndServe(handler func([]byte)) error {
	conn, err := l.listener.Accept()
	if err != nil {
		return err
	}
	l.mu.Lock()
	l.conn = conn
	l.mu.Unlock()

	for {
		buf, err := readFramed(conn)
		if err != nil {
			return err
		}
		handler(buf)
	}
}

func (l *TLSListener) ReadPacket() ([]byte, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}
	l.mu.Lock()
	l.conn = conn
	l.mu.Unlock()
	return readFramed(conn)
}

func (l *TLSListener) WritePacket(data []byte) error {
	l.mu.Lock()
	conn := l.conn
	l.mu.Unlock()
	if conn == nil {
		return fmt.Errorf("no connection")
	}
	return writeFramed(conn, data)
}

func (l *TLSListener) Close() error {
	return l.listener.Close()
}

// NewTLSClient подключается к CosVPN-серверу через TLS
func NewTLSClient(addr string) (*TLSClient, error) {
	config := &tls.Config{
		InsecureSkipVerify: true, // самоподписанный серт
		MinVersion:         tls.VersionTLS13,
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &TLSClient{conn: conn}, nil
}

func (c *TLSClient) WritePacket(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return writeFramed(c.conn, data)
}

func (c *TLSClient) ReadPacket() ([]byte, error) {
	return readFramed(c.conn)
}

func (c *TLSClient) Close() error {
	return c.conn.Close()
}

// Framing: [2 bytes length BE][payload]
func writeFramed(conn net.Conn, data []byte) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(len(data)))
	if _, err := conn.Write(header); err != nil {
		return err
	}
	_, err := conn.Write(data)
	return err
}

func readFramed(conn net.Conn) ([]byte, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(header)
	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func generateSelfSignedCert() (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  key,
	}, nil
}
```

- [ ] **Step 4: Запустить тесты**

Run: `go test ./obfs/ -v`
Expected: PASS — все тесты (obfuscation + TLS)

- [ ] **Step 5: Commit**

```bash
git add obfs/
git commit -m "feat(obfs): add TLS 1.3 transport wrapper"
```

---

### Task 5: Transport — автовыбор режима

**Files:**
- Create: `transport/transport.go`
- Create: `transport/transport_test.go`

- [ ] **Step 1: Написать тест**

Создать `transport/transport_test.go`:
```go
package transport

import (
	"testing"

	"github.com/ochernishov/cosvpn/obfs"
)

func TestAutoTransportFallback(t *testing.T) {
	config := obfs.ObfsConfig{Mode: "auto"}

	tr := NewAutoTransport(config)
	if tr.CurrentMode() != "auto" {
		t.Errorf("expected auto, got %s", tr.CurrentMode())
	}
}

func TestDirectMode(t *testing.T) {
	config := obfs.ObfsConfig{Mode: "direct"}
	tr := NewAutoTransport(config)
	if tr.CurrentMode() != "direct" {
		t.Errorf("expected direct, got %s", tr.CurrentMode())
	}
}

func TestTLSMode(t *testing.T) {
	config := obfs.ObfsConfig{Mode: "tls"}
	tr := NewAutoTransport(config)
	if tr.CurrentMode() != "tls" {
		t.Errorf("expected tls, got %s", tr.CurrentMode())
	}
}
```

- [ ] **Step 2: Реализовать transport/transport.go**

```go
package transport

import (
	"sync"
	"time"

	"github.com/ochernishov/cosvpn/obfs"
)

const autoTimeout = 3 * time.Second

// AutoTransport выбирает режим подключения
type AutoTransport struct {
	config      obfs.ObfsConfig
	currentMode string
	mu          sync.RWMutex
	tlsClient   *obfs.TLSClient
}

func NewAutoTransport(config obfs.ObfsConfig) *AutoTransport {
	mode := config.Mode
	if mode == "" {
		mode = "auto"
	}
	return &AutoTransport{
		config:      config,
		currentMode: mode,
	}
}

func (t *AutoTransport) CurrentMode() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.currentMode
}

func (t *AutoTransport) SetMode(mode string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.currentMode = mode
}

// NeedsTLS возвращает true если текущий режим требует TLS
func (t *AutoTransport) NeedsTLS() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.currentMode == "tls"
}

// SwitchToTLS переключает на TLS (вызывается при fallback в auto-режиме)
func (t *AutoTransport) SwitchToTLS() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.currentMode = "tls"
}

// AutoTimeout возвращает таймаут для direct-режима в auto
func (t *AutoTransport) AutoTimeout() time.Duration {
	return autoTimeout
}
```

- [ ] **Step 3: Тесты**

Run: `go test ./transport/ -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add transport/
git commit -m "feat(transport): add auto/direct/tls mode selection"
```

---

### Task 6: Сборка серверного бинарника и деплой

**Files:**
- Modify: `main.go` (финальные правки)
- Server: `/etc/wireguard/` → cosvpn конфиг

- [ ] **Step 1: Собрать бинарник для Linux**

```bash
cd CosVPN-Go
GOOS=linux GOARCH=amd64 go build -o cosvpn-go -ldflags="-s -w" .
```
Expected: файл `cosvpn-go` (~5-7 MB)

- [ ] **Step 2: Сгенерировать ObfuscationKey**

```bash
head -c 16 /dev/urandom | base64
```
Сохранить результат — это общий ключ для сервера и клиентов.

- [ ] **Step 3: Загрузить бинарник на сервер**

```bash
scp cosvpn-go root@72.56.26.254:/usr/local/bin/cosvpn-go
ssh root@72.56.26.254 'chmod +x /usr/local/bin/cosvpn-go'
```

- [ ] **Step 4: Обновить серверный конфиг**

На сервере обновить `/etc/wireguard/wg0.conf`:
```ini
[Interface]
PrivateKey = <existing>
ListenPort = 443
ObfuscationKey = <base64-ключ-из-step-2>
ObfuscationMode = auto

[Peer]
# ... existing peers ...
```

- [ ] **Step 5: Создать systemd-сервис**

На сервере создать `/etc/systemd/system/cosvpn.service`:
```ini
[Unit]
Description=CosVPN Tunnel
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/cosvpn-go cvpn0
ExecStop=/usr/bin/kill $MAINPID
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

- [ ] **Step 6: Запустить и проверить**

```bash
systemctl stop wg-quick@wg0
systemctl enable cosvpn
systemctl start cosvpn
systemctl status cosvpn
```

- [ ] **Step 7: Commit**

```bash
git commit -m "feat(cosvpn-go): server binary build and deploy config"
```

---

### Task 7: Тестирование — подключение через обфусцированный протокол

- [ ] **Step 1: Создать клиентский конфиг с обфускацией**

```ini
[Interface]
PrivateKey = <client-private-key>
Address = 10.0.0.X/32
DNS = 1.1.1.1, 8.8.8.8
ObfuscationKey = <тот-же-base64-ключ>
ObfuscationMode = auto

[Peer]
PublicKey = <server-public-key>
Endpoint = 72.56.26.254:443
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
```

- [ ] **Step 2: Собрать клиентский бинарник для macOS**

```bash
cd CosVPN-Go
go build -o cosvpn-go .
```

- [ ] **Step 3: Запустить клиент на маке**

```bash
sudo ./cosvpn-go utun
```
Проверить `ifconfig utun*` — должен появиться интерфейс с 10.0.0.X

- [ ] **Step 4: Проверить подключение**

```bash
ping 10.0.0.1          # ping сервера через VPN
curl ifconfig.me       # должен показать 72.56.26.254
```

- [ ] **Step 5: Проверить обфускацию — tcpdump**

На сервере:
```bash
tcpdump -i eth0 port 443 -c 10 -X | head -50
```
Убедиться что пакеты НЕ содержат WireGuard сигнатуры (01 00 00 00, 02 00 00 00 и т.д.)

- [ ] **Step 6: Commit финальный**

```bash
git add -A
git commit -m "feat(cosvpn): working obfuscated VPN protocol — tested on server"
```

---

### Task 8: Обновить ROADMAP и документацию

- [ ] **Step 1: Обновить ROADMAP.md** — добавить новый этап "CosVPN Protocol" как завершённый
- [ ] **Step 2: Обновить deploy-vpn скилл** — добавить команды для cosvpn-go
- [ ] **Step 3: Commit**

```bash
git add ROADMAP.md
git commit -m "docs: update roadmap with CosVPN Protocol stage"
```
