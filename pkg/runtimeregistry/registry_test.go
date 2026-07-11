package runtimeregistry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "runtime.yaml")
	_ = os.WriteFile(p, []byte(`
apiVersion: aesp.runtime/v1
kind: RuntimePlugin
metadata:
  id: example.generic-loop
  version: 1.0.0
spec:
  capabilitiesIn: [tools, streaming]
  capabilitiesOut: [coding, planning]
  sandbox: process
  entrypoint: ./bin/runtime
`), 0o644)
	m, err := LoadManifest(p)
	if err != nil {
		t.Fatal(err)
	}
	if m.Metadata.ID != "example.generic-loop" {
		t.Fatalf("id %s", m.Metadata.ID)
	}
}

func TestSpecMapping(t *testing.T) {
	if SpecMapping().Status != "implemented" {
		t.Fatal("status")
	}
}
