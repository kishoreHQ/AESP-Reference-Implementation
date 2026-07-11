package deploy

import (
	"context"
	"testing"
)

func TestDeployRollback(t *testing.T) {
	e := New()
	s, err := e.Start(context.Background(), "wu", "sha256:abc", "staging", "rolling")
	if err != nil {
		t.Fatal(err)
	}
	s, err = e.Complete(context.Background(), s.ID, false)
	if err != nil || s.Status != StatusFailed {
		t.Fatal(err, s)
	}
	s, err = e.Rollback(context.Background(), s.ID, "health failed")
	if err != nil || s.Status != StatusRolledBack {
		t.Fatal(err, s)
	}
}
