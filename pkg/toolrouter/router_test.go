package toolrouter

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type echoTool struct{}
func (echoTool) Name() string { return "echo" }
func (echoTool) Spec() Spec {
	return Spec{Name: "echo", SideEffectClass: types.SideEffectRead}
}
func (echoTool) Invoke(ctx context.Context, args map[string]any) (Result, error) {
	return Result{Output: args["msg"], Trust: types.TrustVerified}, nil
}

func TestInvoke(t *testing.T) {
	r := New(nil)
	r.Register(echoTool{})
	out, err := r.Invoke(context.Background(), "wu", "echo", map[string]any{"msg": "hi"}, types.TrustAgent)
	if err != nil || out.Output != "hi" {
		t.Fatal(err, out)
	}
	if len(r.Records()) != 1 {
		t.Fatal("record missing")
	}
}
