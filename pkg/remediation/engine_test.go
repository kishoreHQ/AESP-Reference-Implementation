package remediation

import (
	"context"
	"testing"
)

func TestHandlePlaybook(t *testing.T) {
	e := New()
	inc, err := e.Handle(context.Background(), "wu", "service_down on api", SevHigh)
	if err != nil {
		t.Fatal(err)
	}
	if inc.Playbook != "pb.restart-service" || inc.Status != "resolved" {
		t.Fatalf("%+v", inc)
	}
}

func TestCriticalEscalates(t *testing.T) {
	e := New()
	inc, _ := e.Handle(context.Background(), "wu", "other", SevCrit)
	if !inc.NeedsHITL || inc.Status != "escalated" {
		t.Fatalf("%+v", inc)
	}
}
