package router

import (
	"context"
	"fmt"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/capability"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Decision is the result of capability-based routing (INV-03).
type Decision struct {
	ProviderID types.PluginID
	RuntimeID  types.PluginID
	ModelID    string
	Required   []types.Capability
	Reason     string
}

// Router selects provider + runtime without model-name hardcoding.
type Router struct {
	Providers *providerregistry.Registry
	Runtimes  *runtimeregistry.Registry
	Caps      *capability.Engine
}

func New(p *providerregistry.Registry, r *runtimeregistry.Registry) *Router {
	return &Router{Providers: p, Runtimes: r, Caps: capability.New()}
}

func (rt *Router) Route(ctx context.Context, required []types.Capability) (Decision, error) {
	req := rt.Caps.FromMission(required)
	if len(req) == 0 {
		return Decision{}, fmt.Errorf("no capabilities requested")
	}
	provs, err := rt.Providers.FindByCapabilities(ctx, req)
	if err != nil {
		return Decision{}, err
	}
	if len(provs) == 0 {
		return Decision{}, fmt.Errorf("no provider for capabilities %v", req)
	}
	// Prefer first match; production applies policy ranking / cost.
	p := provs[0]
	d, _ := p.Describe(ctx)
	model := ""
	if len(d.Models) > 0 {
		model = d.Models[0].ID
	}
	var runtimeID types.PluginID
	for _, r := range rt.Runtimes.List() {
		rd, err := r.Describe(ctx)
		if err != nil {
			continue
		}
		// Runtime must accept tools if tools required, etc. Simplified: any registered runtime.
		_ = rd
		runtimeID = r.ID()
		break
	}
	if runtimeID == "" {
		return Decision{}, fmt.Errorf("no runtime registered")
	}
	return Decision{
		ProviderID: p.ID(),
		RuntimeID:  runtimeID,
		ModelID:    model,
		Required:   req,
		Reason:     "capability-match",
	}, nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/router",
		AESPSpecs:  []string{"AESP-0015"},
		Invariants: []string{"INV-01", "INV-03"},
		Status:     "stubbed",
	}
}
