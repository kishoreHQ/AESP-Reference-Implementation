package policy

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestUntrustedDenied(t *testing.T) {
	e := New()
	d := e.Evaluate(context.Background(), Request{
		SideEffect: types.SideEffectWriteRemote,
		Trust:      types.TrustUntrusted,
	})
	if d.Effect != Deny {
		t.Fatalf("%+v", d)
	}
}

func TestDestructiveNeedsApproval(t *testing.T) {
	e := New()
	d := e.Evaluate(context.Background(), Request{
		SideEffect: types.SideEffectDestructive,
		Trust:      types.TrustAgent,
	})
	if d.Effect != RequireApproval {
		t.Fatalf("%+v", d)
	}
}
