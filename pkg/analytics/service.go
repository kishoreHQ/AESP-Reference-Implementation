package analytics

import (
	"context"
	"strings"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/sessions"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Report is K7 rollup from journal + sessions — no second bookkeeping system.
type Report struct {
	AgentID      string         `json:"agentId"`
	Sessions     int            `json:"sessions"`
	ToolCalls    int            `json:"toolCalls"`
	Tokens       int64          `json:"tokens"`
	CostUSD      float64        `json:"costUsd"`
	ModelsUsed   map[string]int `json:"modelsUsed"`
	HourHistogram [24]int       `json:"hourHistogram"`
	ErrorRate    float64        `json:"errorRate"`
	Events       int            `json:"eventsScanned"`
}

type Service struct {
	bus  eventbus.Bus
	sess *sessions.Service
}

func New(bus eventbus.Bus, sess *sessions.Service) *Service {
	return &Service{bus: bus, sess: sess}
}

func (s *Service) Agent(ctx context.Context, agentID string) Report {
	r := Report{AgentID: agentID, ModelsUsed: map[string]int{}}
	if s.sess != nil {
		for _, sess := range s.sess.List() {
			if agentID != "" && sess.AgentID != agentID && sess.RuntimeID != agentID {
				continue
			}
			r.Sessions++
			r.ToolCalls += sess.ToolCalls
			r.Tokens += sess.Tokens
			r.CostUSD += sess.CostUSD
			if sess.Model != "" {
				r.ModelsUsed[sess.Model]++
			}
			h := sess.CreatedAt.UTC().Hour()
			r.HourHistogram[h]++
		}
	}
	// Scan journal
	if since, ok := s.bus.(interface {
		Since(ctx context.Context, since int64) ([]eventbus.Event, error)
	}); ok {
		evs, _ := since.Since(ctx, 0)
		errors := 0
		for _, e := range evs {
			r.Events++
			if strings.Contains(e.Type, "error") || strings.Contains(e.Type, "fail") {
				errors++
			}
			if e.Type == "session.model_switch" || strings.Contains(e.Type, "model_switch") {
				if m, ok := e.Data["model"].(string); ok && m != "" {
					r.ModelsUsed[m]++
				}
			}
			if e.Type == "session.tool_call" || strings.Contains(e.Type, "tool") {
				r.ToolCalls++
			}
			h := e.Time.UTC().Hour()
			r.HourHistogram[h]++
		}
		if r.Events > 0 {
			r.ErrorRate = float64(errors) / float64(r.Events)
		}
	}
	_ = types.TrustAgent
	return r
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/analytics",
		AESPSpecs:  []string{"AESP-0011", "AESP-0010"},
		Invariants: []string{"INV-10"},
		Status:     "implemented",
	}
}
