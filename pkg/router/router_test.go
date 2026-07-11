package router

import (
	"context"
	"fmt"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type pMock struct {
	id   types.PluginID
	caps []types.Capability
	pri  int
	fail bool
}

func (m *pMock) ID() types.PluginID { return m.id }
func (m *pMock) Describe(ctx context.Context) (providerregistry.Descriptor, error) {
	return providerregistry.Descriptor{
		ID: m.id, Capabilities: m.caps, Priority: m.pri,
		Models: []providerregistry.ModelInfo{{ID: "m1", Capabilities: m.caps}},
	}, nil
}
func (m *pMock) Complete(ctx context.Context, req providerregistry.CompletionRequest) (providerregistry.CompletionResponse, error) {
	return providerregistry.CompletionResponse{}, nil
}
func (m *pMock) Health(ctx context.Context) error {
	if m.fail {
		return fmt.Errorf("unhealthy")
	}
	return nil
}

type rMock struct{ id types.PluginID }

func (m *rMock) ID() types.PluginID { return m.id }
func (m *rMock) Describe(ctx context.Context) (runtimeregistry.Descriptor, error) {
	return runtimeregistry.Descriptor{ID: m.id, CapabilitiesIn: []types.Capability{"tools", "coding"}}, nil
}
func (m *rMock) Execute(ctx context.Context, env contextenv.Envelope) (runtimeregistry.Result, error) {
	return runtimeregistry.Result{Status: "succeeded"}, nil
}
func (m *rMock) Health(ctx context.Context) error { return nil }

func TestRoute(t *testing.T) {
	pr := providerregistry.New()
	rr := runtimeregistry.New()
	_ = pr.Register(&pMock{id: "prov.a", caps: []types.Capability{"coding", "tools"}, pri: 5})
	_ = rr.Register(&rMock{id: "rt.generic"})
	r := New(pr, rr)
	d, err := r.Route(context.Background(), []types.Capability{"coding"})
	if err != nil {
		t.Fatal(err)
	}
	if d.ProviderID != "prov.a" || d.RuntimeID != "rt.generic" {
		t.Fatalf("%+v", d)
	}
}

func TestFailover(t *testing.T) {
	pr := providerregistry.New()
	rr := runtimeregistry.New()
	_ = pr.Register(&pMock{id: "prov.hi", caps: []types.Capability{"coding"}, pri: 100})
	_ = pr.Register(&pMock{id: "prov.lo", caps: []types.Capability{"coding"}, pri: 1})
	_ = rr.Register(&rMock{id: "rt.generic"})
	r := New(pr, rr)
	d, err := r.RouteWithFailover(context.Background(), []types.Capability{"coding"}, "prov.hi", "simulated")
	if err != nil {
		t.Fatal(err)
	}
	if d.ProviderID != "prov.lo" {
		t.Fatalf("expected prov.lo got %s", d.ProviderID)
	}
	if d.FallbackFrom != "prov.hi" {
		t.Fatalf("fallbackFrom %s", d.FallbackFrom)
	}
}
