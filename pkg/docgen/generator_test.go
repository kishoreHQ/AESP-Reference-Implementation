package docgen

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestFromMission(t *testing.T) {
	g := New()
	d := g.FromMission(context.Background(), types.Mission{
		ID: "wu", Goal: "demo", RequiredCaps: []types.Capability{"coding"},
	}, types.PlanArtifact{Steps: []types.PlanStep{{ID: "s1", Description: "do"}}}, "ok")
	if d.Body == "" || d.Format != "markdown" {
		t.Fatal(d)
	}
}
