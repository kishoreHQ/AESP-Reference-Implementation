package a2a

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// AgentCard describes a peer agent (A2A-aligned, AESP-0015).
type AgentCard struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Capabilities []types.Capability `json:"capabilities"`
	URL          string             `json:"url,omitempty"`
}

type TaskState string

const (
	TaskSubmitted  TaskState = "submitted"
	TaskWorking    TaskState = "working"
	TaskCompleted  TaskState = "completed"
	TaskFailed     TaskState = "failed"
	TaskCanceled   TaskState = "canceled"
)

type Task struct {
	ID        string         `json:"id"`
	State     TaskState      `json:"state"`
	Message   string         `json:"message"`
	Result    string         `json:"result,omitempty"`
	Artifacts []string       `json:"artifacts,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// Registry of peer agents + in-process task exchange (golden fixture surface).
type Registry struct {
	mu     sync.Mutex
	cards  map[string]AgentCard
	tasks  map[string]*Task
	seq    int
}

func New() *Registry {
	return &Registry{cards: map[string]AgentCard{}, tasks: map[string]*Task{}}
}

func (r *Registry) Register(card AgentCard) error {
	if card.ID == "" {
		return fmt.Errorf("agent id required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cards[card.ID] = card
	return nil
}

func (r *Registry) GetCard(id string) (AgentCard, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.cards[id]
	return c, ok
}

func (r *Registry) ListCards() []AgentCard {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]AgentCard, 0, len(r.cards))
	for _, c := range r.cards {
		out = append(out, c)
	}
	return out
}

// SendTask creates a peer task and auto-completes for in-process demos.
func (r *Registry) SendTask(ctx context.Context, agentID, message string) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.cards[agentID]; !ok {
		return nil, fmt.Errorf("unknown agent %s", agentID)
	}
	r.seq++
	id := fmt.Sprintf("a2a_task_%d", r.seq)
	t := &Task{
		ID: id, State: TaskWorking, Message: message, CreatedAt: time.Now().UTC(),
	}
	// In-process peer: complete immediately with capability-safe result
	t.State = TaskCompleted
	t.Result = "peer-ok: " + message
	t.Artifacts = []string{"peer-result"}
	r.tasks[id] = t
	cp := *t
	return &cp, nil
}

func (r *Registry) GetTask(id string) (*Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, false
	}
	cp := *t
	return &cp, true
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/a2a",
		AESPSpecs:  []string{"AESP-0015", "AESP-0003"},
		Invariants: []string{"INV-06"},
		Status:     "implemented",
	}
}
