package host

import (
	"context"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Interface is the sole host interaction surface (INV-11).
// Kernel MUST NOT assume a specific UI or orchestrator.
type Interface interface {
	SubmitMission(ctx context.Context, m types.Mission) (types.WorkUnitID, error)
	CancelMission(ctx context.Context, id types.WorkUnitID, reason string) error
	SubscribeEvents(ctx context.Context, id types.WorkUnitID) (<-chan Event, error)
	ResolveApproval(ctx context.Context, taskID types.HITLTaskID, decision ApprovalDecision) error
	GetArtifact(ctx context.Context, digest types.ArtifactDigest) ([]byte, error)
	GetExecutionTree(ctx context.Context, id types.WorkUnitID) (*ExecutionTree, error)
	Health(ctx context.Context) error
}

// Event is a host-visible mission event.
type Event struct {
	Type       string         `json:"type"`
	ID         types.EventID  `json:"id"`
	WorkUnitID types.WorkUnitID `json:"workUnitId"`
	Data       map[string]any `json:"data,omitempty"`
}

type ApprovalDecision struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment,omitempty"`
	Actor    string `json:"actor"`
}

// ExecutionTree is the audit view for a mission (INV-10).
type ExecutionTree struct {
	WorkUnitID types.WorkUnitID `json:"workUnitId"`
	Agents     []string         `json:"agents"`
	Artifacts  []types.ArtifactDigest `json:"artifacts"`
	CostUSD    float64          `json:"costUSD"`
	Timeline   []Event          `json:"timeline"`
	Failures   []string         `json:"failures,omitempty"`
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/host",
		AESPSpecs:  []string{"AESP-0014", "AESP-0015", "AESP-0011"},
		Invariants: []string{"INV-11", "INV-10"},
		Status:     "stubbed",
	}
}
