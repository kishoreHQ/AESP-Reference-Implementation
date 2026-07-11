package connections

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/adapters"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/credentials"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Kind of connection candidate / registration.
type Kind string

const (
	KindProvider Kind = "provider"
	KindRuntime  Kind = "runtime"
)

type Candidate struct {
	ID          string `json:"id"`
	Kind        Kind   `json:"kind"`
	Name        string `json:"name"`
	Detected    bool   `json:"detected"`
	Version     string `json:"version,omitempty"`
	Detail      string `json:"detail,omitempty"`
	NeedsCred   bool   `json:"needsCredential"`
	Unsandboxed bool   `json:"unsandboxed,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

type Connection struct {
	ID           string    `json:"id"`
	Kind         Kind      `json:"kind"`
	PluginID     string    `json:"pluginId"`
	Name         string    `json:"name"`
	Status       string    `json:"status"` // connected|error|disconnected
	Version      string    `json:"version,omitempty"`
	Unsandboxed  bool      `json:"unsandboxed,omitempty"`
	Capabilities []string  `json:"capabilities,omitempty"`
	LastError    string    `json:"lastError,omitempty"`
	ConnectedAt  time.Time `json:"connectedAt"`
	CredentialID string    `json:"credentialId,omitempty"` // never the secret
}

// Service implements K1 probe + register + handshake.
type Service struct {
	mu    sync.Mutex
	conns map[string]*Connection
	creds *credentials.Broker
	provs *providerregistry.Registry
}

func New(creds *credentials.Broker, provs *providerregistry.Registry) *Service {
	return &Service{
		conns: map[string]*Connection{},
		creds: creds,
		provs: provs,
	}
}

func (s *Service) Probe(ctx context.Context) []Candidate {
	var out []Candidate
	// Runtime adapters
	for _, a := range adapters.Catalog() {
		ok, ver, detail := a.Probe(ctx)
		unsandboxed := a.ID() == "runtime.generic-pty"
		out = append(out, Candidate{
			ID: a.ID(), Kind: KindRuntime, Name: a.ID(), Detected: ok,
			Version: ver, Detail: detail, NeedsCred: false, Unsandboxed: unsandboxed,
			Capabilities: []string{"coding", "tools"},
		})
	}
	// Local providers (no vendor hardcode in logic — declarative endpoints)
	localEndpoints := []struct {
		id, name, url string
		caps          []string
	}{
		{"provider.ollama-local", "local-openai-compat (11434)", "http://127.0.0.1:11434/api/tags", []string{"coding", "local", "tools"}},
		{"provider.lmstudio-local", "local-openai-compat (1234)", "http://127.0.0.1:1234/v1/models", []string{"coding", "local"}},
	}
	client := &http.Client{Timeout: 800 * time.Millisecond}
	for _, ep := range localEndpoints {
		ok := false
		detail := "unreachable"
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, ep.url, nil)
		if resp, err := client.Do(req); err == nil {
			_ = resp.Body.Close()
			ok = resp.StatusCode < 500
			if ok {
				detail = "reachable"
			} else {
				detail = fmt.Sprintf("http %d", resp.StatusCode)
			}
		}
		out = append(out, Candidate{
			ID: ep.id, Kind: KindProvider, Name: ep.name, Detected: ok,
			Detail: detail, NeedsCred: false, Capabilities: ep.caps,
		})
	}
	// Cloud-style provider plugins already registered (capability advertise only)
	if s.provs != nil {
		for _, p := range s.provs.List() {
			d, err := p.Describe(ctx)
			if err != nil {
				continue
			}
			caps := make([]string, 0, len(d.Capabilities))
			for _, c := range d.Capabilities {
				caps = append(caps, string(c))
			}
			out = append(out, Candidate{
				ID: string(d.ID), Kind: KindProvider, Name: string(d.ID),
				Detected: true, Detail: "registered plugin", NeedsCred: !d.Local,
				Capabilities: caps,
			})
		}
	}
	return out
}

type RegisterRequest struct {
	Kind         Kind     `json:"kind"`
	PluginID     string   `json:"pluginId"`
	Name         string   `json:"name"`
	Credential   string   `json:"credential,omitempty"` // write-once secret; never stored in connection record
	Capabilities []string `json:"capabilities,omitempty"`
	Command      []string `json:"command,omitempty"`
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*Connection, error) {
	if req.PluginID == "" {
		return nil, fmt.Errorf("pluginId required")
	}
	name := req.Name
	if name == "" {
		name = req.PluginID
	}
	credID := ""
	if req.Credential != "" && s.creds != nil {
		key := "conn." + req.PluginID
		s.creds.PutSecret(key, req.Credential)
		id, err := s.creds.Issue(ctx, key, req.PluginID, 24*time.Hour)
		if err == nil {
			credID = id
		}
	}

	conn := &Connection{
		ID: fmt.Sprintf("conn_%d", time.Now().UnixNano()), Kind: req.Kind,
		PluginID: req.PluginID, Name: name, Status: "disconnected",
		Capabilities: req.Capabilities, ConnectedAt: time.Now().UTC(),
		CredentialID: credID,
	}

	// Handshake
	switch req.Kind {
	case KindRuntime:
		var ad adapters.RuntimeAdapter
		for _, a := range adapters.Catalog() {
			if a.ID() == req.PluginID {
				ad = a
				break
			}
		}
		if ad == nil {
			ad = adapters.NewGenericPTY()
			conn.PluginID = ad.ID()
		}
		conn.Unsandboxed = ad.ID() == "runtime.generic-pty"
		ok, stdout, stderr, err := ad.Handshake(ctx)
		if !ok || err != nil {
			conn.Status = "error"
			conn.LastError = fmt.Sprintf("%v\nstdout:%s\nstderr:%s", err, stdout, stderr)
			s.store(conn)
			return conn, fmt.Errorf("handshake failed: %s", conn.LastError)
		}
		okp, ver, _ := ad.Probe(ctx)
		if okp {
			conn.Version = ver
		}
		conn.Status = "connected"
	case KindProvider:
		// 1-token style health via registry provider if present
		if s.provs != nil {
			if p, ok := s.provs.Get(types.PluginID(req.PluginID)); ok {
				if err := p.Health(ctx); err != nil {
					conn.Status = "error"
					conn.LastError = err.Error()
					s.store(conn)
					return conn, err
				}
				// capability discovery
				if d, err := p.Describe(ctx); err == nil {
					for _, c := range d.Capabilities {
						conn.Capabilities = append(conn.Capabilities, string(c))
					}
				}
				// 1-token completion smoke
				_, _ = p.Complete(ctx, providerregistry.CompletionRequest{
					Messages:  []providerregistry.Message{{Role: "user", Content: "ping"}},
					MaxTokens: 1,
				})
			}
		}
		conn.Status = "connected"
	default:
		return nil, fmt.Errorf("unknown kind")
	}
	s.store(conn)
	return conn, nil
}

func (s *Service) store(c *Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns[c.ID] = c
}

func (s *Service) List() []Connection {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Connection, 0, len(s.conns))
	for _, c := range s.conns {
		cp := *c
		out = append(out, cp)
	}
	return out
}

func (s *Service) Get(id string) (*Connection, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.conns[id]
	if !ok {
		return nil, false
	}
	cp := *c
	return &cp, true
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/connections",
		AESPSpecs:  []string{"AESP-0015", "AESP-0013"},
		Invariants: []string{"INV-01", "INV-02", "INV-07", "INV-09"},
		Status:     "implemented",
	}
}
