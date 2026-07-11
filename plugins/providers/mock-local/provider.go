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
		ID:       p.ID(),
		Local:    true,
		Priority: 10,
		Capabilities: []types.Capability{"coding", "tools", "streaming", "local", "reasoning", "planning"},
		Models: []providerregistry.ModelInfo{{
			ID:           "mock-local-small",
			Capabilities: []types.Capability{"coding", "tools", "local", "reasoning", "planning"},
		}},
	}, nil
}

func (p *Provider) Complete(ctx context.Context, req providerregistry.CompletionRequest) (providerregistry.CompletionResponse, error) {
	content := "LOCAL_OK: " + lastUser(req)
	// Include success-criteria friendly tokens for demos
	if contains(lastUser(req), "example-complete") || contains(lastUser(req), "goal") {
		content += " example-complete completed"
	}
	return providerregistry.CompletionResponse{
		ProviderID: p.ID(),
		ModelID:    "mock-local-small",
		Content:    content,
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

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && (s == sub || len(s) > 0 && indexOf(s, sub) >= 0))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
