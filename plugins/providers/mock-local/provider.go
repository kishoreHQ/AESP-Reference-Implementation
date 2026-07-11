package mocklocal

import (
	"context"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Provider is a local offline mock for P2 profiles (zero cloud credentials).
type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) ID() types.PluginID { return "provider.mock-local" }

func (p *Provider) Describe(ctx context.Context) (providerregistry.Descriptor, error) {
	return providerregistry.Descriptor{
		ID: p.ID(),
		Local: true,
		Capabilities: []types.Capability{"coding", "tools", "streaming", "local", "reasoning"},
		Models: []providerregistry.ModelInfo{{
			ID: "mock-local-small",
			Capabilities: []types.Capability{"coding", "tools", "local"},
		}},
	}, nil
}

func (p *Provider) Complete(ctx context.Context, req providerregistry.CompletionRequest) (providerregistry.CompletionResponse, error) {
	return providerregistry.CompletionResponse{
		ProviderID: p.ID(),
		ModelID:    "mock-local-small",
		Content:    "mock-local completion for: " + lastUser(req),
		TokensIn:   10,
		TokensOut:  20,
	}, nil
}

func (p *Provider) Health(ctx context.Context) error { return nil }

func lastUser(req providerregistry.CompletionRequest) string {
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" {
			return req.Messages[i].Content
		}
	}
	return ""
}
