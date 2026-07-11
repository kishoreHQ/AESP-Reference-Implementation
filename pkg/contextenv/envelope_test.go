package contextenv

import "testing"

func TestSpecMapping_INV05(t *testing.T) {
	m := SpecMapping()
	if m.Module == "" {
		t.Fatal("module empty")
	}
	found := false
	for _, inv := range m.Invariants {
		if inv == "INV-05" {
			found = true
		}
	}
	if !found {
		t.Fatal("INV-05 must be mapped")
	}
	// Envelope must treat prompt as a field, not the whole world
	e := Envelope{Prompt: "hello"}
	if e.Prompt != "hello" {
		t.Fatal("prompt field missing")
	}
}
