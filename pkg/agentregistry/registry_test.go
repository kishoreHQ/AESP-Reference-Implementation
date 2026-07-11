package agentregistry

import "testing"

func TestRegister(t *testing.T) {
	r := New()
	_ = r.Register(Agent{ID: "agent.1", Roles: []string{"builder"}})
	a, ok := r.Get("agent.1")
	if !ok || a.Status != "registered" {
		t.Fatal(a, ok)
	}
}
