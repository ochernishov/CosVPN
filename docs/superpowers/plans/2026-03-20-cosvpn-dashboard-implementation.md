# CosVPN Dashboard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Встроить веб-дашборд администратора в бинарник cosvpn-go — управление клиентами, настройками обфускации и мониторинг сервера через браузер.

**Architecture:** Go `net/http` сервер в пакете `admin/`, фронтенд через `embed` (HTML/CSS/JS). Бэкенд вызывает `wg` CLI и читает конфиги напрямую. JWT-авторизация через httpOnly cookie.

**Tech Stack:** Go 1.23+, `net/http`, `embed`, `crypto/hmac`, `os/exec`, `encoding/json`, HTML/CSS/JS (vanilla)

**Spec:** `docs/superpowers/specs/2026-03-20-cosvpn-dashboard-design.md`

**Working directory:** `/Users/cos/CosinnDev/React/СosVPN/CosVPN-Go/`

**Go binary:** `/Users/cos/go-sdk/go/bin/go`

---

## File Map

### New files

| File | Responsibility |
|------|---------------|
| `admin/server.go` | HTTP-сервер, роутинг, static serving, middleware |
| `admin/auth.go` | Логин, JWT генерация/валидация, rate limiting |
| `admin/handlers.go` | API handlers: status, clients, settings, logs |
| `admin/wgctl.go` | Обёртка над wg CLI: show, genkey, конфиги |
| `admin/logger.go` | Ring buffer логов событий (in-memory, 100 записей) |
| `admin/static/login.html` | Страница логина |
| `admin/static/index.html` | SPA: dashboard, clients, settings, logs |
| `admin/static/app.js` | UI логика: fetch API, рендер, модалки |
| `admin/static/style.css` | Тёмная тема, карточки, таблицы |

### Modified files

| File | Changes |
|------|---------|
| `main.go` | Запуск admin.StartServer() если задан COSVPN_ADMIN_PASSWORD |

---

### Task 1: Пакет admin — HTTP-сервер и static serving

**Files:**
- Create: `admin/server.go`
- Create: `admin/static/login.html`
- Create: `admin/static/index.html`
- Create: `admin/static/style.css`

- [ ] **Step 1: Создать admin/server.go**

```go
package admin

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func StartServer(addr, password, wgConfigDir string) {
	mux := http.NewServeMux()

	// Static files
	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("GET /", http.FileServer(http.FS(staticFS)))

	// API routes (заглушки — реализуем в следующих задачах)
	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	})
	mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"vpn":{"status":"up"}}`))
	})

	log.Printf("CosVPN Dashboard starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Printf("Dashboard server error: %v", err)
	}
}
```

- [ ] **Step 2: Создать admin/static/login.html**

Минимальная страница логина: поле пароля, кнопка, тёмная тема. При успехе — redirect на index.html.

- [ ] **Step 3: Создать admin/static/index.html**

Каркас SPA с 4 табами: Dashboard, Clients, Settings, Logs. Подключает app.js и style.css.

- [ ] **Step 4: Создать admin/static/style.css**

Тёмная тема: фон #0a0a0f, карточки #1a1a2e, акцент cyan #06b6d4. Стили для таблиц, карточек, кнопок, модалок, табов.

- [ ] **Step 5: Проверить компиляцию и static serving**

```bash
export PATH="/Users/cos/go-sdk/go/bin:$PATH"
cd /Users/cos/CosinnDev/React/СosVPN/CosVPN-Go
go build -o cosvpn-go .
```

- [ ] **Step 6: Commit**

```bash
git add admin/
git commit -m "feat(admin): add embedded web server with static file serving"
```

---

### Task 2: Авторизация — JWT + rate limiting

**Files:**
- Create: `admin/auth.go`
- Modify: `admin/server.go`

- [ ] **Step 1: Создать admin/auth.go**

Реализовать:
- `GenerateJWT(password string) (string, error)` — HS256, expires 24h
- `ValidateJWT(tokenString, password string) bool`
- `AuthMiddleware(next http.Handler, password string) http.Handler` — проверка cookie "token"
- `RateLimiter` struct — map[ip]attempts, сброс через 1 минуту, лимит 5

```go
package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Claims struct {
	Exp int64 `json:"exp"`
}

func GenerateJWT(secret string) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claims := Claims{Exp: time.Now().Add(24 * time.Hour).Unix()}
	claimsJSON, _ := json.Marshal(claims)
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)

	unsigned := header + "." + payload
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigned))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsigned + "." + sig, nil
}

func ValidateJWT(token, secret string) bool {
	parts := strings.SplitN(token, ".", 3)
	if len(parts) != 3 {
		return false
	}

	// Verify signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return false
	}

	// Check expiry
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return false
	}
	return time.Now().Unix() < claims.Exp
}

func AuthMiddleware(next http.Handler, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || !ValidateJWT(cookie.Value, password) {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type RateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean old attempts
	valid := make([]time.Time, 0)
	for _, t := range rl.attempts[ip] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.attempts[ip] = valid
		return false
	}

	rl.attempts[ip] = append(valid, now)
	return true
}
```

- [ ] **Step 2: Обновить server.go — подключить auth middleware и login handler**

Добавить реальный login handler (проверка пароля, выдача JWT cookie).
API endpoints обернуть в AuthMiddleware.

- [ ] **Step 3: Обновить login.html — fetch POST /api/login**

- [ ] **Step 4: Проверить компиляцию**

```bash
go build -o cosvpn-go .
```

- [ ] **Step 5: Commit**

```bash
git add admin/
git commit -m "feat(admin): add JWT auth with rate limiting"
```

---

### Task 3: wgctl — обёртка над WireGuard CLI

**Files:**
- Create: `admin/wgctl.go`

- [ ] **Step 1: Создать admin/wgctl.go**

Реализовать:
- `WgCtl` struct — хранит путь к конфигу (/etc/wireguard)
- `Status() (ServerStatus, error)` — парсинг `wg show`, `uptime`, `free -m`, `df -h /`
- `ListClients() ([]Client, error)` — парсинг `wg show wg0 dump` + чтение конфигов из clients/
- `AddClient(name string) (*NewClient, error)` — генерация ключей, добавление пира, создание .conf
- `RemoveClient(name string) error` — удаление пира и папки
- `GetClientConfig(name string) (string, error)` — чтение .conf файла
- `GenerateQR(name string) ([]byte, error)` — вызов qrencode, возврат PNG
- `GetSettings() (Settings, error)` — чтение текущих настроек из конфига
- `UpdateSettings(s Settings) error` — обновление конфига, перезапуск wg

Все команды выполняются через `os/exec.Command()`.

Структуры:
```go
type ServerStatus struct {
	Uptime   string `json:"uptime"`
	CPU      int    `json:"cpu"`
	RAM      int    `json:"ram"`
	Disk     int    `json:"disk"`
	VPNUp    bool   `json:"vpnUp"`
	Port     int    `json:"port"`
	PublicIP string `json:"publicIP"`
	Total    int    `json:"totalClients"`
	Online   int    `json:"onlineClients"`
	ObfsMode string `json:"obfsMode"`
}

type Client struct {
	Name          string `json:"name"`
	IP            string `json:"ip"`
	Online        bool   `json:"online"`
	LastHandshake string `json:"lastHandshake"`
	TransferUp    int64  `json:"transferUp"`
	TransferDown  int64  `json:"transferDown"`
	PublicKey     string `json:"publicKey"`
}

type NewClient struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Config string `json:"config"`
	QRData []byte `json:"qrData"`
}

type Settings struct {
	ObfsMode   string   `json:"obfuscationMode"`
	ObfsKeySet bool     `json:"obfuscationKeySet"`
	Port       int      `json:"listenPort"`
	DNS        []string `json:"dns"`
	MTU        int      `json:"mtu"`
	Subnet     string   `json:"subnet"`
}
```

- [ ] **Step 2: Проверить компиляцию**

```bash
go build -o cosvpn-go .
```

- [ ] **Step 3: Commit**

```bash
git add admin/wgctl.go
git commit -m "feat(admin): add wgctl — WireGuard CLI wrapper for dashboard"
```

---

### Task 4: Event Logger — ring buffer

**Files:**
- Create: `admin/logger.go`

- [ ] **Step 1: Создать admin/logger.go**

```go
package admin

import (
	"sync"
	"time"
)

type LogEntry struct {
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`    // "connect", "disconnect", "error", "settings"
	Client  string    `json:"client"`
	Details string    `json:"details"`
}

type EventLogger struct {
	mu      sync.RWMutex
	entries []LogEntry
	maxSize int
}

func NewEventLogger(maxSize int) *EventLogger {
	return &EventLogger{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

func (l *EventLogger) Add(entryType, client, details string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Time:    time.Now(),
		Type:    entryType,
		Client:  client,
		Details: details,
	}

	if len(l.entries) >= l.maxSize {
		l.entries = l.entries[1:]
	}
	l.entries = append(l.entries, entry)
}

func (l *EventLogger) Get(limit int) []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if limit <= 0 || limit > len(l.entries) {
		limit = len(l.entries)
	}

	// Return most recent first
	result := make([]LogEntry, limit)
	for i := 0; i < limit; i++ {
		result[i] = l.entries[len(l.entries)-1-i]
	}
	return result
}
```

- [ ] **Step 2: Commit**

```bash
git add admin/logger.go
git commit -m "feat(admin): add event logger with ring buffer"
```

---

### Task 5: API handlers — полная реализация

**Files:**
- Modify: `admin/handlers.go` (создать)
- Modify: `admin/server.go`

- [ ] **Step 1: Создать admin/handlers.go**

Реализовать все API handlers:
- `HandleStatus(wg *WgCtl) http.HandlerFunc` — GET /api/status
- `HandleListClients(wg *WgCtl) http.HandlerFunc` — GET /api/clients
- `HandleAddClient(wg *WgCtl, log *EventLogger) http.HandlerFunc` — POST /api/clients
- `HandleDeleteClient(wg *WgCtl, log *EventLogger) http.HandlerFunc` — DELETE /api/clients/{name}
- `HandleClientQR(wg *WgCtl) http.HandlerFunc` — GET /api/clients/{name}/qr
- `HandleClientConf(wg *WgCtl) http.HandlerFunc` — GET /api/clients/{name}/conf
- `HandleGetSettings(wg *WgCtl) http.HandlerFunc` — GET /api/settings
- `HandleUpdateSettings(wg *WgCtl, log *EventLogger) http.HandlerFunc` — PUT /api/settings
- `HandleLogs(log *EventLogger) http.HandlerFunc` — GET /api/logs

Каждый handler: парсит запрос → вызывает wgctl → возвращает JSON.

- [ ] **Step 2: Обновить server.go — подключить все routes**

```go
func StartServer(addr, password, wgConfigDir string) {
	wg := NewWgCtl(wgConfigDir)
	eventLog := NewEventLogger(100)
	rl := NewRateLimiter(5, time.Minute)

	mux := http.NewServeMux()

	// Static
	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("GET /", http.FileServer(http.FS(staticFS)))

	// Public
	mux.HandleFunc("POST /api/login", HandleLogin(password, rl))

	// Protected API
	api := http.NewServeMux()
	api.HandleFunc("GET /api/status", HandleStatus(wg))
	api.HandleFunc("GET /api/clients", HandleListClients(wg))
	api.HandleFunc("POST /api/clients", HandleAddClient(wg, eventLog))
	api.HandleFunc("DELETE /api/clients/{name}", HandleDeleteClient(wg, eventLog))
	api.HandleFunc("GET /api/clients/{name}/qr", HandleClientQR(wg))
	api.HandleFunc("GET /api/clients/{name}/conf", HandleClientConf(wg))
	api.HandleFunc("GET /api/settings", HandleGetSettings(wg))
	api.HandleFunc("PUT /api/settings", HandleUpdateSettings(wg, eventLog))
	api.HandleFunc("GET /api/logs", HandleLogs(eventLog))

	mux.Handle("/api/", AuthMiddleware(api, password))
	// Re-register login outside auth
	mux.HandleFunc("POST /api/login", HandleLogin(password, rl))

	log.Printf("CosVPN Dashboard on %s", addr)
	http.ListenAndServe(addr, mux)
}
```

- [ ] **Step 3: Проверить компиляцию**

```bash
go build -o cosvpn-go .
```

- [ ] **Step 4: Commit**

```bash
git add admin/
git commit -m "feat(admin): add all API handlers — status, clients, settings, logs"
```

---

### Task 6: Фронтенд — полный UI

**Files:**
- Modify: `admin/static/index.html`
- Create: `admin/static/app.js`
- Modify: `admin/static/style.css`
- Modify: `admin/static/login.html`

- [ ] **Step 1: login.html — полная страница логина**

Тёмная тема, центрированная форма, лого CosVPN, поле пароля, кнопка. fetch POST /api/login.

- [ ] **Step 2: index.html — SPA с 4 табами**

Header: "CosVPN" + 4 таба (Dashboard, Clients, Settings, Logs) + кнопка Logout.
4 секции (div), переключаются через JS.

- [ ] **Step 3: app.js — логика UI**

Реализовать:
- Tab switching
- Dashboard: fetch /api/status → рендер карточек
- Clients: fetch /api/clients → таблица, кнопки действий
- Add Client: модалка с формой → POST /api/clients → показать QR
- Delete Client: confirm → DELETE /api/clients/{name}
- QR: модалка с img src=/api/clients/{name}/qr
- Download conf: window.open(/api/clients/{name}/conf)
- Settings: fetch /api/settings → форма → PUT /api/settings
- Logs: fetch /api/logs → таблица, auto-refresh 10s
- Error handling: показ ошибок
- Logout: очистить cookie, redirect на login

- [ ] **Step 4: style.css — полные стили**

Тёмная тема:
- Фон: #0a0a0f, карточки: #1a1a2e, бордеры: #2a2a4e
- Акцент: #06b6d4 (cyan)
- Текст: #e2e8f0, secondary: #94a3b8
- Таблицы, кнопки, модалки, формы, табы, статус-индикаторы
- Адаптив для мобильных
- Моноширинный шрифт для IP, ключей

- [ ] **Step 5: Проверить компиляцию**

```bash
go build -o cosvpn-go .
```

- [ ] **Step 6: Commit**

```bash
git add admin/static/
git commit -m "feat(admin): add complete dashboard UI — dark theme, all pages"
```

---

### Task 7: Интеграция в main.go и деплой

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Добавить запуск dashboard в main.go**

```go
import "github.com/ochernishov/cosvpn/admin"

// В func main(), после создания device:
adminPassword := os.Getenv("COSVPN_ADMIN_PASSWORD")
adminPort := os.Getenv("COSVPN_ADMIN_PORT")
if adminPort == "" {
	adminPort = "8443"
}
if adminPassword != "" {
	go admin.StartServer(":"+adminPort, adminPassword, "/etc/wireguard")
	logger.Verbosef("CosVPN Dashboard started on :%s", adminPort)
}
```

- [ ] **Step 2: Собрать для Linux и задеплоить**

```bash
export PATH="/Users/cos/go-sdk/go/bin:$PATH"
GOOS=linux GOARCH=amd64 go build -o cosvpn-go-linux -ldflags="-s -w" .
```

Деплой на сервер:
```bash
export SSHPASS='dypBSvSWR9?WSm'
cat cosvpn-go-linux | sshpass -e ssh root@72.56.26.254 'cat > /usr/local/bin/cosvpn-go && chmod +x /usr/local/bin/cosvpn-go'
```

- [ ] **Step 3: Настроить env на сервере**

```bash
sshpass -e ssh root@72.56.26.254 '
echo "COSVPN_ADMIN_PASSWORD=<сгенерировать>" >> /etc/environment
echo "COSVPN_ADMIN_PORT=8443" >> /etc/environment
'
```

- [ ] **Step 4: Открыть порт и перезапустить**

```bash
sshpass -e ssh root@72.56.26.254 '
ufw allow 8443/tcp 2>/dev/null
systemctl restart wg-quick@wg0
'
```

- [ ] **Step 5: Проверить доступ**

Открыть в браузере: `http://72.56.26.254:8443`

- [ ] **Step 6: Commit**

```bash
git add main.go
git commit -m "feat: integrate dashboard into main.go, deploy to server"
```
