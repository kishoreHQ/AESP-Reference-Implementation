package memory

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestWriteRequiresTrust(t *testing.T) {
	m := New()
	err := m.Write(context.Background(), Item{Text: "x"})
	if err == nil {
		t.Fatal("expected trust required")
	}
}

func TestMayAuthorize(t *testing.T) {
	if MayAuthorizePrivileged(types.TrustUntrusted) {
		t.Fatal("untrusted must not authorize")
	}
	if !MayAuthorizePrivileged(types.TrustVerified) {
		t.Fatal("verified may authorize")
	}
}
