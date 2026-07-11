package agentregistry

import (
	"fmt"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Agent struct {
	ID     types.PrincipalID
	Roles  []string
	Status string
}

type Registry struct {
	mu   sync.RWMutex
	byID map[types.PrincipalID]*Agent
}

func New() *Registry { return &Registry{byID: map[types.PrincipalID]*Agent{}} }

func (r *Registry) Register(a Agent) error {
	if a.ID == "" {
		return fmt.Errorf("id required")
	}
	if a.Status == "" {
		a.Status = "registered"
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := a
	r.byID[a.ID] = &cp
	return nil
}

func (r *Registry) Get(id types.PrincipalID) (*Agent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.byID[id]
	return a, ok
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/agentregistry",
		AESPSpecs:  []string{"AESP-0001", "AESP-0002"},
		Invariants: []string{"INV-02"},
		Status:     "implemented",
	}
}
