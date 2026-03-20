package transport

import (
	"sync"
	"time"

	"github.com/ochernishov/cosvpn/obfs"
)

const AutoTimeout = 3 * time.Second

// AutoTransport выбирает режим подключения
type AutoTransport struct {
	config      obfs.ObfsConfig
	currentMode string
	mu          sync.RWMutex
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

func (t *AutoTransport) NeedsTLS() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.currentMode == "tls"
}

func (t *AutoTransport) SwitchToTLS() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.currentMode = "tls"
}

func (t *AutoTransport) GetAutoTimeout() time.Duration {
	return AutoTimeout
}
