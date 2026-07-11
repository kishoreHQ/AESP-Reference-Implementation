package artifact

import (
	"context"
	"testing"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func TestPutGet(t *testing.T) {
	s := New()
	d, err := s.Put(context.Background(), []byte("hello"), Meta{WorkUnit: "wu", Trust: types.TrustVerified})
	if err != nil {
		t.Fatal(err)
	}
	b, m, err := s.Get(context.Background(), d)
	if err != nil || string(b) != "hello" || m.Digest != d {
		t.Fatal(err, b, m)
	}
}
