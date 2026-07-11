package approval

import (
	"context"
	"testing"
)

func TestExpireDoesNotApprove(t *testing.T) {
	s := New()
	id, _ := s.Request(context.Background(), "wu", "deploy")
	_ = s.Expire(context.Background(), id)
	task, _ := s.Get(id)
	if task.Status != Expired {
		t.Fatalf("status %s", task.Status)
	}
	if task.Status == Approved {
		t.Fatal("auto-approve forbidden")
	}
}
