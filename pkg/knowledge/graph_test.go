package knowledge

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestUpsertQuery(t *testing.T) {
	g := New()
	_ = g.Upsert(context.Background(), Triple{Subject: "svc", Predicate: "depends_on", Object: "db", Trust: types.TrustVerified})
	got, _ := g.Query(context.Background(), "svc")
	if len(got) != 1 {
		t.Fatal(got)
	}
}
