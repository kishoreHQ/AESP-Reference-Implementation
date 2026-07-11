package providerregistry

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type mockProv struct {
	id   types.PluginID
	caps []types.Capability
}

func (m *mockProv) ID() types.PluginID { return m.id }
func (m *mockProv) Describe(ctx context.Context) (Descriptor, error) {
	return Descriptor{ID: m.id, Capabilities: m.caps}, nil
}
func (m *mockProv) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	return CompletionResponse{ProviderID: m.id, ModelID: "mock-1", Content: "ok"}, nil
}
func (m *mockProv) Health(ctx context.Context) error { return nil }

func TestFindByCapabilities(t *testing.T) {
	r := New()
	_ = r.Register(&mockProv{id: "p.local", caps: []types.Capability{"coding", "local", "tools"}})
	_ = r.Register(&mockProv{id: "p.remote", caps: []types.Capability{"vision", "tools"}})
	found, err := r.FindByCapabilities(context.Background(), []types.Capability{"coding", "local"})
	if err != nil || len(found) != 1 || found[0].ID() != "p.local" {
		t.Fatalf("got %v err %v", found, err)
	}
}

func TestSpecMapping(t *testing.T) {
	m := SpecMapping()
	if len(m.Invariants) < 2 {
		t.Fatal("invariants")
	}
}
