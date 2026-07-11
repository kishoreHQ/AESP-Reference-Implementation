package routines

import (
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
	"github.com/robfig/cron/v3"
)

type Routine struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Schedule     string    `json:"schedule"` // cron
	Prompt       string    `json:"prompt,omitempty"`
	RuntimeID    string    `json:"runtimeId,omitempty"`
	Capabilities []string  `json:"capabilities,omitempty"`
	Paused       bool      `json:"paused"`
	LastRunAt    time.Time `json:"lastRunAt,omitempty"`
	LastStatus   string    `json:"lastStatus,omitempty"`
	LastMission  string    `json:"lastMissionId,omitempty"`
	NextFireAt   time.Time `json:"nextFireAt,omitempty"`
	History      []Run     `json:"history,omitempty"`
}

type Run struct {
	At        time.Time `json:"at"`
	MissionID string    `json:"missionId"`
	Status    string    `json:"status"`
}

type Service struct {
	mu       sync.Mutex
	routines map[string]*Routine
	parser   cron.Parser
}

func New() *Service {
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s := &Service{routines: map[string]*Routine{}, parser: p}
	// seed demo routine
	_, _ = s.Create(CreateReq{
		Name: "Hourly health sweep", Schedule: "@every 1h",
		Prompt: "Run provider health and log status", Capabilities: []string{"tools"},
	})
	return s
}

type CreateReq struct {
	Name         string   `json:"name"`
	Schedule     string   `json:"schedule"`
	Prompt       string   `json:"prompt"`
	RuntimeID    string   `json:"runtime"`
	Capabilities []string `json:"capabilities"`
}

func (s *Service) Create(req CreateReq) (*Routine, error) {
	if req.Name == "" || req.Schedule == "" {
		return nil, fmt.Errorf("name and schedule required")
	}
	sched, err := s.parser.Parse(req.Schedule)
	if err != nil {
		return nil, fmt.Errorf("invalid cron: %w", err)
	}
	r := &Routine{
		ID: fmt.Sprintf("rtn_%d", time.Now().UnixNano()), Name: req.Name, Schedule: req.Schedule,
		Prompt: req.Prompt, RuntimeID: req.RuntimeID, Capabilities: req.Capabilities,
		NextFireAt: sched.Next(time.Now().UTC()),
	}
	s.mu.Lock()
	s.routines[r.ID] = r
	s.mu.Unlock()
	return r, nil
}

func (s *Service) List() []Routine {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Routine, 0, len(s.routines))
	for _, r := range s.routines {
		s.refreshNext(r)
		out = append(out, *r)
	}
	return out
}

func (s *Service) refreshNext(r *Routine) {
	if r.Paused {
		return
	}
	if sched, err := s.parser.Parse(r.Schedule); err == nil {
		r.NextFireAt = sched.Next(time.Now().UTC())
	}
}

func (s *Service) NextFires(id string, n int) ([]time.Time, error) {
	s.mu.Lock()
	r, ok := s.routines[id]
	s.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	sched, err := s.parser.Parse(r.Schedule)
	if err != nil {
		return nil, err
	}
	var out []time.Time
	t := time.Now().UTC()
	for i := 0; i < n; i++ {
		t = sched.Next(t)
		out = append(out, t)
	}
	return out, nil
}

func (s *Service) Pause(id string, paused bool) (*Routine, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.routines[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	r.Paused = paused
	s.refreshNext(r)
	cp := *r
	return &cp, nil
}

// Fire marks a routine run (links mission id).
func (s *Service) Fire(id, missionID, status string) (*Routine, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.routines[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	r.LastRunAt = time.Now().UTC()
	r.LastStatus = status
	r.LastMission = missionID
	r.History = append(r.History, Run{At: r.LastRunAt, MissionID: missionID, Status: status})
	s.refreshNext(r)
	cp := *r
	return &cp, nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/routines",
		AESPSpecs:  []string{"AESP-0005", "AESP-0012"},
		Invariants: []string{"INV-03", "INV-10"},
		Status:     "implemented",
	}
}
