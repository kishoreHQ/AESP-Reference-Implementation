package evaluation

import (
	"context"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Campaign is an offline eval harness distinct from the agent harness (AGENT-RUNTIME).
type Campaign struct {
	ID      string
	Cases   []Case
}

type Case struct {
	ID       string
	Input    string
	Expect   string
	Required []types.Capability
}

type Result struct {
	CaseID string
	Passed bool
	Score  float64
	Notes  string
}

type Harness struct{}

func New() *Harness { return &Harness{} }

func (h *Harness) Run(ctx context.Context, c Campaign, run func(context.Context, Case) (string, error)) []Result {
	var out []Result
	for _, tc := range c.Cases {
		got, err := run(ctx, tc)
		passed := err == nil && (tc.Expect == "" || got == tc.Expect)
		score := 0.0
		if passed {
			score = 1.0
		}
		note := ""
		if err != nil {
			note = err.Error()
		}
		out = append(out, Result{CaseID: tc.ID, Passed: passed, Score: score, Notes: note})
	}
	return out
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/evaluation",
		AESPSpecs:  []string{"AESP-0010"},
		Invariants: []string{"INV-10"},
		Status:     "stubbed",
	}
}
