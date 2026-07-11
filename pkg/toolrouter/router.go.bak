package toolrouter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/policy"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Tool interface {
	Name() string
	Spec() Spec
	Invoke(ctx context.Context, args map[string]any) (Result, error)
}

type Spec struct {
	Name            string
	Description     string
	SideEffectClass types.SideEffectClass
	ApprovalRequired bool
}

type Result struct {
	Output any
	Trust  types.TrustLabel
}

type InvocationRecord struct {
	ID         string
	Tool       string
	WorkUnitID types.WorkUnitID
	Args       map[string]any
	ResultTrust types.TrustLabel
	Allowed    bool
	At         time.Time
	Error      string
}

type Router struct {
	mu     sync.Mutex
	tools  map[string]Tool
	policy *policy.Engine
	log    []InvocationRecord
}

func New(p *policy.Engine) *Router {
	if p == nil {
		p = policy.New()
	}
	return &Router{tools: map[string]Tool{}, policy: p}
}

func (r *Router) Register(t Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[t.Name()] = t
}

func (r *Router) Invoke(ctx context.Context, workUnit types.WorkUnitID, name string, args map[string]any, trust types.TrustLabel) (Result, error) {
	r.mu.Lock()
	t, ok := r.tools[name]
	r.mu.Unlock()
	if !ok {
		return Result{}, fmt.Errorf("unknown tool %s", name)
	}
	sp := t.Spec()
	dec := r.policy.Evaluate(ctx, policy.Request{
		Action:     "tool.invoke",
		SideEffect: sp.SideEffectClass,
		Trust:      trust,
		Resource:   name,
	})
	rec := InvocationRecord{ID: fmt.Sprintf("ti_%d", time.Now().UnixNano()), Tool: name, WorkUnitID: workUnit, Args: args, At: time.Now().UTC()}
	if dec.Effect == policy.Deny {
		rec.Allowed = false
		rec.Error = dec.Reason
		r.append(rec)
		return Result{}, fmt.Errorf("denied: %s", dec.Reason)
	}
	if dec.Effect == policy.RequireApproval || sp.ApprovalRequired {
		rec.Allowed = false
		rec.Error = "requires_approval"
		r.append(rec)
		return Result{}, fmt.Errorf("requires_approval: %s", dec.Reason)
	}
	out, err := t.Invoke(ctx, args)
	if err != nil {
		rec.Error = err.Error()
		r.append(rec)
		return Result{}, err
	}
	if out.Trust == "" {
		out.Trust = types.TrustUntrusted
	}
	rec.Allowed = true
	rec.ResultTrust = out.Trust
	r.append(rec)
	return out, nil
}

func (r *Router) append(rec InvocationRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.log = append(r.log, rec)
}

func (r *Router) Records() []InvocationRecord {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]InvocationRecord{}, r.log...)
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/toolrouter",
		AESPSpecs:  []string{"AESP-0015", "AESP-0013", "AESP-0010"},
		Invariants: []string{"INV-06"},
		Status:     "stubbed",
	}
}
