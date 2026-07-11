package builtin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/knowledge"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/memory"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/toolrouter"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// RegisterAll registers standard tools used by the reference Agent OS.
func RegisterAll(r *toolrouter.Router, mem *memory.Memory, kg *knowledge.Graph, workspace string) {
	if workspace == "" {
		workspace = "."
	}
	r.Register(&echoTool{})
	r.Register(&memoryWriteTool{mem: mem})
	r.Register(&memoryReadTool{mem: mem})
	r.Register(&kgUpsertTool{kg: kg})
	r.Register(&workspaceWriteTool{root: workspace})
	r.Register(&workspaceReadTool{root: workspace})
	r.Register(&failTool{})
}

type echoTool struct{}

func (echoTool) Name() string { return "echo" }
func (echoTool) Spec() toolrouter.Spec {
	return toolrouter.Spec{Name: "echo", Description: "Echo a message", SideEffectClass: types.SideEffectRead}
}
func (echoTool) Invoke(ctx context.Context, args map[string]any) (toolrouter.Result, error) {
	return toolrouter.Result{Output: args["msg"], Trust: types.TrustVerified}, nil
}

type memoryWriteTool struct{ mem *memory.Memory }

func (t *memoryWriteTool) Name() string { return "memory.write" }
func (t *memoryWriteTool) Spec() toolrouter.Spec {
	return toolrouter.Spec{Name: "memory.write", Description: "Write unified memory with trust label", SideEffectClass: types.SideEffectWriteWorkspace}
}
func (t *memoryWriteTool) Invoke(ctx context.Context, args map[string]any) (toolrouter.Result, error) {
	text, _ := args["text"].(string)
	trustS, _ := args["trust"].(string)
	trust := types.TrustLabel(trustS)
	if trust == "" {
		trust = types.TrustAgent
	}
	tenant, _ := args["tenant"].(string)
	wu, _ := args["workUnitId"].(string)
	err := t.mem.Write(ctx, memory.Item{
		ID: fmt.Sprintf("mem_%v", args["id"]), Text: text, Trust: trust,
		Tenant: types.TenantID(tenant), WorkUnit: types.WorkUnitID(wu), Scope: "session",
	})
	if err != nil {
		return toolrouter.Result{}, err
	}
	return toolrouter.Result{Output: map[string]any{"written": true, "trust": string(trust)}, Trust: types.TrustSystem}, nil
}

type memoryReadTool struct{ mem *memory.Memory }

func (t *memoryReadTool) Name() string { return "memory.read" }
func (t *memoryReadTool) Spec() toolrouter.Spec {
	return toolrouter.Spec{Name: "memory.read", Description: "Query unified memory", SideEffectClass: types.SideEffectRead}
}
func (t *memoryReadTool) Invoke(ctx context.Context, args map[string]any) (toolrouter.Result, error) {
	tenant, _ := args["tenant"].(string)
	items, err := t.mem.Query(ctx, memory.Query{Tenant: types.TenantID(tenant), Limit: 20})
	if err != nil {
		return toolrouter.Result{}, err
	}
	return toolrouter.Result{Output: items, Trust: types.TrustSystem}, nil
}

type kgUpsertTool struct{ kg *knowledge.Graph }

func (t *kgUpsertTool) Name() string { return "kg.upsert" }
func (t *kgUpsertTool) Spec() toolrouter.Spec {
	return toolrouter.Spec{Name: "kg.upsert", Description: "Upsert knowledge triple", SideEffectClass: types.SideEffectWriteWorkspace}
}
func (t *kgUpsertTool) Invoke(ctx context.Context, args map[string]any) (toolrouter.Result, error) {
	s, _ := args["subject"].(string)
	p, _ := args["predicate"].(string)
	o, _ := args["object"].(string)
	trustS, _ := args["trust"].(string)
	trust := types.TrustLabel(trustS)
	if trust == "" {
		trust = types.TrustAgent
	}
	err := t.kg.Upsert(ctx, knowledge.Triple{Subject: s, Predicate: p, Object: o, Trust: trust})
	if err != nil {
		return toolrouter.Result{}, err
	}
	return toolrouter.Result{Output: map[string]any{"upserted": true}, Trust: types.TrustSystem}, nil
}

type workspaceWriteTool struct {
	root string
	mu   sync.Mutex
}

func (t *workspaceWriteTool) Name() string { return "workspace.write" }
func (t *workspaceWriteTool) Spec() toolrouter.Spec {
	return toolrouter.Spec{Name: "workspace.write", Description: "Write file under workspace", SideEffectClass: types.SideEffectWriteWorkspace}
}
func (t *workspaceWriteTool) Invoke(ctx context.Context, args map[string]any) (toolrouter.Result, error) {
	rel, _ := args["path"].(string)
	content, _ := args["content"].(string)
	if rel == "" {
		return toolrouter.Result{}, fmt.Errorf("path required")
	}
	path := filepath.Join(t.root, filepath.Clean("/"+rel))
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return toolrouter.Result{}, err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return toolrouter.Result{}, err
	}
	return toolrouter.Result{Output: map[string]any{"path": rel}, Trust: types.TrustAgent}, nil
}

type workspaceReadTool struct{ root string }

func (t *workspaceReadTool) Name() string { return "workspace.read" }
func (t *workspaceReadTool) Spec() toolrouter.Spec {
	return toolrouter.Spec{Name: "workspace.read", Description: "Read file under workspace", SideEffectClass: types.SideEffectRead}
}
func (t *workspaceReadTool) Invoke(ctx context.Context, args map[string]any) (toolrouter.Result, error) {
	rel, _ := args["path"].(string)
	path := filepath.Join(t.root, filepath.Clean("/"+rel))
	b, err := os.ReadFile(path)
	if err != nil {
		return toolrouter.Result{}, err
	}
	return toolrouter.Result{Output: string(b), Trust: types.TrustRetrieved}, nil
}

// failTool always fails — used for retry/rollback demos.
type failTool struct {
	n int
	mu sync.Mutex
}

func (t *failTool) Name() string { return "flaky.step" }
func (t *failTool) Spec() toolrouter.Spec {
	return toolrouter.Spec{Name: "flaky.step", Description: "Fails first N calls then succeeds", SideEffectClass: types.SideEffectWriteWorkspace}
}
func (t *failTool) Invoke(ctx context.Context, args map[string]any) (toolrouter.Result, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.n++
	failTimes := 1
	if v, ok := args["failTimes"].(float64); ok {
		failTimes = int(v)
	}
	if v, ok := args["failTimes"].(int); ok {
		failTimes = v
	}
	if t.n <= failTimes {
		return toolrouter.Result{}, fmt.Errorf("flaky failure attempt %d", t.n)
	}
	return toolrouter.Result{Output: "recovered", Trust: types.TrustVerified}, nil
}
