package admin

import (
	"sync"
	"time"
)

type LogEntry struct {
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
	Client  string    `json:"client"`
	Details string    `json:"details"`
}

type EventLogger struct {
	mu      sync.RWMutex
	entries []LogEntry
	maxSize int
}

func NewEventLogger(maxSize int) *EventLogger {
	return &EventLogger{entries: make([]LogEntry, 0, maxSize), maxSize: maxSize}
}

func (l *EventLogger) Add(entryType, client, details string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry := LogEntry{Time: time.Now(), Type: entryType, Client: client, Details: details}
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
	result := make([]LogEntry, limit)
	for i := 0; i < limit; i++ {
		result[i] = l.entries[len(l.entries)-1-i]
	}
	return result
}
