package remediation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Severity string

const (
	SevLow    Severity = "low"
	SevMedium Severity = "medium"
	SevHigh   Severity = "high"
	SevCrit   Severity = "critical"
)

type Incident struct {
	ID         string
	WorkUnitID types.WorkUnitID
	Alert      string
	Severity   Severity
	Status     string // open|healing|resolved|escalated
	Playbook   string
	Actions    []string
	CreatedAt  time.Time
	NeedsHITL  bool
}

type Playbook struct {
	ID       string
	Match    string // alert substring
	Actions  []string
	NeedHITL bool
	MaxBlast string
}

// Engine runs remediation playbooks (AESP-0012).
type Engine struct {
	mu         sync.Mutex
	playbooks  []Playbook
	incidents  map[string]*Incident
	seq        int
}

func New() *Engine {
	e := &Engine{incidents: map[string]*Incident{}}
	e.Register(Playbook{
		ID: "pb.restart-service", Match: "service_down",
		Actions: []string{"isolate", "restart", "verify"},
	})
	e.Register(Playbook{
		ID: "pb.budget-exhaust", Match: "budget",
		Actions: []string{"stop-loop", "notify"}, NeedHITL: true,
	})
	e.Register(Playbook{
		ID: "pb.generic", Match: "*",
		Actions: []string{"diagnose", "mitigate", "verify"},
	})
	return e
}

func (e *Engine) Register(p Playbook) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.playbooks = append(e.playbooks, p)
}

func (e *Engine) Handle(ctx context.Context, workUnit types.WorkUnitID, alert string, sev Severity) (*Incident, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	pb := e.match(alert)
	e.seq++
	id := fmt.Sprintf("inc_%d", e.seq)
	inc := &Incident{
		ID: id, WorkUnitID: workUnit, Alert: alert, Severity: sev,
		Status: "healing", Playbook: pb.ID, Actions: append([]string{}, pb.Actions...),
		CreatedAt: time.Now().UTC(), NeedsHITL: pb.NeedHITL || sev == SevCrit,
	}
	if inc.NeedsHITL {
		inc.Status = "escalated"
	} else {
		inc.Status = "resolved"
	}
	e.incidents[id] = inc
	return cloneInc(inc), nil
}

func (e *Engine) match(alert string) Playbook {
	var generic *Playbook
	for i := range e.playbooks {
		p := &e.playbooks[i]
		if p.Match == "*" {
			generic = p
			continue
		}
		if contains(alert, p.Match) {
			return *p
		}
	}
	if generic != nil {
		return *generic
	}
	return Playbook{ID: "pb.none", Actions: []string{"escalate"}}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func cloneInc(i *Incident) *Incident {
	cp := *i
	cp.Actions = append([]string{}, i.Actions...)
	return &cp
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/remediation",
		AESPSpecs:  []string{"AESP-0012", "AESP-0014", "AESP-0011"},
		Invariants: []string{"INV-10"},
		Status:     "implemented",
	}
}
