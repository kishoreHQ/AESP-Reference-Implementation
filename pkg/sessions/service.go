package sessions

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/adapters"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/memory"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Status string

const (
	StatusIdle     Status = "idle"
	StatusWorking  Status = "working"
	StatusPaused   Status = "paused"
	StatusError    Status = "error"
	StatusStopped  Status = "stopped"
)

type Session struct {
	ID           string    `json:"id"`
	RuntimeID    string    `json:"runtimeId"`
	AgentID      string    `json:"agentId,omitempty"`
	ProviderID   string    `json:"providerId,omitempty"`
	Status       Status    `json:"status"`
	Unsandboxed  bool      `json:"unsandboxed"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	LastMessage  string    `json:"lastMessage,omitempty"`
	Model        string    `json:"model,omitempty"`
	Tokens       int64     `json:"tokens"`
	CostUSD      float64   `json:"costUsd"`
	ToolCalls    int       `json:"toolCalls"`
}

type Service struct {
	mu       sync.Mutex
	sessions map[string]*Session
	handles  map[string]adapters.SessionHandle
	bus      eventbus.Bus
	mem      *memory.Memory
	seqLocal int64
}

func New(bus eventbus.Bus, mem *memory.Memory) *Service {
	return &Service{
		sessions: map[string]*Session{},
		handles:  map[string]adapters.SessionHandle{},
		bus:      bus,
		mem:      mem,
	}
}

type CreateRequest struct {
	RuntimeID  string   `json:"runtime"`
	AgentID    string   `json:"agent,omitempty"`
	ProviderID string   `json:"provider,omitempty"`
	Caps       []string `json:"capabilities,omitempty"`
	Command    []string `json:"command,omitempty"`
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Session, error) {
	if req.RuntimeID == "" {
		req.RuntimeID = "runtime.generic-pty"
	}
	var ad adapters.RuntimeAdapter
	for _, a := range adapters.Catalog() {
		if a.ID() == req.RuntimeID {
			ad = a
			break
		}
	}
	if ad == nil {
		ad = adapters.NewGenericPTY()
		req.RuntimeID = ad.ID()
	}
	h, err := ad.Launch(ctx, adapters.LaunchOpts{Command: req.Command, Unsandboxed: req.RuntimeID == "runtime.generic-pty"})
	if err != nil {
		return nil, err
	}
	sess := &Session{
		ID: h.ID(), RuntimeID: req.RuntimeID, AgentID: req.AgentID, ProviderID: req.ProviderID,
		Status: StatusWorking, Unsandboxed: h.Unsandboxed(),
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	s.mu.Lock()
	s.sessions[sess.ID] = sess
	s.handles[sess.ID] = h
	s.mu.Unlock()

	s.emit(sess.ID, "session.status", map[string]any{"status": sess.Status})
	go s.pump(sess.ID, h)
	return sess, nil
}

func (s *Service) pump(id string, h adapters.SessionHandle) {
	for ev := range h.Events() {
		s.mu.Lock()
		sess := s.sessions[id]
		if sess != nil {
			sess.UpdatedAt = time.Now().UTC()
			switch ev.Kind {
			case "tool_call":
				sess.ToolCalls++
				s.emit(id, "session.tool_call", map[string]any{"tool": ev.Tool, "input": ev.Input, "output": ev.Output})
			case "model_switch":
				sess.Model = ev.Model
				s.emit(id, "session.model_switch", map[string]any{"model": ev.Model, "provider": ev.Provider})
			case "status":
				if ev.Text == "stopped" || ev.Text == "exited" {
					sess.Status = StatusStopped
				}
				s.emit(id, "session.status", map[string]any{"status": sess.Status, "text": ev.Text})
			case "output", "raw", "step", "error":
				s.emit(id, "session.output", map[string]any{"kind": ev.Kind, "text": ev.Text, "raw": ev.Raw})
			}
			sess.Tokens += ev.Tokens
			sess.CostUSD += ev.CostUSD
		}
		s.mu.Unlock()
	}
}

func (s *Service) Message(ctx context.Context, id, text string) error {
	s.mu.Lock()
	h := s.handles[id]
	sess := s.sessions[id]
	s.mu.Unlock()
	if h == nil || sess == nil {
		return fmt.Errorf("session not found")
	}
	sess.LastMessage = text
	sess.Status = StatusWorking
	sess.UpdatedAt = time.Now().UTC()
	// Unified memory write — agent-derived trust (never verified by default)
	if s.mem != nil {
		_ = s.mem.Write(ctx, memory.Item{
			ID: fmt.Sprintf("sess_msg_%d", time.Now().UnixNano()),
			Text: text, Trust: types.TrustAgent, Scope: "session",
			WorkUnit: types.WorkUnitID(id),
		})
	}
	s.emit(id, "session.output", map[string]any{"kind": "user", "text": text})
	// Simulate model switch for demos when provider set
	if sess.ProviderID != "" && sess.ToolCalls%3 == 1 {
		s.emit(id, "session.model_switch", map[string]any{
			"model": "routed-model", "provider": sess.ProviderID, "reason": "capability-route",
		})
	}
	return h.Send(ctx, text)
}

func (s *Service) Stop(ctx context.Context, id string) error {
	s.mu.Lock()
	h := s.handles[id]
	sess := s.sessions[id]
	s.mu.Unlock()
	if h == nil {
		return fmt.Errorf("session not found")
	}
	err := h.Stop(ctx)
	if sess != nil {
		sess.Status = StatusStopped
		s.emit(id, "session.status", map[string]any{"status": StatusStopped})
	}
	return err
}

func (s *Service) List() []Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Session, 0, len(s.sessions))
	for _, sess := range s.sessions {
		out = append(out, *sess)
	}
	return out
}

func (s *Service) Get(id string) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[id]
	if !ok {
		return nil, false
	}
	cp := *sess
	return &cp, true
}

func (s *Service) emit(sessionID, typ string, data map[string]any) {
	if s.bus == nil {
		return
	}
	if data == nil {
		data = map[string]any{}
	}
	data["sessionId"] = sessionID
	_ = s.bus.Publish(context.Background(), eventbus.Event{
		Type: typ, WorkUnitID: types.WorkUnitID(sessionID), Data: data, Time: time.Now().UTC(),
	})
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/sessions",
		AESPSpecs:  []string{"AESP-0001", "AESP-0004", "AESP-0015"},
		Invariants: []string{"INV-04", "INV-05", "INV-09", "INV-10"},
		Status:     "implemented",
	}
}
