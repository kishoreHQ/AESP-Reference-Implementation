package policy

import (
	"context"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Effect string

const (
	Allow            Effect = "allow"
	Deny             Effect = "deny"
	RequireApproval  Effect = "require_approval"
)

type Request struct {
	Action         string
	Principal      types.PrincipalID
	Tenant         types.TenantID
	SideEffect     types.SideEffectClass
	Trust          types.TrustLabel
	Resource       string
}

type Decision struct {
	Effect      Effect
	Obligations []string
	Reason      string
}

type Engine struct{}

func New() *Engine { return &Engine{} }

// Evaluate implements fail-closed defaults for production side effects.
func (e *Engine) Evaluate(ctx context.Context, req Request) Decision {
	if req.Trust == types.TrustUntrusted || req.Trust == types.TrustPoisonSuspect {
		if req.SideEffect != types.SideEffectRead {
			return Decision{Effect: Deny, Reason: "untrusted content cannot authorize privileged actions"}
		}
	}
	switch req.SideEffect {
	case types.SideEffectDestructive, types.SideEffectAdmin:
		return Decision{Effect: RequireApproval, Reason: "destructive/admin requires HITL"}
	case types.SideEffectWriteRemote:
		return Decision{Effect: RequireApproval, Reason: "remote write requires HITL", Obligations: []string{"audit"}}
	default:
		return Decision{Effect: Allow, Reason: "default allow for read/workspace under policy"}
	}
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/policy",
		AESPSpecs:  []string{"AESP-0013", "AESP-0002", "AESP-0014"},
		Invariants: []string{"INV-06"},
		Status:     "implemented",
	}
}
