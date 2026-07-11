package credentials

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Broker issues short-lived handles (INV-07). Raw secrets never logged.
type Broker struct {
	mu      sync.Mutex
	secrets map[string]string // requirement key -> secret (in-memory for local profile)
	handles map[string]handle
}

type handle struct {
	audience string
	key      string
	expires  time.Time
}

func New() *Broker {
	return &Broker{secrets: map[string]string{}, handles: map[string]handle{}}
}

func (b *Broker) PutSecret(requirementKey, secret string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.secrets[requirementKey] = secret
}

func (b *Broker) Issue(ctx context.Context, requirementKey, audience string, ttl time.Duration) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.secrets[requirementKey]; !ok {
		return "", fmt.Errorf("secret not configured for %s", requirementKey)
	}
	id := fmt.Sprintf("cred_%d", time.Now().UnixNano())
	b.handles[id] = handle{audience: audience, key: requirementKey, expires: time.Now().Add(ttl)}
	return id, nil
}

// Resolve returns the secret for a handle if audience matches. Callers MUST not log return value.
func (b *Broker) Resolve(ctx context.Context, handleID, audience string) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	h, ok := b.handles[handleID]
	if !ok || time.Now().After(h.expires) {
		return "", fmt.Errorf("invalid or expired handle")
	}
	if h.audience != audience {
		return "", fmt.Errorf("audience mismatch")
	}
	return b.secrets[h.key], nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/credentials",
		AESPSpecs:  []string{"AESP-0013", "AESP-0015"},
		Invariants: []string{"INV-07"},
		Status:     "stubbed",
	}
}
