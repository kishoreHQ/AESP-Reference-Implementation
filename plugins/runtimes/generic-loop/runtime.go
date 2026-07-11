package genericloop

import (
	"context"
	"fmt"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Runtime is a minimal agent harness that uses the Context Envelope (INV-05).
type Runtime struct{}

func New() *Runtime { return &Runtime{} }

func (r *Runtime) ID() types.PluginID { return "runtime.generic-loop" }

func (r *Runtime) Describe(ctx context.Context) (runtimeregistry.Descriptor, error) {
	return runtimeregistry.Descriptor{
		ID: r.ID(), Version: "1.0.0",
		CapabilitiesIn:  []types.Capability{"tools", "streaming", "reasoning", "coding"},
		CapabilitiesOut: []types.Capability{"coding", "planning"},
		Sandbox:         "process",
	}, nil
}

func (r *Runtime) Execute(ctx context.Context, env contextenv.Envelope) (runtimeregistry.Result, error) {
	// Enforce stopping conditions minimally via budget.
	if env.Budget.MaxSteps <= 0 {
		env.Budget.MaxSteps = 5
	}
	out := fmt.Sprintf("completed goal=%q tools=%d memory=%d", env.Mission.Goal, len(env.Tools), len(env.Memory))
	return runtimeregistry.Result{
		Status: "succeeded", Output: out, StepsUsed: 1, TokensUsed: 1,
	}, nil
}

func (r *Runtime) Health(ctx context.Context) error { return nil }
