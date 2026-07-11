package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/agentos"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/host"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Server exposes Host Interface over HTTP (P1/P3 profile).
type Server struct {
	sys *agentos.System
}

func New(sys *agentos.System) *Server {
	return &Server{sys: sys}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/v1/conformance", s.conformance)
	mux.HandleFunc("/v1/missions", s.missions)
	mux.HandleFunc("/v1/missions/", s.missionSub)
	mux.HandleFunc("/v1/approvals/", s.approvals)
	return mux
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if err := s.sys.Host().Health(r.Context()); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"status": "ok"})
}

func (s *Server) conformance(w http.ResponseWriter, r *http.Request) {
	// Avoid importing pkg/conformance here (import cycle with ModuleMappings).
	// Clients should use `aespd conformance` or inject a reporter.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"hint":   "run: aespd conformance",
		"status": "ok",
		"host":   "httpapi",
	})
}

func (s *Server) missions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var m types.Mission
	if err := json.Unmarshal(body, &m); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	res, err := s.sys.RunMission(r.Context(), m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, res)
}

func (s *Server) missionSub(w http.ResponseWriter, r *http.Request) {
	// /v1/missions/{id}/tree or /v1/missions/{id}
	path := strings.TrimPrefix(r.URL.Path, "/v1/missions/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "missing id", 400)
		return
	}
	id := types.WorkUnitID(parts[0])
	if len(parts) >= 2 && parts[1] == "tree" {
		tree, err := s.sys.Host().GetExecutionTree(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		writeJSON(w, tree)
		return
	}
	if len(parts) >= 2 && parts[1] == "events" {
		events, err := s.sys.Bus.Replay(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, events)
		return
	}
	m, ok := s.sys.Kernel.GetMission(id)
	if !ok {
		http.Error(w, "not found", 404)
		return
	}
	writeJSON(w, m)
}

func (s *Server) approvals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	// /v1/approvals/{taskId}
	taskID := types.HITLTaskID(strings.TrimPrefix(r.URL.Path, "/v1/approvals/"))
	var dec host.ApprovalDecision
	if err := json.NewDecoder(r.Body).Decode(&dec); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := s.sys.Approval.Resolve(r.Context(), taskID, dec.Approved, dec.Actor); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = s.sys.Host().ResolveApproval(r.Context(), taskID, dec)
	writeJSON(w, map[string]any{"ok": true})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// SpecMapping documents AESP coverage for the HTTP host adapter.
func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/httpapi",
		AESPSpecs:  []string{"AESP-0014", "AESP-0015", "AESP-0011"},
		Invariants: []string{"INV-11"},
		Status:     "implemented",
	}
}

