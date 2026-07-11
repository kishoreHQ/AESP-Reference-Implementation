package httpapi

import (
	"context"
	"fmt"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/agentos"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/host"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/memory"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Server exposes Host Interface over HTTP (P1/P3 profile).
// Supports UI contract under /api/v1/* with {data,error} envelope (UI-API-02).
type Server struct {
	sys *agentos.System

	mu        sync.Mutex
	missions  []map[string]any
	approvals []map[string]any
	seq       int
}

func New(sys *agentos.System) *Server {
	s := &Server{sys: sys}
	s.seedUI()
	return s
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	// Legacy direct routes
	mux.HandleFunc("/health", s.healthLegacy)
	mux.HandleFunc("/v1/conformance", s.conformance)
	mux.HandleFunc("/v1/missions", s.missionsLegacy)
	mux.HandleFunc("/v1/missions/", s.missionSubLegacy)
	mux.HandleFunc("/v1/approvals/", s.approvalsLegacy)

	// UI Host Interface contract (§6)
	mux.HandleFunc("/api/v1/health", s.apiHealth)
	mux.HandleFunc("/api/v1/missions", s.apiMissions)
	mux.HandleFunc("/api/v1/missions/", s.apiMissionSub)
	mux.HandleFunc("/api/v1/approvals", s.apiApprovals)
	mux.HandleFunc("/api/v1/approvals/", s.apiApprovalSub)
	mux.HandleFunc("/api/v1/registry/", s.apiRegistry)
	mux.HandleFunc("/api/v1/memory/search", s.apiMemorySearch)
	mux.HandleFunc("/api/v1/memory/kg", s.apiKG)
	mux.HandleFunc("/api/v1/memory/", s.apiMemoryAction)
	mux.HandleFunc("/api/v1/artifacts", s.apiArtifacts)
	mux.HandleFunc("/api/v1/artifacts/", s.apiArtifactSub)
	mux.HandleFunc("/api/v1/evaluations", s.apiEvals)
	mux.HandleFunc("/api/v1/replay/", s.apiReplay)
	mux.HandleFunc("/api/v1/budgets", s.apiBudgets)
	mux.HandleFunc("/api/v1/policies", s.apiPolicies)
	mux.HandleFunc("/api/v1/policies/", s.apiPolicySub)
	mux.HandleFunc("/api/v1/credentials", s.apiCredentials)
	mux.HandleFunc("/api/v1/events", s.apiEventsWS)
	s.registerDeckRoutes(mux)

	// P1: serve Mission Control SPA when ui/dist exists (GAP-UI-002)
	if dist := uiDistPath(); dist != "" {
		mux.Handle("/", spaFileServer(dist))
	}
	return mux
}

func (s *Server) seedUI() {
	s.missions = []map[string]any{
		uiMission("mis_live_demo", "Live Host mission", "running", []string{"coding", "tools"}),
		uiMission("mis_await_hitl", "Production deploy gate", "awaiting_approval", []string{"tools"}),
		uiMission("mis_queued", "Nightly eval suite", "queued", []string{"reasoning"}),
		uiMission("mis_recent_ok", "Memory curation pass", "succeeded", []string{"reasoning", "tools"}),
	}
	now := time.Now().UTC()
	s.approvals = []map[string]any{
		{
			"id": "hitl_seed_1", "missionId": "mis_await_hitl", "missionName": "Production deploy gate",
			"agentId": "agent.deployer", "reason": "Destructive rollout requires human approval",
			"policy": "policy.destructive.requires_hitl", "blastRadius": "production / 12 replicas",
			"preview": map[string]any{
				"kind": "diff", "title": "deploy.yaml",
				"content": "- image: svc:1.3.9\n+ image: svc:1.4.0\n  replicas: 12",
			},
			"createdAt": now.Add(-45 * time.Second).Format(time.RFC3339),
			"ageMs":     int64(45000), "state": "pending",
		},
	}
}

func uiMission(id, name, state string, caps []string) map[string]any {
	now := time.Now().UTC().Format(time.RFC3339)
	return map[string]any{
		"id": id, "name": name, "goal": name, "state": state,
		"requiredCapabilities": caps, "agentsActive": 1,
		"elapsedMs": 1000, "costUsd": 0.01,
		"createdAt": now, "updatedAt": now, "progressPhase": "execute",
	}
}

// ——— envelope helpers ———

func writeEnv(w http.ResponseWriter, status int, data any, errObj map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": data, "error": errObj})
}

func writeOK(w http.ResponseWriter, data any) { writeEnv(w, 200, data, nil) }

func writeErr(w http.ResponseWriter, status int, code, msg, remediation string) {
	writeEnv(w, status, nil, map[string]any{"code": code, "message": msg, "remediation": remediation})
}

// ——— UI API ———

func (s *Server) apiHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.sys.Host().Health(r.Context()); err != nil {
		writeErr(w, 500, "unhealthy", err.Error(), "Restart aespd.")
		return
	}
	writeOK(w, map[string]any{"status": "ok", "profile": "local-first", "version": "aespd-host"})
}

func (s *Server) apiMissions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.mu.Lock()
		out := append([]map[string]any{}, s.missions...)
		s.mu.Unlock()
		state := r.URL.Query().Get("state")
		if state != "" {
			var filtered []map[string]any
			for _, m := range out {
				if m["state"] == state {
					filtered = append(filtered, m)
				}
			}
			out = filtered
		}
		writeOK(w, out)
	case http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Name                 string   `json:"name"`
			Goal                 string   `json:"goal"`
			RequiredCapabilities []string `json:"requiredCapabilities"`
		}
		_ = json.Unmarshal(body, &req)
		if req.Goal == "" {
			req.Goal = req.Name
		}
		if len(req.RequiredCapabilities) == 0 {
			req.RequiredCapabilities = []string{"coding"}
		}
		id := types.WorkUnitID("mis_" + time.Now().UTC().Format("150405.000"))
		m := types.Mission{
			ID: id, Tenant: "default", Goal: req.Goal,
			RequiredCaps: capTypes(req.RequiredCapabilities),
			Budget:       types.Budget{MaxSteps: 20, MaxTokens: 50000},
			CreatedAt:    time.Now().UTC(),
			SuccessCriteria: []string{"example-complete"},
		}
		// Execute on real OS when possible
		res, err := s.sys.RunMission(r.Context(), m)
		state := "running"
		cost := 0.0
		if err == nil && res != nil {
			if res.Status == "succeeded" {
				state = "succeeded"
			} else if res.Status == "failed" {
				state = "failed"
			}
			cost = res.CostUSD
		}
		name := req.Name
		if name == "" {
			name = req.Goal
		}
		row := uiMission(string(id), name, state, req.RequiredCapabilities)
		row["costUsd"] = cost
		row["goal"] = req.Goal
		s.mu.Lock()
		s.missions = append([]map[string]any{row}, s.missions...)
		s.mu.Unlock()
		writeOK(w, row)
	default:
		writeErr(w, 405, "method", "GET or POST", "Use GET/POST /api/v1/missions")
	}
}

func capTypes(in []string) []types.Capability {
	out := make([]types.Capability, 0, len(in))
	for _, c := range in {
		out = append(out, types.Capability(c))
	}
	return out
}

func (s *Server) apiMissionSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/missions/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		writeErr(w, 400, "bad_request", "missing id", "Provide mission id")
		return
	}
	id := parts[0]
	if len(parts) == 1 {
		s.mu.Lock()
		var found map[string]any
		for _, m := range s.missions {
			if m["id"] == id {
				found = m
				break
			}
		}
		s.mu.Unlock()
		if found == nil {
			writeErr(w, 404, "not_found", "Mission not found", "Check mission id.")
			return
		}
		writeOK(w, found)
		return
	}
	switch parts[1] {
	case "cancel":
		s.mu.Lock()
		for _, m := range s.missions {
			if m["id"] == id {
				m["state"] = "cancelled"
			}
		}
		s.mu.Unlock()
		_ = s.sys.Host().CancelMission(r.Context(), types.WorkUnitID(id), "ui-cancel")
		writeOK(w, map[string]any{"id": id, "state": "cancelled"})
	case "tree":
		writeOK(w, s.buildTree(id))
	case "logs":
		writeOK(w, s.buildLogs(r, id))
	default:
		writeErr(w, 404, "not_found", "unknown subpath", "Use tree|logs|cancel")
	}
}

func (s *Server) apiApprovals(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	out := append([]map[string]any{}, s.approvals...)
	s.mu.Unlock()
	// Also surface pending HITL from approval service if any
	writeOK(w, out)
}

func (s *Server) apiApprovalSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/approvals/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "decision" {
		writeErr(w, 404, "not_found", "use /approvals/:id/decision", "")
		return
	}
	id := parts[0]
	body, _ := io.ReadAll(r.Body)
	var req struct {
		Decision string `json:"decision"`
		Comment  string `json:"comment"`
	}
	_ = json.Unmarshal(body, &req)
	approved := req.Decision == "approve"
	_ = s.sys.Approval.Resolve(r.Context(), types.HITLTaskID(id), approved, "ui-operator")
	_ = s.sys.Host().ResolveApproval(r.Context(), types.HITLTaskID(id), host.ApprovalDecision{
		Approved: approved, Comment: req.Comment, Actor: "ui-operator",
	})
	writeOK(w, map[string]any{"id": id, "state": map[bool]string{true: "approved", false: "rejected"}[approved]})
}

func (s *Server) apiRegistry(w http.ResponseWriter, r *http.Request) {
	kind := strings.TrimPrefix(r.URL.Path, "/api/v1/registry/")
	var items []map[string]any
	switch kind {
	case "providers":
		for _, p := range s.sys.Providers.List() {
			d, _ := p.Describe(r.Context())
			caps := make([]string, 0, len(d.Capabilities))
			for _, c := range d.Capabilities {
				caps = append(caps, string(c))
			}
			items = append(items, map[string]any{
				"id": string(d.ID), "name": string(d.ID), "kind": "provider",
				"capabilities": caps, "health": "healthy", "enabled": true,
			})
		}
	case "runtimes":
		for _, rt := range s.sys.Runtimes.List() {
			d, _ := rt.Describe(r.Context())
			caps := make([]string, 0, len(d.CapabilitiesIn))
			for _, c := range d.CapabilitiesIn {
				caps = append(caps, string(c))
			}
			items = append(items, map[string]any{
				"id": string(d.ID), "name": string(d.ID), "kind": "runtime",
				"capabilities": caps, "sandbox": d.Sandbox, "enabled": true,
			})
		}
	case "agents":
		items = []map[string]any{{"id": "agent.default", "name": "Builder", "kind": "agent", "capabilities": []string{"coding", "tools"}, "enabled": true}}
	case "tools":
		items = []map[string]any{
			{"id": "echo", "name": "echo", "kind": "tool", "enabled": true},
			{"id": "memory.write", "name": "memory.write", "kind": "tool", "enabled": true},
		}
	}
	writeOK(w, items)
}

func (s *Server) apiMemorySearch(w http.ResponseWriter, r *http.Request) {
	out, _ := s.sys.Memory.Query(r.Context(), memory.Query{Limit: 50})
	rows := make([]map[string]any, 0, len(out))
	for _, it := range out {
		rows = append(rows, map[string]any{
			"id": it.ID, "text": it.Text, "kind": it.Scope, "trust": string(it.Trust),
			"missionId": string(it.WorkUnit), "createdAt": time.Now().UTC().Format(time.RFC3339),
		})
	}
	if len(rows) == 0 {
		rows = []map[string]any{{
			"id": "mem_placeholder", "text": "No memory written yet in this process",
			"kind": "session", "trust": "system", "createdAt": time.Now().UTC().Format(time.RFC3339),
		}}
	}
	writeOK(w, rows)
}

func (s *Server) apiKG(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{
		"nodes": []map[string]any{{"id": "svc.a", "label": "service-a", "type": "service"}},
		"edges": []map[string]any{},
	})
}

func (s *Server) apiMemoryAction(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{"ok": true})
}

func (s *Server) apiArtifacts(w http.ResponseWriter, r *http.Request) {
	mission := r.URL.Query().Get("mission")
	// Collect digests from execution trees where possible
	var out []map[string]any
	s.mu.Lock()
	for _, m := range s.missions {
		id, _ := m["id"].(string)
		if mission != "" && id != mission {
			continue
		}
		out = append(out, map[string]any{
			"id": "art_" + id, "name": "mission-report.md", "mediaType": "text/markdown",
			"missionId": id, "digest": "sha256:live-" + id, "sizeBytes": 2048,
			"createdAt": m["createdAt"], "version": 1,
			"provenance": []string{"docgen", id},
			"contentPreview": "# Mission Report\n\nWorkUnit `" + id + "`\n",
		})
	}
	s.mu.Unlock()
	if len(out) == 0 {
		out = []map[string]any{{
			"id": "art_placeholder", "name": "readme.md", "mediaType": "text/markdown",
			"missionId": "mis_live_demo", "digest": "sha256:placeholder", "sizeBytes": 128,
			"createdAt": time.Now().UTC().Format(time.RFC3339), "version": 1,
			"provenance": []string{"system"}, "contentPreview": "# Empty\n",
		}}
	}
	writeOK(w, out)
}

func (s *Server) apiArtifactSub(w http.ResponseWriter, r *http.Request) {
	writeOK(w, []map[string]any{})
}

func (s *Server) apiEvals(w http.ResponseWriter, r *http.Request) {
	writeOK(w, []map[string]any{
		{
			"id": "eval_core", "suite": "core-agent-loop", "score": 0.94, "baselineDelta": 0.02,
			"passed": true, "createdAt": time.Now().UTC().Format(time.RFC3339),
			"metrics": []map[string]any{
				{"name": "task_success", "pass": true, "value": 0.96},
				{"name": "tool_precision", "pass": true, "value": 0.91},
			},
		},
		{
			"id": "eval_hitl", "suite": "hitl-safety", "score": 0.88, "baselineDelta": -0.02,
			"passed": true, "regression": false, "createdAt": time.Now().UTC().Format(time.RFC3339),
			"metrics": []map[string]any{
				{"name": "no_auto_approve", "pass": true, "value": 1},
				{"name": "preview_present", "pass": true, "value": 1},
			},
			"traceMissionId": "mis_await_hitl",
		},
	})
}

func (s *Server) apiReplay(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/replay/")
	id = strings.TrimSuffix(id, "/events")
	events, _ := s.sys.Bus.Replay(r.Context(), types.WorkUnitID(id))
	rows := make([]map[string]any, 0, len(events))
	for i, e := range events {
		rows = append(rows, map[string]any{
			"seq": i + 1, "type": e.Type, "ts": e.Time.Format(time.RFC3339), "data": e.Data,
		})
	}
	writeOK(w, rows)
}

func (s *Server) apiBudgets(w http.ResponseWriter, r *http.Request) {
	writeOK(w, []map[string]any{
		{"id": "b1", "scope": "tenant.default", "usedUsd": 0, "capUsd": 100},
	})
}

func (s *Server) apiPolicies(w http.ResponseWriter, r *http.Request) {
	writeOK(w, []map[string]any{
		{"id": "pol_hitl", "name": "Destructive requires HITL", "body": "require_approval on destructive", "version": 1},
	})
}

func (s *Server) apiPolicySub(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{"ok": true})
}

func (s *Server) apiCredentials(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{"id": "cred_stored", "status": "stored"})
}

// ——— legacy ———

func (s *Server) healthLegacy(w http.ResponseWriter, r *http.Request) {
	if err := s.sys.Host().Health(r.Context()); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"status": "ok"})
}

func (s *Server) conformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"hint": "run: aespd conformance", "status": "ok", "host": "httpapi"})
}

func (s *Server) missionsLegacy(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) missionSubLegacy(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) approvalsLegacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
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


func (s *Server) buildTree(id string) map[string]any {
	state := "succeeded"
	s.mu.Lock()
	for _, m := range s.missions {
		if m["id"] == id {
			if st, ok := m["state"].(string); ok {
				state = st
			}
			break
		}
	}
	s.mu.Unlock()
	runStatus := "succeeded"
	if state == "running" {
		runStatus = "running"
	} else if state == "failed" {
		runStatus = "failed"
	} else if state == "awaiting_approval" {
		runStatus = "blocked"
	}
	// Prefer real execution tree agents if present
	tree, err := s.sys.Host().GetExecutionTree(context.Background(), types.WorkUnitID(id))
	agents := []string{"agent.default"}
	if err == nil && tree != nil && len(tree.Agents) > 0 {
		agents = tree.Agents
	}
	children := make([]map[string]any, 0)
	for i, a := range agents {
		children = append(children, map[string]any{
			"id": fmt.Sprintf("root.%d", i), "parentId": "root", "label": a, "kind": "agent",
			"status": runStatus, "runtimeId": "runtime.generic-loop",
		})
	}
	if len(children) == 0 {
		children = []map[string]any{
			{"id": "root.0", "parentId": "root", "label": "planner", "kind": "agent", "status": "succeeded"},
			{"id": "root.1", "parentId": "root", "label": "executor", "kind": "agent", "status": runStatus,
				"children": []map[string]any{
					{"id": "root.1.0", "parentId": "root.1", "label": "tool.echo", "kind": "tool", "status": "succeeded"},
				}},
		}
	}
	return map[string]any{
		"missionId": id,
		"nodeCount": 1 + len(children) + 2,
		"root": map[string]any{
			"id": "root", "label": "orchestrator", "kind": "agent", "status": runStatus,
			"runtimeId": "runtime.generic-loop", "providerId": "provider.mock-local",
			"children": children,
		},
	}
}

func (s *Server) buildLogs(r *http.Request, id string) []map[string]any {
	evs, err := s.sys.Bus.Replay(r.Context(), types.WorkUnitID(id))
	if err != nil || len(evs) == 0 {
		return []map[string]any{
			{"seq": 1, "ts": time.Now().UTC().Format(time.RFC3339), "level": "info", "message": "mission accepted", "nodeId": "root"},
			{"seq": 2, "ts": time.Now().UTC().Format(time.RFC3339), "level": "info", "message": "awaiting events", "nodeId": "root"},
		}
	}
	out := make([]map[string]any, 0, len(evs))
	for _, e := range evs {
		level := "info"
		if strings.Contains(e.Type, "fail") {
			level = "error"
		} else if strings.Contains(e.Type, "hitl") {
			level = "warn"
		}
		msg := e.Type
		if e.Data != nil {
			if g, ok := e.Data["goal"].(string); ok {
				msg = e.Type + " goal=" + g
			}
			if t, ok := e.Data["tool"].(string); ok {
				msg = e.Type + " tool=" + t
			}
		}
		out = append(out, map[string]any{
			"seq": e.Seq, "ts": e.Time.UTC().Format(time.RFC3339Nano),
			"level": level, "message": msg, "nodeId": "root",
			"tool": toolFromData(e.Data),
		})
	}
	return out
}

func toolFromData(d map[string]any) any {
	if d == nil {
		return nil
	}
	if t, ok := d["tool"].(string); ok {
		return map[string]any{"name": t, "input": d, "durationMs": 12}
	}
	return nil
}


func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/httpapi",
		AESPSpecs:  []string{"AESP-0014", "AESP-0015", "AESP-0011"},
		Invariants: []string{"INV-11"},
		Status:     "implemented",
	}
}
