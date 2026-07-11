package credentials

import (
	"context"
	"testing"
	"time"
)

func TestIssueResolve(t *testing.T) {
	b := New()
	b.PutSecret("provider.default", "s3cret")
	h, err := b.Issue(context.Background(), "provider.default", "prov.a", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Resolve(context.Background(), h, "prov.a")
	if err != nil || s != "s3cret" {
		t.Fatal(err, s)
	}
	_, err = b.Resolve(context.Background(), h, "other")
	if err == nil {
		t.Fatal("audience mismatch expected")
	}
}
