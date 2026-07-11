package eventbus

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestPublishReplay(t *testing.T) {
	b := NewMemoryBus()
	_ = b.Publish(context.Background(), Event{Type: "aesp.runtime.step.completed", WorkUnitID: "wu"})
	ev, err := b.Replay(context.Background(), types.WorkUnitID("wu"))
	if err != nil || len(ev) != 1 {
		t.Fatalf("replay %v %v", ev, err)
	}
}
