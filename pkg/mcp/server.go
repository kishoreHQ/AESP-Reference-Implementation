package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/toolrouter"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Minimal MCP-aligned tool surface (AESP-0015 INT).
// Not a full wire protocol — in-process server with golden fixture parity.

type ToolDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type CallRequest struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type CallResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError"`
	Trust   types.TrustLabel `json:"trust"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Server exposes tools via MCP-shaped list/call API.
type Server struct {
	mu    sync.Mutex
	tools map[string]ToolDef
	router *toolrouter.Router
	workUnit types.WorkUnitID
}

func NewServer(r *toolrouter.Router) *Server {
	return &Server{tools: map[string]ToolDef{}, router: r}
}

func (s *Server) RegisterTool(def ToolDef) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[def.Name] = def
}

func (s *Server) SetWorkUnit(id types.WorkUnitID) { s.workUnit = id }

func (s *Server) ListTools(ctx context.Context) []ToolDef {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ToolDef, 0, len(s.tools))
	for _, t := range s.tools {
		out = append(out, t)
	}
	return out
}

func (s *Server) CallTool(ctx context.Context, req CallRequest) (CallResult, error) {
	s.mu.Lock()
	_, ok := s.tools[req.Name]
	s.mu.Unlock()
	if !ok {
		return CallResult{IsError: true, Content: []ContentBlock{{Type: "text", Text: "unknown tool"}}, Trust: types.TrustUntrusted}, fmt.Errorf("unknown tool %s", req.Name)
	}
	if s.router == nil {
		return CallResult{IsError: true, Content: []ContentBlock{{Type: "text", Text: "no router"}}, Trust: types.TrustUntrusted}, fmt.Errorf("no router")
	}
	res, err := s.router.Invoke(ctx, s.workUnit, req.Name, req.Arguments, types.TrustAgent)
	if err != nil {
		return CallResult{IsError: true, Content: []ContentBlock{{Type: "text", Text: err.Error()}}, Trust: types.TrustUntrusted}, err
	}
	b, _ := json.Marshal(res.Output)
	return CallResult{
		Content: []ContentBlock{{Type: "text", Text: string(b)}},
		Trust:   res.Trust,
	}, nil
}

// Client is an MCP-shaped client binding to a Server.
type Client struct {
	Server *Server
}

func (c *Client) Initialize(ctx context.Context) (map[string]any, error) {
	return map[string]any{
		"protocolVersion": "2024-11-05",
		"serverInfo":      map[string]any{"name": "aesp-ref-mcp", "version": "1.0.0"},
		"capabilities":    map[string]any{"tools": map[string]any{}},
	}, nil
}

func (c *Client) ListTools(ctx context.Context) ([]ToolDef, error) {
	return c.Server.ListTools(ctx), nil
}

func (c *Client) CallTool(ctx context.Context, name string, args map[string]any) (CallResult, error) {
	return c.Server.CallTool(ctx, CallRequest{Name: name, Arguments: args})
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/mcp",
		AESPSpecs:  []string{"AESP-0015", "AESP-0013"},
		Invariants: []string{"INV-06"},
		Status:     "implemented",
	}
}
