package capability

import (
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Engine resolves intent into required capabilities (INV-03).
// Never routes by vendor model name strings in kernel policy.
type Engine struct{}

func New() *Engine { return &Engine{} }

// FromMission extracts required capabilities from mission declaration.
func (e *Engine) FromMission(required []types.Capability) []types.Capability {
	out := make([]types.Capability, 0, len(required))
	seen := map[types.Capability]bool{}
	for _, c := range required {
		if c == "" || seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

// Compatible reports whether advertised caps satisfy required.
func Compatible(advertised, required []types.Capability) bool {
	set := map[types.Capability]bool{}
	for _, c := range advertised {
		set[c] = true
	}
	for _, r := range required {
		if !set[r] {
			return false
		}
	}
	return true
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/capability",
		AESPSpecs:  []string{"AESP-0015", "AESP-0001"},
		Invariants: []string{"INV-03"},
		Status:     "stubbed",
	}
}
