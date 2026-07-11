package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/board"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/connections"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/routines"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/sessions"
)

func (s *Server) registerDeckRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/connections/probe", s.apiConnProbe)
	mux.HandleFunc("/api/v1/connections", s.apiConnections)
	mux.HandleFunc("/api/v1/sessions", s.apiSessions)
	mux.HandleFunc("/api/v1/sessions/", s.apiSessionSub)
	mux.HandleFunc("/api/v1/boards", s.apiBoards)
	mux.HandleFunc("/api/v1/tasks", s.apiTasks)
	mux.HandleFunc("/api/v1/tasks/", s.apiTaskSub)
	mux.HandleFunc("/api/v1/routines", s.apiRoutines)
	mux.HandleFunc("/api/v1/routines/", s.apiRoutineSub)
	mux.HandleFunc("/api/v1/goals", s.apiGoals)
	mux.HandleFunc("/api/v1/journal", s.apiJournal)
	mux.HandleFunc("/api/v1/analytics/agents/", s.apiAnalyticsAgent)
}

func (s *Server) apiConnProbe(w http.ResponseWriter, r *http.Request) {
	if s.sys.Connections == nil {
		writeErr(w, 503, "unavailable", "connections service not ready", "Restart aespd")
		return
	}
	writeOK(w, s.sys.Connections.Probe(r.Context()))
}

func (s *Server) apiConnections(w http.ResponseWriter, r *http.Request) {
	if s.sys.Connections == nil {
		writeErr(w, 503, "unavailable", "connections service not ready", "Restart aespd")
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeOK(w, s.sys.Connections.List())
	case http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		var req connections.RegisterRequest
		if err := json.Unmarshal(body, &req); err != nil {
			writeErr(w, 400, "bad_request", err.Error(), "Send JSON body")
			return
		}
		c, err := s.sys.Connections.Register(r.Context(), req)
		if err != nil {
			// still return connection with error detail for UI remediation
			if c != nil {
				writeEnv(w, 422, c, map[string]any{"code": "handshake_failed", "message": err.Error(), "remediation": c.LastError})
				return
			}
			writeErr(w, 422, "handshake_failed", err.Error(), "Check PATH, version, credentials")
			return
		}
		writeOK(w, c)
	default:
		writeErr(w, 405, "method", "GET or POST", "")
	}
}

func (s *Server) apiSessions(w http.ResponseWriter, r *http.Request) {
	if s.sys.Sessions == nil {
		writeErr(w, 503, "unavailable", "sessions not ready", "")
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeOK(w, s.sys.Sessions.List())
	case http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		var req sessions.CreateRequest
		_ = json.Unmarshal(body, &req)
		sess, err := s.sys.Sessions.Create(r.Context(), req)
		if err != nil {
			writeErr(w, 500, "session_error", err.Error(), "Check runtime adapter / PATH")
			return
		}
		writeOK(w, sess)
	default:
		writeErr(w, 405, "method", "GET or POST", "")
	}
}

func (s *Server) apiSessionSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/sessions/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		writeErr(w, 400, "bad_request", "missing id", "")
		return
	}
	id := parts[0]
	if len(parts) == 1 {
		sess, ok := s.sys.Sessions.Get(id)
		if !ok {
			writeErr(w, 404, "not_found", "session not found", "")
			return
		}
		writeOK(w, sess)
		return
	}
	switch parts[1] {
	case "message":
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Text string `json:"text"`
		}
		_ = json.Unmarshal(body, &req)
		if err := s.sys.Sessions.Message(r.Context(), id, req.Text); err != nil {
			writeErr(w, 400, "message_failed", err.Error(), "")
			return
		}
		writeOK(w, map[string]any{"ok": true})
	case "stop":
		_ = s.sys.Sessions.Stop(r.Context(), id)
		writeOK(w, map[string]any{"ok": true, "status": "stopped"})
	default:
		writeErr(w, 404, "not_found", "use message|stop", "")
	}
}

func (s *Server) apiBoards(w http.ResponseWriter, r *http.Request) {
	if s.sys.Board == nil {
		writeErr(w, 503, "unavailable", "board not ready", "")
		return
	}
	writeOK(w, s.sys.Board.ListBoards())
}

func (s *Server) apiTasks(w http.ResponseWriter, r *http.Request) {
	if s.sys.Board == nil {
		writeErr(w, 503, "unavailable", "board not ready", "")
		return
	}
	switch r.Method {
	case http.MethodGet:
		boardID := r.URL.Query().Get("board")
		if boardID == "" {
			boardID = "board_default"
		}
		writeOK(w, s.sys.Board.ListTasks(boardID))
	case http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		var req struct {
			BoardID      string   `json:"boardId"`
			Title        string   `json:"title"`
			Capabilities []string `json:"capabilities"`
			Assignee     string   `json:"assignee"`
		}
		_ = json.Unmarshal(body, &req)
		if req.BoardID == "" {
			req.BoardID = "board_default"
		}
		t, err := s.sys.Board.CreateTask(req.BoardID, req.Title, req.Capabilities, req.Assignee)
		if err != nil {
			writeErr(w, 400, "bad_request", err.Error(), "")
			return
		}
		writeOK(w, t)
	default:
		writeErr(w, 405, "method", "GET or POST", "")
	}
}

func (s *Server) apiTaskSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	parts := strings.Split(path, "/")
	if len(parts) < 1 {
		writeErr(w, 400, "bad_request", "missing id", "")
		return
	}
	id := parts[0]
	if len(parts) >= 2 && parts[1] == "claim" {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			AgentID string   `json:"agentId"`
			Caps    []string `json:"capabilities"`
		}
		_ = json.Unmarshal(body, &req)
		if req.AgentID == "" {
			req.AgentID = "agent.default"
		}
		t, err := s.sys.Board.Claim(req.AgentID, req.Caps)
		if err != nil {
			writeErr(w, 404, "no_task", err.Error(), "Add queued tasks")
			return
		}
		writeOK(w, t)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var req struct {
		Column   string `json:"column"`
		Assignee string `json:"assignee"`
	}
	_ = json.Unmarshal(body, &req)
	t, err := s.sys.Board.MoveTask(id, board.Column(req.Column), req.Assignee)
	if err != nil {
		writeErr(w, 400, "bad_request", err.Error(), "")
		return
	}
	writeOK(w, t)
}

func (s *Server) apiRoutines(w http.ResponseWriter, r *http.Request) {
	if s.sys.Routines == nil {
		writeErr(w, 503, "unavailable", "routines not ready", "")
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeOK(w, s.sys.Routines.List())
	case http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		var req routines.CreateReq
		_ = json.Unmarshal(body, &req)
		rt, err := s.sys.Routines.Create(req)
		if err != nil {
			writeErr(w, 400, "bad_request", err.Error(), "Use valid cron e.g. @every 1h or 0 * * * *")
			return
		}
		writeOK(w, rt)
	default:
		writeErr(w, 405, "method", "GET or POST", "")
	}
}

func (s *Server) apiRoutineSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/routines/")
	parts := strings.Split(path, "/")
	id := parts[0]
	if len(parts) >= 2 {
		switch parts[1] {
		case "pause":
			rt, err := s.sys.Routines.Pause(id, true)
			if err != nil {
				writeErr(w, 404, "not_found", err.Error(), "")
				return
			}
			writeOK(w, rt)
		case "resume":
			rt, err := s.sys.Routines.Pause(id, false)
			if err != nil {
				writeErr(w, 404, "not_found", err.Error(), "")
				return
			}
			writeOK(w, rt)
		case "next":
			fires, err := s.sys.Routines.NextFires(id, 3)
			if err != nil {
				writeErr(w, 400, "bad_request", err.Error(), "")
				return
			}
			writeOK(w, fires)
		case "fire":
			// manual fire for demos — creates a lightweight mission id link
			mid := "mis_rtn_" + time.Now().UTC().Format("150405")
			rt, err := s.sys.Routines.Fire(id, mid, "succeeded")
			if err != nil {
				writeErr(w, 404, "not_found", err.Error(), "")
				return
			}
			writeOK(w, rt)
		default:
			writeErr(w, 404, "not_found", "pause|resume|next|fire", "")
		}
		return
	}
	writeErr(w, 404, "not_found", "missing action", "")
}

func (s *Server) apiGoals(w http.ResponseWriter, r *http.Request) {
	if s.sys.Goals == nil {
		writeErr(w, 503, "unavailable", "goals not ready", "")
		return
	}
	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Title string `json:"title"`
		}
		_ = json.Unmarshal(body, &req)
		writeOK(w, s.sys.Goals.CreateGoal(req.Title, time.Time{}))
		return
	}
	writeOK(w, s.sys.Goals.ListGoals())
}

func (s *Server) apiJournal(w http.ResponseWriter, r *http.Request) {
	if s.sys.Goals == nil {
		writeErr(w, 503, "unavailable", "journal not ready", "")
		return
	}
	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Text string `json:"text"`
		}
		_ = json.Unmarshal(body, &req)
		writeOK(w, s.sys.Goals.AddJournal(r.Context(), req.Text))
		return
	}
	writeOK(w, s.sys.Goals.ListJournal())
}

func (s *Server) apiAnalyticsAgent(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/analytics/agents/")
	if id == "" {
		id = "agent.default"
	}
	if s.sys.Analytics == nil {
		writeErr(w, 503, "unavailable", "analytics not ready", "")
		return
	}
	writeOK(w, s.sys.Analytics.Agent(r.Context(), id))
}
