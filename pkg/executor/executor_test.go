package executor

import "testing"

func TestSpecMapping(t *testing.T) {
	if SpecMapping().Module != "pkg/executor" {
		t.Fatal("map")
	}
}
