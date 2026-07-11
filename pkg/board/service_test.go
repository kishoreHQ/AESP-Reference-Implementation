package board

import "testing"

func TestClaimQueued(t *testing.T) {
	s := New()
	// ensure a queued task
	_, _ = s.CreateTask("board_default", "Claim me", []string{"coding"}, "")
	// move first backlog? seed has queued
	task, err := s.Claim("agent.builder", []string{"coding"})
	if err != nil {
		t.Fatal(err)
	}
	if task.Column != InProgress || task.Assignee != "agent.builder" {
		t.Fatalf("%+v", task)
	}
}
