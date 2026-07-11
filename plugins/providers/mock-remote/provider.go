package mockremote

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Provider simulates a remote multi-capability provider for failover demos.
// No real vendor SDK — abstract plugin only.
type Provider struct {
	// Unhealthy when >0
	unhealthy atomic.Bool
}

func New() *Provider { return &Provider{} }

func (p *Provider) ID() types.PluginID { return "provider.mock-remote" }

func (p *Provider) SetUnhealthy(v bool) { p.unhealthy.Store(v) }

func (p *Provider) Describe(ctx context.Context) (providerregistry.Descriptor, error) {
	return providerregistry.Descriptor{
		ID:       p.ID(),
		Local:    false,
		Priority: 100, // preferred when healthy
		Capabilities: []types.Capability{"coding", "tools", "streaming", "vision", "reasoning", "planning"},
		Models: []providerregistry.ModelInfo{{
			ID:           "mock-remote-large",
			Capabilities: []types.Capability{"coding", "tools", "vision", "reasoning", "planning"},
		}},
	}, nil
}

func (p *Provider) Complete(ctx context.Context, req providerregistry.CompletionRequest) (providerregistry.CompletionResponse, error) {
	if p.unhealthy.Load() {
		return providerregistry.CompletionResponse{}, fmt.Errorf("provider.mock-remote unavailable")
	}
	return providerregistry.CompletionResponse{
		ProviderID: p.ID(),
		ModelID:    "mock-remote-large",
		Content:    "REMOTE_OK example-complete completed",
		TokensIn:   12,
		TokensOut:  30,
		CostUSD:    0.001,
	}, nil
}

func (p *Provider) Health(ctx context.Context) error {
	if p.unhealthy.Load() {
		return fmt.Errorf("provider.mock-remote health check failed")
	}
	return nil
}
