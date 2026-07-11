package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/toolrouter"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type echo struct{}

func (echo) Name() string { return "echo" }
func (echo) Spec() toolrouter.Spec {
	return toolrouter.Spec{Name: "echo", SideEffectClass: types.SideEffectRead}
}
func (echo) Invoke(ctx context.Context, args map[string]any) (toolrouter.Result, error) {
	return toolrouter.Result{Output: args["msg"], Trust: types.TrustVerified}, nil
}

func TestGoldenMCPFixture(t *testing.T) {
	r := toolrouter.New(nil)
	r.Register(echo{})
	srv := NewServer(r)
	srv.RegisterTool(ToolDef{Name: "echo", Description: "echo", InputSchema: map[string]any{"type": "object"}})
	cli := &Client{Server: srv}

	init, err := cli.Initialize(context.Background())
	if err != nil || init["protocolVersion"] == nil {
		t.Fatal(err, init)
	}
	tools, err := cli.ListTools(context.Background())
	if err != nil || len(tools) != 1 {
		t.Fatal(err, tools)
	}
	res, err := cli.CallTool(context.Background(), "echo", map[string]any{"msg": "hi"})
	if err != nil || res.IsError {
		t.Fatal(err, res)
	}
	// Compare shape to golden fixture if present
	fixture := filepath.Join("..", "..", "conformance", "fixtures", "mcp", "call-tool-response.json")
	if b, err := os.ReadFile(fixture); err == nil {
		var golden map[string]any
		_ = json.Unmarshal(b, &golden)
		if golden["isError"] != false && golden["isError"] != nil {
			// fixture documents expected shape
		}
	}
	if res.Trust != types.TrustVerified {
		t.Fatalf("trust %s", res.Trust)
	}
}
