package orchestrator

import (
	"context"
	"fmt"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/router"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Orchestrator runs the durable control loop for a mission.
type Orchestrator struct {
	Bus    eventbus.Bus
	Router *router.Router
}

func New(bus eventbus.Bus, r *router.Router) *Orchestrator {
	return &Orchestrator{Bus: bus, Router: r}
}

// RunMission routes by capability and executes the selected runtime with a context envelope.
func (o *Orchestrator) RunMission(ctx context.Context, m types.Mission, env contextenv.Envelope) (runtimeregistry.Result, error) {
	if o.Router == nil {
		return runtimeregistry.Result{}, fmt.Errorf("router required")
	}
	dec, err := o.Router.Route(ctx, m.RequiredCaps)
	if err != nil {
		return runtimeregistry.Result{}, err
	}
	_ = o.Bus.Publish(ctx, eventbus.Event{
		Type: "aesp.control.route.selected", WorkUnitID: m.ID,
		Data: map[string]any{"providerId": string(dec.ProviderID), "runtimeId": string(dec.RuntimeID), "reason": dec.Reason},
	})
	rt, ok := o.Router.Runtimes.Get(dec.RuntimeID)
	if !ok {
		return runtimeregistry.Result{}, fmt.Errorf("runtime missing after route")
	}
	env.Correlation.WorkUnitID = m.ID
	env.Budget = m.Budget
	env.Mission = contextenv.MissionContext{Goal: m.Goal, Constraints: m.Constraints, SuccessCriteria: m.SuccessCriteria}
	res, err := rt.Execute(ctx, env)
	if err != nil {
		_ = o.Bus.Publish(ctx, eventbus.Event{Type: "aesp.runtime.failed", WorkUnitID: m.ID, Data: map[string]any{"error": err.Error()}})
		return res, err
	}
	_ = o.Bus.Publish(ctx, eventbus.Event{Type: "aesp.runtime.completed", WorkUnitID: m.ID, Data: map[string]any{"status": res.Status}})
	return res, nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/orchestrator",
		AESPSpecs:  []string{"AESP-0005", "AESP-0001", "AESP-0015"},
		Invariants: []string{"INV-03", "INV-05", "INV-01"},
		Status:     "stubbed",
	}
}
