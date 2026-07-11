package missionload

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "mission.yaml")
	_ = os.WriteFile(p, []byte(`
apiVersion: aesp.mission/v1
kind: Mission
metadata:
  id: example.test
spec:
  goal: "test"
  requiredCapabilities:
    - coding
  successCriteria:
    - example-complete
  budget:
    maxSteps: 5
`), 0o644)
	m, _, err := Load(p)
	if err != nil || m.ID != "example.test" || len(m.RequiredCaps) != 1 {
		t.Fatal(err, m)
	}
}

func TestRejectModelName(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.yaml")
	_ = os.WriteFile(p, []byte(`
metadata:
  id: bad
spec:
  goal: x
  requiredCapabilities: [gpt-4]
`), 0o644)
	_, _, err := Load(p)
	if err == nil {
		t.Fatal("expected INV-03 rejection")
	}
}
