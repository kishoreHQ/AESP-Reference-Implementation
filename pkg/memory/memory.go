package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Store is the unified memory subsystem (INV-04). No per-runtime silos.
type Store interface {
	Write(ctx context.Context, item Item) error
	Query(ctx context.Context, q Query) ([]Item, error)
}

type Item struct {
	ID        string
	Tenant    types.TenantID
	Text      string
	Trust     types.TrustLabel
	Scope     string // working|session|semantic
	WorkUnit  types.WorkUnitID
	RuntimeID types.PluginID // recorded for audit; NOT a silo key for isolation of truth
}

type Query struct {
	Tenant types.TenantID
	Scope  string
	Text   string
	Limit  int
}

type Memory struct {
	mu    sync.Mutex
	items []Item
}

func New() *Memory { return &Memory{} }

func (m *Memory) Write(ctx context.Context, item Item) error {
	if item.Trust == "" {
		return fmt.Errorf("trust label required on every write (INV-04)")
	}
	if item.Trust == types.TrustUntrusted || item.Trust == types.TrustPoisonSuspect {
		// still stored, but never elevates privileges by itself
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = append(m.items, item)
	return nil
}

func (m *Memory) Query(ctx context.Context, q Query) ([]Item, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []Item
	for _, it := range m.items {
		if q.Tenant != "" && it.Tenant != q.Tenant {
			continue
		}
		if q.Scope != "" && it.Scope != q.Scope {
			continue
		}
		out = append(out, it)
		if q.Limit > 0 && len(out) >= q.Limit {
			break
		}
	}
	return out, nil
}

// MayAuthorizePrivileged is false for untrusted/retrieved/poison labels.
func MayAuthorizePrivileged(t types.TrustLabel) bool {
	switch t {
	case types.TrustSystem, types.TrustVerified:
		return true
	default:
		return false
	}
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/memory",
		AESPSpecs:  []string{"AESP-0004", "AESP-0013"},
		Invariants: []string{"INV-04"},
		Status:     "implemented",
	}
}
