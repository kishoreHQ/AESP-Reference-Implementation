package host

import "testing"

func TestSpecMapping_INV11(t *testing.T) {
	m := SpecMapping()
	ok := false
	for _, inv := range m.Invariants {
		if inv == "INV-11" {
			ok = true
		}
	}
	if !ok {
		t.Fatal("INV-11 required")
	}
}
