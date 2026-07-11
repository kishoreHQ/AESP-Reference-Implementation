package planner

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestPlan(t *testing.T) {
	p := New()
	plan := p.Plan(context.Background(), types.Mission{ID: "wu", Goal: "g", RequiredCaps: []types.Capability{"coding"}}, 1)
	if plan.Revision != 1 || len(plan.Steps) < 2 {
		t.Fatalf("%+v", plan)
	}
}
