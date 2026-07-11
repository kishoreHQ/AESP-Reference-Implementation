package kernel

import (
	"context"
	"testing"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestSubmitMission_RequiresCapabilities(t *testing.T) {
	k := New(nil)
	_, err := k.SubmitMission(context.Background(), types.Mission{
		ID:   "wu_1",
		Goal: "test",
	})
	if err == nil {
		t.Fatal("expected error when requiredCapabilities empty (INV-03)")
	}
}

func TestSubmitMission_OK(t *testing.T) {
	k := New(nil)
	id, err := k.SubmitMission(context.Background(), types.Mission{
		ID:           "wu_2",
		Goal:         "demo",
		RequiredCaps: []types.Capability{"coding"},
		Budget:       types.Budget{MaxSteps: 10},
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "wu_2" {
		t.Fatalf("id %s", id)
	}
	tree, err := k.GetExecutionTree(context.Background(), id)
	if err != nil || tree.WorkUnitID != id {
		t.Fatal("execution tree missing")
	}
}

func TestSpecMapping(t *testing.T) {
	if SpecMapping().Module != "pkg/kernel" {
		t.Fatal("mapping")
	}
}
