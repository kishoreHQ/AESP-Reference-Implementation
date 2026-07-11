package agentos

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestRunMission_Single(t *testing.T) {
	sys := New(Config{AutoApprove: true, Workspace: t.TempDir()})
	res, err := sys.RunMission(context.Background(), types.Mission{
		ID:              "example.01-single-agent",
		Goal:            "Single-agent task",
		RequiredCaps:    []types.Capability{"coding", "tools"},
		SuccessCriteria: []string{"example-complete"},
		Budget:          types.Budget{MaxSteps: 10, MaxTokens: 1000},
		CreatedAt:       time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "succeeded" {
		t.Fatalf("status %s output %s", res.Status, res.Output)
	}
	if res.ProviderID == "" || res.RuntimeID == "" {
		t.Fatal("routing missing")
	}
	if len(res.Artifacts) < 2 {
		t.Fatal("expected artifacts")
	}
	// Events must include accept + route + complete
	typesSeen := map[string]bool{}
	for _, e := range res.Events {
		typesSeen[e.Type] = true
	}
	for _, need := range []string{"aesp.control.mission.accepted", "aesp.control.route.selected", "aesp.runtime.completed"} {
		if !typesSeen[need] {
			t.Fatalf("missing event %s in %v", need, typesSeen)
		}
	}
}

func TestRunMission_Failover(t *testing.T) {
	sys := New(Config{AutoApprove: true, Workspace: t.TempDir(), RemoteUnhealthy: false})
	res, err := sys.RunMission(context.Background(), types.Mission{
		ID:              "example.09-provider-failover",
		Goal:            "Provider-failover workflow",
		RequiredCaps:    []types.Capability{"coding", "tools"},
		SuccessCriteria: []string{"example-complete"},
		Budget:          types.Budget{MaxSteps: 10},
		Labels:          map[string]string{"scenario": "failover"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ProviderID != "provider.mock-local" {
		t.Fatalf("expected local after failover, got %s", res.ProviderID)
	}
	found := false
	for _, e := range res.Events {
		if e.Type == "aesp.provider.health.failed" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected health.failed event")
	}
}

func TestRunMission_MemoryAndKG(t *testing.T) {
	sys := New(Config{Workspace: t.TempDir()})
	_, err := sys.RunMission(context.Background(), types.Mission{
		ID: "example.05-memory-update", Goal: "Memory-update",
		RequiredCaps: []types.Capability{"reasoning", "tools"},
		SuccessCriteria: []string{"example-complete"},
		Budget: types.Budget{MaxSteps: 5},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = sys.RunMission(context.Background(), types.Mission{
		ID: "example.06-kg-update", Goal: "KG-update",
		RequiredCaps: []types.Capability{"reasoning"},
		SuccessCriteria: []string{"example-complete"},
		Budget: types.Budget{MaxSteps: 5},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunMission_HITLNoAutoApproveOnExpire(t *testing.T) {
	sys := New(Config{AutoApprove: false, Workspace: t.TempDir()})
	res, err := sys.RunMission(context.Background(), types.Mission{
		ID: "example.08-hitl", Goal: "HITL",
		RequiredCaps: []types.Capability{"tools"},
		SuccessCriteria: []string{"example-complete"},
		Budget: types.Budget{MaxSteps: 5},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "succeeded" {
		t.Fatal(res.Status, res.Output)
	}
}

func TestRunMission_Rollback(t *testing.T) {
	sys := New(Config{Workspace: t.TempDir()})
	res, err := sys.RunMission(context.Background(), types.Mission{
		ID: "example.10-rollback-retry", Goal: "Rollback",
		RequiredCaps: []types.Capability{"tools"},
		SuccessCriteria: []string{"example-complete"},
		Budget: types.Budget{MaxSteps: 10},
		Labels: map[string]string{"scenario": "rollback"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Output, "recovered") && res.Status != "succeeded" {
		t.Fatalf("%s %s", res.Status, res.Output)
	}
}

func TestRunMission_AllScenarios(t *testing.T) {
	cases := []struct {
		id   string
		caps []types.Capability
	}{
		{"example.01-single-agent", []types.Capability{"coding", "tools"}},
		{"example.02-multi-agent", []types.Capability{"coding", "tools", "planning"}},
		{"example.03-code-generation", []types.Capability{"coding", "tools"}},
		{"example.04-review-approval", []types.Capability{"coding"}},
		{"example.05-memory-update", []types.Capability{"reasoning", "tools"}},
		{"example.06-kg-update", []types.Capability{"reasoning"}},
		{"example.07-remediation", []types.Capability{"tools", "reasoning"}},
		{"example.08-hitl", []types.Capability{"tools"}},
		{"example.09-provider-failover", []types.Capability{"coding", "tools"}},
		{"example.10-rollback-retry", []types.Capability{"tools"}},
	}
	sys := New(Config{AutoApprove: true, Workspace: t.TempDir()})
	// reset remote healthy between runs for non-failover
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			sys.Remote.SetUnhealthy(false)
			sys.Providers.MarkHealthy("provider.mock-remote")
			labels := map[string]string{}
			if strings.Contains(tc.id, "failover") {
				labels["scenario"] = "failover"
			}
			if strings.Contains(tc.id, "rollback") {
				labels["scenario"] = "rollback"
			}
			res, err := sys.RunMission(context.Background(), types.Mission{
				ID: types.WorkUnitID(tc.id), Goal: tc.id,
				RequiredCaps: tc.caps, SuccessCriteria: []string{"example-complete"},
				Budget: types.Budget{MaxSteps: 15, MaxTokens: 2000},
				Labels: labels,
			})
			if err != nil {
				t.Fatal(err)
			}
			if res.Status != "succeeded" {
				t.Fatalf("status=%s out=%s", res.Status, res.Output)
			}
		})
	}
}
