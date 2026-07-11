package providerregistry

import (
	"context"
	"fmt"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type mockProv struct {
	id       types.PluginID
	caps     []types.Capability
	priority int
	fail     bool
}

func (m *mockProv) ID() types.PluginID { return m.id }
func (m *mockProv) Describe(ctx context.Context) (Descriptor, error) {
	return Descriptor{ID: m.id, Capabilities: m.caps, Priority: m.priority}, nil
}
func (m *mockProv) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	return CompletionResponse{ProviderID: m.id, ModelID: "mock-1", Content: "ok"}, nil
}
func (m *mockProv) Health(ctx context.Context) error {
	if m.fail {
		return fmt.Errorf("down")
	}
	return nil
}

func TestFindByCapabilities(t *testing.T) {
	r := New()
	_ = r.Register(&mockProv{id: "p.local", caps: []types.Capability{"coding", "local", "tools"}, priority: 1})
	_ = r.Register(&mockProv{id: "p.remote", caps: []types.Capability{"vision", "tools"}, priority: 10})
	found, err := r.FindByCapabilities(context.Background(), []types.Capability{"coding", "local"})
	if err != nil || len(found) != 1 || found[0].ID() != "p.local" {
		t.Fatalf("got %v err %v", found, err)
	}
}

func TestUnhealthyExcluded(t *testing.T) {
	r := New()
	_ = r.Register(&mockProv{id: "p.a", caps: []types.Capability{"coding"}, priority: 10, fail: true})
	_ = r.Register(&mockProv{id: "p.b", caps: []types.Capability{"coding"}, priority: 1})
	found, _ := r.FindByCapabilities(context.Background(), []types.Capability{"coding"})
	if len(found) != 1 || found[0].ID() != "p.b" {
		t.Fatalf("expected failover to p.b, got %v", found)
	}
}

func TestSpecMapping(t *testing.T) {
	if SpecMapping().Status != "implemented" {
		t.Fatal("status")
	}
}
