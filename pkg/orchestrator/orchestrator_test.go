package orchestrator

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/router"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type pMock struct{}
func (pMock) ID() types.PluginID { return "prov.mock" }
func (pMock) Describe(ctx context.Context) (providerregistry.Descriptor, error) {
	return providerregistry.Descriptor{ID: "prov.mock", Capabilities: []types.Capability{"coding"}}, nil
}
func (pMock) Complete(ctx context.Context, req providerregistry.CompletionRequest) (providerregistry.CompletionResponse, error) {
	return providerregistry.CompletionResponse{}, nil
}
func (pMock) Health(ctx context.Context) error { return nil }

type rMock struct{}
func (rMock) ID() types.PluginID { return "rt.mock" }
func (rMock) Describe(ctx context.Context) (runtimeregistry.Descriptor, error) {
	return runtimeregistry.Descriptor{ID: "rt.mock"}, nil
}
func (rMock) Execute(ctx context.Context, env contextenv.Envelope) (runtimeregistry.Result, error) {
	return runtimeregistry.Result{Status: "succeeded", Output: env.Mission.Goal}, nil
}
func (rMock) Health(ctx context.Context) error { return nil }

func TestRunMission(t *testing.T) {
	pr := providerregistry.New()
	rr := runtimeregistry.New()
	_ = pr.Register(pMock{})
	_ = rr.Register(rMock{})
	bus := eventbus.NewMemoryBus()
	o := New(bus, router.New(pr, rr))
	res, err := o.RunMission(context.Background(), types.Mission{
		ID: "wu", Goal: "build", RequiredCaps: []types.Capability{"coding"},
	}, contextenv.Envelope{})
	if err != nil || res.Status != "succeeded" {
		t.Fatal(err, res)
	}
}
