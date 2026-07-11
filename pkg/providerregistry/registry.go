package providerregistry

import (
	"context"
	"fmt"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Provider is a model-inference plugin (INV-01). Zero vendor names in core.
type Provider interface {
	ID() types.PluginID
	Describe(ctx context.Context) (Descriptor, error)
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
	Health(ctx context.Context) error
}

type Descriptor struct {
	ID           types.PluginID     `json:"id"`
	Capabilities []types.Capability `json:"capabilities"`
	Models       []ModelInfo        `json:"models"`
	Local        bool               `json:"local"`
	Priority     int                `json:"priority"` // higher preferred when healthy
}

type ModelInfo struct {
	ID           string             `json:"id"`
	Capabilities []types.Capability `json:"capabilities"`
}

type CompletionRequest struct {
	Model            string
	Messages         []Message
	Tools            []map[string]any
	MaxTokens        int
	CredentialHandle string
	Correlation      map[string]string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	ProviderID types.PluginID `json:"providerId"`
	ModelID    string         `json:"modelId"`
	Content    string         `json:"content"`
	TokensIn   int64          `json:"tokensIn"`
	TokensOut  int64          `json:"tokensOut"`
	CostUSD    float64        `json:"costUSD"`
}

type Registry struct {
	mu        sync.RWMutex
	byID      map[types.PluginID]Provider
	unhealthy map[types.PluginID]string
}

func New() *Registry {
	return &Registry{
		byID:      make(map[types.PluginID]Provider),
		unhealthy: make(map[types.PluginID]string),
	}
}

func (r *Registry) Register(p Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := p.ID()
	if id == "" {
		return fmt.Errorf("provider id required")
	}
	r.byID[id] = p
	return nil
}

func (r *Registry) Get(id types.PluginID) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.byID[id]
	return p, ok
}

func (r *Registry) List() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Provider, 0, len(r.byID))
	for _, p := range r.byID {
		out = append(out, p)
	}
	return out
}

func (r *Registry) MarkUnhealthy(id types.PluginID, reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.unhealthy[id] = reason
}

func (r *Registry) MarkHealthy(id types.PluginID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.unhealthy, id)
}

func (r *Registry) IsHealthy(id types.PluginID) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, bad := r.unhealthy[id]
	return !bad
}

// FindByCapabilities returns healthy providers advertising all required caps (INV-03).
// Sorted by priority desc; unhealthy excluded unless includeUnhealthy.
func (r *Registry) FindByCapabilities(ctx context.Context, required []types.Capability) ([]Provider, error) {
	type scored struct {
		p Provider
		pr int
	}
	var list []scored
	for _, p := range r.List() {
		if !r.IsHealthy(p.ID()) {
			continue
		}
		if err := p.Health(ctx); err != nil {
			r.MarkUnhealthy(p.ID(), err.Error())
			continue
		}
		d, err := p.Describe(ctx)
		if err != nil {
			continue
		}
		if hasAll(d.Capabilities, required) {
			list = append(list, scored{p: p, pr: d.Priority})
		}
	}
	// simple priority sort
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j].pr > list[i].pr {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	out := make([]Provider, len(list))
	for i, s := range list {
		out[i] = s.p
	}
	return out, nil
}

func hasAll(have, need []types.Capability) bool {
	set := map[types.Capability]bool{}
	for _, c := range have {
		set[c] = true
	}
	for _, n := range need {
		if !set[n] {
			return false
		}
	}
	return true
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/providerregistry",
		AESPSpecs:  []string{"AESP-0015"},
		Invariants: []string{"INV-01", "INV-02", "INV-03"},
		Status:     "implemented",
	}
}
