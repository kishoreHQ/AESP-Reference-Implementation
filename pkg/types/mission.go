package types

import "time"

// Mission is a host-submitted unit of work (Host Interface).
type Mission struct {
	ID              WorkUnitID            `json:"id"`
	Tenant          TenantID              `json:"tenant"`
	Goal            string                `json:"goal"`
	Constraints     []string              `json:"constraints,omitempty"`
	SuccessCriteria []string              `json:"successCriteria,omitempty"`
	RequiredCaps    []Capability          `json:"requiredCapabilities"`
	Budget          Budget                `json:"budget"`
	Labels          map[string]string     `json:"labels,omitempty"`
	CreatedAt       time.Time             `json:"createdAt"`
}

// Budget limits for a mission/session.
type Budget struct {
	MaxSteps   int     `json:"maxSteps"`
	MaxTokens  int64   `json:"maxTokens"`
	MaxCostUSD float64 `json:"maxCostUSD"`
	MaxWallSec int64   `json:"maxWallSec"`
}

// PlanArtifact is versioned planning output (INT-REQ-075/076).
type PlanArtifact struct {
	WorkUnitID      WorkUnitID `json:"workUnitId"`
	Revision        int        `json:"revision"`
	Goal            string     `json:"goal"`
	Steps           []PlanStep `json:"steps"`
	Assumptions     []string   `json:"assumptions,omitempty"`
	SuccessCriteria []string   `json:"successCriteria,omitempty"`
	Digest          ArtifactDigest `json:"digest,omitempty"`
}

type PlanStep struct {
	ID           string       `json:"id"`
	Description  string       `json:"description"`
	Capabilities []Capability `json:"capabilities,omitempty"`
	ToolHints    []string     `json:"toolHints,omitempty"`
}
