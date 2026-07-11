package capability

import (
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestCompatible(t *testing.T) {
	if !Compatible([]types.Capability{"a", "b"}, []types.Capability{"a"}) {
		t.Fatal("expected compatible")
	}
	if Compatible([]types.Capability{"a"}, []types.Capability{"b"}) {
		t.Fatal("expected incompatible")
	}
}
