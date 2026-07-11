package replay

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
)

func TestJournal(t *testing.T) {
	j := New()
	_ = j.Append(context.Background(), eventbus.Event{Type: "aesp.runtime.step.completed", WorkUnitID: "wu"})
	if len(j.Events("wu")) != 1 {
		t.Fatal("missing")
	}
}
