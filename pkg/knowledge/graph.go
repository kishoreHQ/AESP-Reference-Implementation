package knowledge

import (
	"context"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Triple struct {
	Subject   string
	Predicate string
	Object    string
	Trust     types.TrustLabel
	WorkUnit  types.WorkUnitID
}

type Graph struct {
	mu      sync.Mutex
	triples []Triple
}

func New() *Graph { return &Graph{} }

func (g *Graph) Upsert(ctx context.Context, t Triple) error {
	if t.Trust == "" {
		t.Trust = types.TrustAgent
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.triples = append(g.triples, t)
	return nil
}

func (g *Graph) Query(ctx context.Context, subject string) ([]Triple, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []Triple
	for _, t := range g.triples {
		if subject == "" || t.Subject == subject {
			out = append(out, t)
		}
	}
	return out, nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/knowledge",
		AESPSpecs:  []string{"AESP-0006", "AESP-0004"},
		Invariants: []string{"INV-04"},
		Status:     "implemented",
	}
}
