package reviewer

import (
	"context"
	"strings"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Outcome struct {
	Passed  bool
	Findings []string
}

// Reviewer checks outputs against success criteria (eval + policy surface).
type Reviewer struct{}

func New() *Reviewer { return &Reviewer{} }

func (r *Reviewer) Review(ctx context.Context, criteria []string, output string) Outcome {
	var findings []string
	for _, c := range criteria {
		if c == "" {
			continue
		}
		if !strings.Contains(strings.ToLower(output), strings.ToLower(c)) {
			findings = append(findings, "missing criterion: "+c)
		}
	}
	return Outcome{Passed: len(findings) == 0, Findings: findings}
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/reviewer",
		AESPSpecs:  []string{"AESP-0010", "AESP-0014"},
		Invariants: []string{"INV-10"},
		Status:     "stubbed",
	}
}
