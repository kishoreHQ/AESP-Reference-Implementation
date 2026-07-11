package mockremote

import (
	"context"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Provider simulates a remote multi-capability provider for failover demos.
// No real vendor SDK — abstract plugin only.
type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) ID() types.PluginID { return "provider.mock-remote" }

func (p *Provider) Describe(ctx context.Context) (providerregistry.Descriptor, error) {
	return providerregistry.Descriptor{
		ID: p.ID(),
		Local: false,
		Capabilities: []types.Capability{"coding", "tools", "streaming", "vision", "reasoning"},
		Models: []providerregistry.ModelInfo{{
			ID: "mock-remote-large",
			Capabilities: []types.Capability{"coding", "tools", "vision", "reasoning"},
		}},
	}, nil
}

func (p *Provider) Complete(ctx context.Context, req providerregistry.CompletionRequest) (providerregistry.CompletionResponse, error) {
	return providerregistry.CompletionResponse{
		ProviderID: p.ID(),
		ModelID:    "mock-remote-large",
		Content:    "mock-remote completion",
		TokensIn:   12,
		TokensOut:  30,
		CostUSD:    0.001,
	}, nil
}

func (p *Provider) Health(ctx context.Context) error { return nil }
