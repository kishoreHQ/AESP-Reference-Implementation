package deploy

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Status for deployment sessions (AESP-0009).
type Status string

const (
	StatusPending   Status = "pending"
	StatusRolling   Status = "rolling"
	StatusHealthy   Status = "healthy"
	StatusFailed    Status = "failed"
	StatusRolledBack Status = "rolled_back"
)

type Session struct {
	ID          string
	WorkUnitID  types.WorkUnitID
	Artifact    types.ArtifactDigest
	Environment string
	Status      Status
	Strategy    string // rolling|bluegreen|canary|recreate
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Evidence    []string
	Error       string
}

// Engine orchestrates deploy sessions with rollback (DEP-REQ family).
type Engine struct {
	mu       sync.Mutex
	sessions map[string]*Session
	seq      int
}

func New() *Engine {
	return &Engine{sessions: map[string]*Session{}}
}

func (e *Engine) Start(ctx context.Context, workUnit types.WorkUnitID, artifact types.ArtifactDigest, env, strategy string) (*Session, error) {
	if artifact == "" {
		return nil, fmt.Errorf("artifact digest required")
	}
	if strategy == "" {
		strategy = "rolling"
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.seq++
	id := fmt.Sprintf("dep_%d", e.seq)
	now := time.Now().UTC()
	s := &Session{
		ID: id, WorkUnitID: workUnit, Artifact: artifact, Environment: env,
		Status: StatusRolling, Strategy: strategy, CreatedAt: now, UpdatedAt: now,
		Evidence: []string{"gate:artifact-pinned"},
	}
	e.sessions[id] = s
	return clone(s), nil
}

func (e *Engine) Complete(ctx context.Context, id string, healthy bool) (*Session, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	s, ok := e.sessions[id]
	if !ok {
		return nil, fmt.Errorf("unknown session")
	}
	s.UpdatedAt = time.Now().UTC()
	if healthy {
		s.Status = StatusHealthy
		s.Evidence = append(s.Evidence, "gate:health-pass")
	} else {
		s.Status = StatusFailed
		s.Error = "health check failed"
	}
	return clone(s), nil
}

func (e *Engine) Rollback(ctx context.Context, id string, reason string) (*Session, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	s, ok := e.sessions[id]
	if !ok {
		return nil, fmt.Errorf("unknown session")
	}
	s.Status = StatusRolledBack
	s.Error = reason
	s.UpdatedAt = time.Now().UTC()
	s.Evidence = append(s.Evidence, "rollback:"+reason)
	return clone(s), nil
}

func (e *Engine) Get(id string) (*Session, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	s, ok := e.sessions[id]
	if !ok {
		return nil, false
	}
	return clone(s), true
}

func clone(s *Session) *Session {
	cp := *s
	cp.Evidence = append([]string{}, s.Evidence...)
	return &cp
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/deploy",
		AESPSpecs:  []string{"AESP-0009", "AESP-0010"},
		Invariants: []string{"INV-10"},
		Status:     "implemented",
	}
}
