package connections

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/credentials"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
)

func TestProbeAndRegisterGenericPTY(t *testing.T) {
	s := New(credentials.New(), providerregistry.New())
	cands := s.Probe(context.Background())
	if len(cands) == 0 {
		t.Fatal("expected candidates")
	}
	var found bool
	for _, c := range cands {
		if c.ID == "runtime.generic-pty" {
			found = true
		}
	}
	if !found {
		t.Fatal("generic-pty missing from probe")
	}
	conn, err := s.Register(context.Background(), RegisterRequest{
		Kind: KindRuntime, PluginID: "runtime.generic-pty", Name: "PTY",
	})
	if err != nil {
		t.Fatal(err)
	}
	if conn.Status != "connected" {
		t.Fatalf("status %s err %s", conn.Status, conn.LastError)
	}
	if !conn.Unsandboxed {
		t.Fatal("pty should be unsandboxed")
	}
}
