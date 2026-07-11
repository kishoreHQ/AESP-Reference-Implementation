package artifact

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

type Meta struct {
	Digest    types.ArtifactDigest
	WorkUnit  types.WorkUnitID
	Producer  string
	Trust     types.TrustLabel
	MediaType string
}

type Store struct {
	mu   sync.Mutex
	blob map[types.ArtifactDigest][]byte
	meta map[types.ArtifactDigest]Meta
}

func New() *Store {
	return &Store{blob: map[types.ArtifactDigest][]byte{}, meta: map[types.ArtifactDigest]Meta{}}
}

func DigestOf(b []byte) types.ArtifactDigest {
	h := sha256.Sum256(b)
	return types.ArtifactDigest("sha256:" + hex.EncodeToString(h[:]))
}

func (s *Store) Put(ctx context.Context, b []byte, m Meta) (types.ArtifactDigest, error) {
	d := DigestOf(b)
	m.Digest = d
	if m.Trust == "" {
		m.Trust = types.TrustAgent
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blob[d] = append([]byte{}, b...)
	s.meta[d] = m
	return d, nil
}

func (s *Store) Get(ctx context.Context, d types.ArtifactDigest) ([]byte, Meta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.blob[d]
	if !ok {
		return nil, Meta{}, fmt.Errorf("not found")
	}
	return append([]byte{}, b...), s.meta[d], nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/artifact",
		AESPSpecs:  []string{"AESP-0007", "AESP-0009", "AESP-0010"},
		Invariants: []string{"INV-10"},
		Status:     "stubbed",
	}
}
