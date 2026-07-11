package eventbus

import (
	"context"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Event struct {
	Type       string
	ID         types.EventID
	WorkUnitID types.WorkUnitID
	SessionID  types.SessionID
	TraceID    types.TraceID
	Time       time.Time
	Data       map[string]any
	Seq        int64 // monotonic per bus (UI-RT-01)
}

type Bus interface {
	Publish(ctx context.Context, e Event) error
	Subscribe(ctx context.Context, workUnitFilter string) (<-chan Event, error)
	Replay(ctx context.Context, workUnitID types.WorkUnitID) ([]Event, error)
	// Since returns all events with Seq > since (global journal for WS reconnect).
	Since(ctx context.Context, since int64) ([]Event, error)
	Seq() int64
}

type MemoryBus struct {
	mu   sync.Mutex
	log  []Event
	subs map[string][]chan Event
	seq  int64
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{subs: make(map[string][]chan Event)}
}

func (b *MemoryBus) Publish(ctx context.Context, e Event) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.seq++
	e.Seq = b.seq
	if e.ID == "" {
		e.ID = types.EventID(time.Now().UTC().Format("20060102T150405.000000000"))
	}
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}
	b.log = append(b.log, e)
	key := string(e.WorkUnitID)
	for _, ch := range b.subs[key] {
		select {
		case ch <- e:
		default:
		}
	}
	for _, ch := range b.subs[""] {
		select {
		case ch <- e:
		default:
		}
	}
	return nil
}

func (b *MemoryBus) Subscribe(ctx context.Context, workUnitFilter string) (<-chan Event, error) {
	ch := make(chan Event, 128)
	b.mu.Lock()
	b.subs[workUnitFilter] = append(b.subs[workUnitFilter], ch)
	b.mu.Unlock()
	go func() {
		<-ctx.Done()
	}()
	return ch, nil
}

func (b *MemoryBus) Replay(ctx context.Context, workUnitID types.WorkUnitID) ([]Event, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []Event
	for _, e := range b.log {
		if e.WorkUnitID == workUnitID {
			out = append(out, e)
		}
	}
	return out, nil
}

func (b *MemoryBus) Since(ctx context.Context, since int64) ([]Event, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []Event
	for _, e := range b.log {
		if e.Seq > since {
			out = append(out, e)
		}
	}
	return out, nil
}

func (b *MemoryBus) Seq() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.seq
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/eventbus",
		AESPSpecs:  []string{"AESP-0011", "AESP-0003"},
		Invariants: []string{"INV-10"},
		Status:     "implemented",
	}
}
