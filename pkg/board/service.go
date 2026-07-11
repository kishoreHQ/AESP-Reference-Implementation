package board

import (
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Column string

const (
	Backlog    Column = "backlog"
	Queued     Column = "queued"
	InProgress Column = "in_progress"
	Review     Column = "review"
	Done       Column = "done"
)

type Task struct {
	ID           string    `json:"id"`
	BoardID      string    `json:"boardId"`
	Title        string    `json:"title"`
	Column       Column    `json:"column"`
	Assignee     string    `json:"assignee,omitempty"` // agent id or empty = any capable
	Capabilities []string  `json:"capabilities,omitempty"`
	MissionID    string    `json:"missionId,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Board struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks,omitempty"`
}

type Service struct {
	mu     sync.Mutex
	boards map[string]*Board
	tasks  map[string]*Task
}

func New() *Service {
	s := &Service{boards: map[string]*Board{}, tasks: map[string]*Task{}}
	b := &Board{ID: "board_default", Name: "Command deck"}
	s.boards[b.ID] = b
	// seed
	s.addTask(b.ID, "Wire live session stream", Queued, []string{"coding", "tools"}, "")
	s.addTask(b.ID, "Review deploy plan", Review, []string{"tools"}, "")
	s.addTask(b.ID, "Triage provider health", Backlog, []string{"reasoning"}, "")
	return s
}

func (s *Service) addTask(boardID, title string, col Column, caps []string, assignee string) *Task {
	t := &Task{
		ID: fmt.Sprintf("task_%d", time.Now().UnixNano()), BoardID: boardID, Title: title,
		Column: col, Capabilities: caps, Assignee: assignee,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	s.tasks[t.ID] = t
	return t
}

func (s *Service) ListBoards() []Board {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Board, 0, len(s.boards))
	for _, b := range s.boards {
		bb := *b
		bb.Tasks = s.tasksFor(b.ID)
		out = append(out, bb)
	}
	return out
}

func (s *Service) tasksFor(boardID string) []Task {
	var out []Task
	for _, t := range s.tasks {
		if t.BoardID == boardID {
			out = append(out, *t)
		}
	}
	return out
}

func (s *Service) CreateTask(boardID, title string, caps []string, assignee string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.boards[boardID]; !ok {
		return nil, fmt.Errorf("board not found")
	}
	t := s.addTask(boardID, title, Backlog, caps, assignee)
	cp := *t
	return &cp, nil
}

func (s *Service) MoveTask(id string, col Column, assignee string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	t.Column = col
	if assignee != "" {
		t.Assignee = assignee
	}
	t.UpdatedAt = time.Now().UTC()
	cp := *t
	return &cp, nil
}

// Claim lets an agent pull next queued task (K4).
func (s *Service) Claim(agentID string, caps []string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.tasks {
		if t.Column != Queued {
			continue
		}
		if t.Assignee != "" && t.Assignee != agentID && t.Assignee != "any" {
			continue
		}
		t.Column = InProgress
		t.Assignee = agentID
		t.UpdatedAt = time.Now().UTC()
		cp := *t
		return &cp, nil
	}
	return nil, fmt.Errorf("no claimable task")
}

func (s *Service) ListTasks(boardID string) []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tasksFor(boardID)
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/board",
		AESPSpecs:  []string{"AESP-0005"},
		Invariants: []string{"INV-03", "INV-10"},
		Status:     "implemented",
	}
}
