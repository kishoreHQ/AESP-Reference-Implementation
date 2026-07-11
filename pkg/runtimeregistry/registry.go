package runtimeregistry

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
	"gopkg.in/yaml.v3"
)

// Runtime is an agent harness plugin (INV-01, INV-09).
type Runtime interface {
	ID() types.PluginID
	Describe(ctx context.Context) (Descriptor, error)
	Execute(ctx context.Context, env contextenv.Envelope) (Result, error)
	Health(ctx context.Context) error
}

type Descriptor struct {
	ID            types.PluginID     `json:"id"`
	Version       string             `json:"version"`
	CapabilitiesIn  []types.Capability `json:"capabilitiesIn"`
	CapabilitiesOut []types.Capability `json:"capabilitiesOut"`
	Sandbox       string             `json:"sandbox"`
}

type Result struct {
	Status     string                 `json:"status"` // succeeded|failed|needs_approval|budget
	Output     string                 `json:"output,omitempty"`
	Artifacts  []types.ArtifactDigest `json:"artifacts,omitempty"`
	StepsUsed  int                    `json:"stepsUsed"`
	TokensUsed int64                  `json:"tokensUsed"`
	CostUSD    float64                `json:"costUSD"`
}

// Manifest is runtime.yaml (INV-09).
type Manifest struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		ID      string `yaml:"id"`
		Version string `yaml:"version"`
	} `yaml:"metadata"`
	Spec struct {
		CapabilitiesIn  []string `yaml:"capabilitiesIn"`
		CapabilitiesOut []string `yaml:"capabilitiesOut"`
		Sandbox         string   `yaml:"sandbox"`
		Entrypoint      string   `yaml:"entrypoint"`
	} `yaml:"spec"`
}

type Registry struct {
	mu   sync.RWMutex
	byID map[types.PluginID]Runtime
}

func New() *Registry {
	return &Registry{byID: make(map[types.PluginID]Runtime)}
}

func (r *Registry) Register(rt Runtime) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rt.ID() == "" {
		return fmt.Errorf("runtime id required")
	}
	r.byID[rt.ID()] = rt
	return nil
}

func (r *Registry) Get(id types.PluginID) (Runtime, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rt, ok := r.byID[id]
	return rt, ok
}

func (r *Registry) List() []Runtime {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Runtime, 0, len(r.byID))
	for _, rt := range r.byID {
		out = append(out, rt)
	}
	return out
}

// LoadManifest parses runtime.yaml without loading code (discovery).
func LoadManifest(path string) (*Manifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if m.Metadata.ID == "" {
		return nil, fmt.Errorf("metadata.id required")
	}
	return &m, nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/runtimeregistry",
		AESPSpecs:  []string{"AESP-0001", "AESP-0005", "AESP-0015"},
		Invariants: []string{"INV-01", "INV-02", "INV-09"},
		Status:     "implemented",
	}
}
