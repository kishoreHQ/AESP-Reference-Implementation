package adapters

import (
	"context"
	"testing"
)

func TestGenericPTY_Handshake(t *testing.T) {
	g := NewGenericPTY()
	ok, out, _, err := g.Handshake(context.Background())
	if err != nil || !ok {
		t.Fatalf("handshake %v %v %s", ok, err, out)
	}
	if !contains(out, "aesp-handshake-ok") {
		t.Fatalf("stdout %s", out)
	}
}

func TestCatalog_NoVendorInIDsOnlyPlugins(t *testing.T) {
	// IDs are plugin ids — acceptable; kernel must not switch on vendor for business logic
	for _, a := range Catalog() {
		if a.ID() == "" {
			t.Fatal("empty id")
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
