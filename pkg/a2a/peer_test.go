package a2a

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestA2AGoldenFlow(t *testing.T) {
	r := New()
	_ = r.Register(AgentCard{
		ID: "agent.reviewer", Name: "Reviewer",
		Capabilities: []types.Capability{"coding"}, Description: "reviews code",
	})
	cards := r.ListCards()
	if len(cards) != 1 {
		t.Fatal(cards)
	}
	task, err := r.SendTask(context.Background(), "agent.reviewer", "review PR")
	if err != nil || task.State != TaskCompleted {
		t.Fatal(err, task)
	}
	// golden fixture shape
	fixture := filepath.Join("..", "..", "conformance", "fixtures", "a2a", "task-completed.json")
	if b, err := os.ReadFile(fixture); err == nil {
		var g map[string]any
		_ = json.Unmarshal(b, &g)
		if g["state"] != "completed" {
			t.Fatalf("golden state %v", g["state"])
		}
	}
}
