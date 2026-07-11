package goals

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/memory"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Goal struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Target    time.Time `json:"targetDate,omitempty"`
	Progress  float64   `json:"progress"`
	Missions  []string  `json:"linkedMissions,omitempty"`
	Tasks     []string  `json:"linkedTasks,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type JournalEntry struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	Trust     string    `json:"trust"`
}

// Service is memory-backed goals + journal (K6) — not a separate store.
type Service struct {
	mu      sync.Mutex
	goals   map[string]*Goal
	journal []JournalEntry
	mem     *memory.Memory
}

func New(mem *memory.Memory) *Service {
	s := &Service{goals: map[string]*Goal{}, mem: mem}
	s.goals["goal_1"] = &Goal{
		ID: "goal_1", Title: "Ship live agent cockpit", Progress: 0.35,
		CreatedAt: time.Now().UTC(), Missions: []string{"mis_live_demo"},
	}
	return s
}

func (s *Service) ListGoals() []Goal {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Goal, 0, len(s.goals))
	for _, g := range s.goals {
		out = append(out, *g)
	}
	return out
}

func (s *Service) CreateGoal(title string, target time.Time) *Goal {
	s.mu.Lock()
	defer s.mu.Unlock()
	g := &Goal{ID: fmt.Sprintf("goal_%d", time.Now().UnixNano()), Title: title, Target: target, CreatedAt: time.Now().UTC()}
	s.goals[g.ID] = g
	return g
}

func (s *Service) AdvanceFromMission(missionID string, delta float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, g := range s.goals {
		for _, m := range g.Missions {
			if m == missionID {
				g.Progress += delta
				if g.Progress > 1 {
					g.Progress = 1
				}
			}
		}
	}
}

func (s *Service) LinkMission(goalID, missionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, ok := s.goals[goalID]
	if !ok {
		return fmt.Errorf("goal not found")
	}
	g.Missions = append(g.Missions, missionID)
	return nil
}

func (s *Service) AddJournal(ctx context.Context, text string) JournalEntry {
	e := JournalEntry{
		ID: fmt.Sprintf("j_%d", time.Now().UnixNano()), Text: text,
		CreatedAt: time.Now().UTC(), Trust: string(types.TrustAgent),
	}
	s.mu.Lock()
	s.journal = append([]JournalEntry{e}, s.journal...)
	s.mu.Unlock()
	if s.mem != nil {
		_ = s.mem.Write(ctx, memory.Item{
			ID: e.ID, Text: text, Trust: types.TrustAgent, Scope: "session",
		})
	}
	return e
}

func (s *Service) ListJournal() []JournalEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]JournalEntry{}, s.journal...)
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/goals",
		AESPSpecs:  []string{"AESP-0004"},
		Invariants: []string{"INV-04", "INV-10"},
		Status:     "implemented",
	}
}
