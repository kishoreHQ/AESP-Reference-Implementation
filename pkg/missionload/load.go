package missionload

import (
	"fmt"
	"os"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
	"gopkg.in/yaml.v3"
)

type fileMission struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		ID              string `yaml:"id"`
		ProfilePortable bool   `yaml:"profilePortable"`
	} `yaml:"metadata"`
	Spec struct {
		Goal                  string   `yaml:"goal"`
		RequiredCapabilities  []string `yaml:"requiredCapabilities"`
		SuccessCriteria       []string `yaml:"successCriteria"`
		Constraints           []string `yaml:"constraints"`
		Budget                struct {
			MaxSteps   int     `yaml:"maxSteps"`
			MaxTokens  int64   `yaml:"maxTokens"`
			MaxCostUSD float64 `yaml:"maxCostUSD"`
			MaxWallSec int64   `yaml:"maxWallSec"`
		} `yaml:"budget"`
		Labels map[string]string `yaml:"labels"`
		// Scenario hints for reference runner demos
		Scenario string `yaml:"scenario"` // single|multi|codegen|approval|memory|kg|remediation|hitl|failover|rollback
	} `yaml:"spec"`
}

// Load reads a mission.yaml into types.Mission.
func Load(path string) (types.Mission, string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return types.Mission{}, "", err
	}
	var fm fileMission
	if err := yaml.Unmarshal(b, &fm); err != nil {
		return types.Mission{}, "", err
	}
	if fm.Metadata.ID == "" {
		return types.Mission{}, "", fmt.Errorf("metadata.id required")
	}
	if len(fm.Spec.RequiredCapabilities) == 0 {
		return types.Mission{}, "", fmt.Errorf("requiredCapabilities required (INV-03)")
	}
	caps := make([]types.Capability, 0, len(fm.Spec.RequiredCapabilities))
	for _, c := range fm.Spec.RequiredCapabilities {
		caps = append(caps, types.Capability(c))
	}
	// Detect vendor-like model routing anti-pattern
	for _, c := range caps {
		s := string(c)
		if s == "gpt-4" || s == "claude" || s == "gemini" {
			return types.Mission{}, "", fmt.Errorf("model-name routing forbidden: %s (INV-03)", s)
		}
	}
	m := types.Mission{
		ID: types.WorkUnitID(fm.Metadata.ID),
		Tenant: "default",
		Goal: fm.Spec.Goal,
		Constraints: fm.Spec.Constraints,
		SuccessCriteria: fm.Spec.SuccessCriteria,
		RequiredCaps: caps,
		Budget: types.Budget{
			MaxSteps: fm.Spec.Budget.MaxSteps,
			MaxTokens: fm.Spec.Budget.MaxTokens,
			MaxCostUSD: fm.Spec.Budget.MaxCostUSD,
			MaxWallSec: fm.Spec.Budget.MaxWallSec,
		},
		Labels: fm.Spec.Labels,
		CreatedAt: time.Now().UTC(),
	}
	if m.Budget.MaxSteps == 0 {
		m.Budget.MaxSteps = 20
	}
	if m.Labels == nil {
		m.Labels = map[string]string{}
	}
	if fm.Spec.Scenario != "" {
		m.Labels["scenario"] = fm.Spec.Scenario
	}
	return m, fm.Spec.Scenario, nil
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/missionload",
		AESPSpecs:  []string{"AESP-0001", "AESP-0015"},
		Invariants: []string{"INV-03", "INV-11"},
		Status:     "implemented",
	}
}
