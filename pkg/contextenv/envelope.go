package contextenv

import "github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"

// Envelope is the unified context passed to every runtime (INV-05).
// Prompt is one field — not the entire context.
type Envelope struct {
	Workspace   WorkspaceContext       `json:"workspace"`
	Mission     MissionContext         `json:"mission"`
	Memory      []MemoryItem           `json:"memory"`
	Knowledge   []KnowledgeFact        `json:"knowledge"`
	Artifacts   []types.ArtifactDigest `json:"artifacts"`
	Policies    []PolicyObligation     `json:"policies"`
	Preferences map[string]string      `json:"preferences,omitempty"`
	Credentials []CredentialHandle     `json:"credentials"`
	Tools       []ToolSpec             `json:"tools"`
	Budget      types.Budget           `json:"budget"`
	Security    SecurityContext        `json:"security"`
	Prompt      string                 `json:"prompt"`
	Correlation Correlation            `json:"correlation"`
}

type WorkspaceContext struct {
	Root string `json:"root"`
	VCSRef string `json:"vcsRef,omitempty"`
}

type MissionContext struct {
	Goal            string   `json:"goal"`
	Constraints     []string `json:"constraints,omitempty"`
	SuccessCriteria []string `json:"successCriteria,omitempty"`
}

type MemoryItem struct {
	ID    string           `json:"id"`
	Text  string           `json:"text"`
	Trust types.TrustLabel `json:"trust"`
	Scope string           `json:"scope,omitempty"`
}

type KnowledgeFact struct {
	Subject string `json:"subject"`
	Predicate string `json:"predicate"`
	Object  string `json:"object"`
}

type PolicyObligation struct {
	ID   string `json:"id"`
	Kind string `json:"kind"` // deny | require-approval | redact | rate-limit
	Detail string `json:"detail,omitempty"`
}

type CredentialHandle struct {
	ID       string `json:"id"`
	Audience string `json:"audience"` // provider|runtime|tool plugin id
	// Raw secret material MUST NOT appear here.
}

type ToolSpec struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	InputSchema       map[string]any         `json:"inputSchema,omitempty"`
	SideEffectClass   types.SideEffectClass  `json:"sideEffectClass"`
	ApprovalRequired  bool                   `json:"approvalRequired"`
	RequiredCaps      []types.Capability     `json:"requiredCapabilities,omitempty"`
}

type SecurityContext struct {
	Tenant         types.TenantID   `json:"tenant"`
	Classification string           `json:"classification,omitempty"`
	Principal      types.PrincipalID `json:"principal"`
}

type Correlation struct {
	WorkUnitID types.WorkUnitID `json:"workUnitId"`
	SessionID  types.SessionID  `json:"sessionId"`
	TraceID    types.TraceID    `json:"traceId"`
}

// SpecMapping for conformance.
func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/contextenv",
		AESPSpecs:  []string{"AESP-0004", "AESP-0006", "AESP-0013", "AESP-0015"},
		Invariants: []string{"INV-05"},
		Status:     "stubbed",
	}
}
