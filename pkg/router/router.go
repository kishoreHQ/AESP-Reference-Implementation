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
	ProviderID   types.PluginID
	RuntimeID    types.PluginID
	ModelID      string
	Required     []types.Capability
	Reason       string
	FallbackFrom types.PluginID `json:"fallbackFrom,omitempty"`
	Candidates   []types.PluginID
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
		return Decision{}, fmt.Errorf("no healthy provider for capabilities %v", req)
	}
	cands := make([]types.PluginID, 0, len(provs))
	for _, p := range provs {
		cands = append(cands, p.ID())
	}
	p := provs[0]
	d, _ := p.Describe(ctx)
	model := ""
	if len(d.Models) > 0 {
		model = d.Models[0].ID
	}
	runtimeID, err := rt.pickRuntime(ctx, req)
	if err != nil {
		return Decision{}, err
	}
	return Decision{
		ProviderID: p.ID(),
		RuntimeID:  runtimeID,
		ModelID:    model,
		Required:   req,
		Reason:     "capability-match",
		Candidates: cands,
	}, nil
}

// RouteWithFailover marks failed provider unhealthy and re-routes.
func (rt *Router) RouteWithFailover(ctx context.Context, required []types.Capability, failed types.PluginID, reason string) (Decision, error) {
	if failed != "" {
		rt.Providers.MarkUnhealthy(failed, reason)
	}
	dec, err := rt.Route(ctx, required)
	if err != nil {
		return dec, err
	}
	dec.FallbackFrom = failed
	dec.Reason = "capability-match-failover"
	return dec, nil
}

func (rt *Router) pickRuntime(ctx context.Context, required []types.Capability) (types.PluginID, error) {
	var runtimeID types.PluginID
	for _, r := range rt.Runtimes.List() {
		if err := r.Health(ctx); err != nil {
			continue
		}
		rd, err := r.Describe(ctx)
		if err != nil {
			continue
		}
		// Prefer runtimes whose CapabilitiesIn cover tools if tools required.
		needTools := false
		for _, c := range required {
			if c == "tools" {
				needTools = true
			}
		}
		if needTools {
			ok := false
			for _, c := range rd.CapabilitiesIn {
				if c == "tools" {
					ok = true
				}
			}
			if !ok {
				continue
			}
		}
		runtimeID = r.ID()
		break
	}
	if runtimeID == "" {
		// fallback any healthy runtime
		for _, r := range rt.Runtimes.List() {
			if r.Health(ctx) == nil {
				return r.ID(), nil
			}
		}
		return "", fmt.Errorf("no runtime registered")
	}
	return runtimeID, nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/router",
		AESPSpecs:  []string{"AESP-0015"},
		Invariants: []string{"INV-01", "INV-03"},
		Status:     "implemented",
	}
}
