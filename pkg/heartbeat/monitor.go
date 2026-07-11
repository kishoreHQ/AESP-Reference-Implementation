package heartbeat

import (
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Monitor tracks last beat per session/runtime (ADT-10).
type Monitor struct {
	mu       sync.Mutex
	beats    map[string]time.Time
	timeout  time.Duration
}

func New(timeout time.Duration) *Monitor {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Monitor{beats: map[string]time.Time{}, timeout: timeout}
}

func (m *Monitor) Beat(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.beats[id] = time.Now().UTC()
}

// Status returns "ok" or "error" if last beat older than timeout (or never).
func (m *Monitor) Status(id string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.beats[id]
	if !ok {
		return "unknown"
	}
	if time.Since(t) > m.timeout {
		return "error"
	}
	return "ok"
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module: "pkg/heartbeat", AESPSpecs: []string{"AESP-0011", "AESP-0012"},
		Invariants: []string{"INV-10"}, Status: "stubbed",
	}
}
