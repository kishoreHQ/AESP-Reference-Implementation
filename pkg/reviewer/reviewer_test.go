package reviewer

import (
	"context"
	"testing"
)

func TestReview(t *testing.T) {
	r := New()
	o := r.Review(context.Background(), []string{"ok"}, "status ok")
	if !o.Passed {
		t.Fatal(o)
	}
}
