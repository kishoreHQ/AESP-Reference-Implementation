package docgen

import (
	"context"
	"fmt"
	"strings"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Document is a generated documentation artifact (AESP-0008).
type Document struct {
	Title   string
	Body    string
	Format  string // markdown
	Source  string
	WorkUnit types.WorkUnitID
}

// Generator produces living docs from structured inputs.
type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) FromMission(ctx context.Context, m types.Mission, plan types.PlanArtifact, result string) Document {
	var b strings.Builder
	b.WriteString("# Mission Report\n\n")
	b.WriteString(fmt.Sprintf("**WorkUnit:** `%s`\n\n", m.ID))
	b.WriteString(fmt.Sprintf("**Goal:** %s\n\n", m.Goal))
	b.WriteString("## Required capabilities\n\n")
	for _, c := range m.RequiredCaps {
		b.WriteString(fmt.Sprintf("- `%s`\n", c))
	}
	b.WriteString("\n## Plan\n\n")
	for _, s := range plan.Steps {
		b.WriteString(fmt.Sprintf("1. **%s**: %s\n", s.ID, s.Description))
	}
	b.WriteString("\n## Result\n\n")
	b.WriteString(result)
	b.WriteString("\n")
	return Document{
		Title: "Mission Report: " + string(m.ID),
		Body: b.String(), Format: "markdown", Source: "docgen.mission", WorkUnit: m.ID,
	}
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/docgen",
		AESPSpecs:  []string{"AESP-0008", "AESP-0007"},
		Invariants: []string{"INV-10"},
		Status:     "implemented",
	}
}
