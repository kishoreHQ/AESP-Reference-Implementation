package types

// Correlation and identity types shared across the kernel.
// Vendor-neutral. AESP-0001 / ARCHITECTURE correlation keys.

type TenantID string
type PrincipalID string
type WorkUnitID string
type SessionID string
type TraceID string
type ArtifactDigest string
type HITLTaskID string
type EventID string
type PluginID string
type Capability string

// TrustLabel is applied to memory writes and tool results (INV-04, trust-model).
type TrustLabel string

const (
	TrustSystem        TrustLabel = "system"
	TrustVerified      TrustLabel = "verified"
	TrustAgent         TrustLabel = "agent"
	TrustRetrieved     TrustLabel = "retrieved"
	TrustUntrusted     TrustLabel = "untrusted"
	TrustPoisonSuspect TrustLabel = "poison-suspect"
)

// SideEffectClass classifies tool risk.
type SideEffectClass string

const (
	SideEffectRead           SideEffectClass = "read"
	SideEffectWriteWorkspace SideEffectClass = "write-workspace"
	SideEffectWriteRemote    SideEffectClass = "write-remote"
	SideEffectAdmin          SideEffectClass = "admin"
	SideEffectDestructive    SideEffectClass = "destructive"
)
