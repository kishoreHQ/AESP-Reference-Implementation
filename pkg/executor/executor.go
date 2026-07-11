package executor

import (
	"context"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Executor runs plan steps via a runtime plugin.
type Executor struct {
	Runtime runtimeregistry.Runtime
}

func (e *Executor) ExecutePlan(ctx context.Context, plan types.PlanArtifact, env contextenv.Envelope) (runtimeregistry.Result, error) {
	env.Mission.Goal = plan.Goal
	env.Mission.SuccessCriteria = plan.SuccessCriteria
	return e.Runtime.Execute(ctx, env)
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/executor",
		AESPSpecs:  []string{"AESP-0005", "AESP-0015"},
		Invariants: []string{"INV-01", "INV-05"},
		Status:     "stubbed",
	}
}
