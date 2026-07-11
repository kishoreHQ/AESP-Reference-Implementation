package evaluation

import (
	"context"
	"testing"
)

func TestRun(t *testing.T) {
	h := New()
	res := h.Run(context.Background(), Campaign{ID: "c1", Cases: []Case{{ID: "t1", Expect: "ok"}}},
		func(ctx context.Context, c Case) (string, error) { return "ok", nil })
	if len(res) != 1 || !res[0].Passed {
		t.Fatal(res)
	}
}
