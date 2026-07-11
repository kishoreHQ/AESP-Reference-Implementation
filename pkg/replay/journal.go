package replay

import (
	"context"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Journal stores events for deterministic control-plane replay (INV-10).
type Journal struct {
	mu  sync.Mutex
	log map[types.WorkUnitID][]eventbus.Event
}

func New() *Journal { return &Journal{log: map[types.WorkUnitID][]eventbus.Event{}} }

func (j *Journal) Append(ctx context.Context, e eventbus.Event) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.log[e.WorkUnitID] = append(j.log[e.WorkUnitID], e)
	return nil
}

func (j *Journal) Events(id types.WorkUnitID) []eventbus.Event {
	j.mu.Lock()
	defer j.mu.Unlock()
	return append([]eventbus.Event{}, j.log[id]...)
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/replay",
		AESPSpecs:  []string{"AESP-0011", "AESP-0005"},
		Invariants: []string{"INV-10"},
		Status:     "stubbed",
	}
}
