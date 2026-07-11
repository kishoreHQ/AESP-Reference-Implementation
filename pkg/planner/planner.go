package planner

import (
	"context"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Planner produces versioned plan artifacts (INT-REQ-075).
type Planner struct{}

func New() *Planner { return &Planner{} }

func (p *Planner) Plan(ctx context.Context, m types.Mission, revision int) types.PlanArtifact {
	steps := []types.PlanStep{{ID: "s1", Description: "analyze goal", Capabilities: m.RequiredCaps}}
	if len(m.RequiredCaps) > 0 {
		steps = append(steps, types.PlanStep{ID: "s2", Description: "execute with tools", Capabilities: m.RequiredCaps})
	}
	steps = append(steps, types.PlanStep{ID: "s3", Description: "verify success criteria"})
	return types.PlanArtifact{
		WorkUnitID: m.ID, Revision: revision, Goal: m.Goal,
		Steps: steps, SuccessCriteria: m.SuccessCriteria, Assumptions: m.Constraints,
	}
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:         "pkg/planner",
		AESPSpecs:      []string{"AESP-0015", "AESP-0005"},
		RequirementIDs: []string{"INT-REQ-075", "INT-REQ-076"},
		Invariants:     []string{"INV-10"},
		Status:         "stubbed",
	}
}
