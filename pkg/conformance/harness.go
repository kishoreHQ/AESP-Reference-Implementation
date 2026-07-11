package conformance

import (
	"fmt"
	"sort"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/a2a"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/agentos"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/agentregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/approval"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/artifact"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/capability"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/credentials"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/deploy"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/docgen"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/evaluation"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/executor"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/host"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/kernel"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/knowledge"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/mcp"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/memory"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/missionload"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/orchestrator"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/planner"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/policy"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/remediation"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/replay"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/reviewer"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/router"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/toolrouter"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// MUST is a tracked AESP requirement for conformance enumeration.
type MUST struct {
	ID     string
	Spec   string
	Title  string
	Status string // implemented | stubbed | missing | gap-filed
	Module string
}

// Catalog returns the reference implementation mapping of high-priority MUSTs.
func Catalog() []MUST {
	return []MUST{
		{ID: "CORE-WORKUNIT", Spec: "AESP-0001", Title: "WorkUnit identity and mission admission", Status: "implemented", Module: "pkg/kernel"},
		{ID: "CORE-ROLES", Spec: "AESP-0002", Title: "Agent principal registry and roles", Status: "implemented", Module: "pkg/agentregistry"},
		{ID: "CORE-EVENTS", Spec: "AESP-0003", Title: "Event envelope and bus", Status: "implemented", Module: "pkg/eventbus"},
		{ID: "MEM-UNIFIED", Spec: "AESP-0004", Title: "Unified memory with trust labels", Status: "implemented", Module: "pkg/memory"},
		{ID: "WF-ORCH", Spec: "AESP-0005", Title: "Orchestrated mission execution loop", Status: "implemented", Module: "pkg/agentos"},
		{ID: "KG-GRAPH", Spec: "AESP-0006", Title: "Knowledge graph upsert/query", Status: "implemented", Module: "pkg/knowledge"},
		{ID: "CG-ARTIFACT", Spec: "AESP-0007", Title: "Content-addressed artifacts", Status: "implemented", Module: "pkg/artifact"},
		{ID: "DOC-GEN", Spec: "AESP-0008", Title: "Documentation generator pipeline", Status: "implemented", Module: "pkg/docgen"},
		{ID: "DEP-ROLLOUT", Spec: "AESP-0009", Title: "Deployment session orchestration", Status: "implemented", Module: "pkg/deploy"},
		{ID: "TEST-EVAL", Spec: "AESP-0010", Title: "Evaluation harness distinct from agent harness", Status: "implemented", Module: "pkg/evaluation"},
		{ID: "OBS-EVENTS", Spec: "AESP-0011", Title: "Mission event journal / replay", Status: "implemented", Module: "pkg/replay"},
		{ID: "REM-PLAYBOOK", Spec: "AESP-0012", Title: "Remediation playbook engine", Status: "implemented", Module: "pkg/remediation"},
		{ID: "SEC-POLICY", Spec: "AESP-0013", Title: "Policy fail-closed for untrusted privilege", Status: "implemented", Module: "pkg/policy"},
		{ID: "HITL-NO-AUTO", Spec: "AESP-0014", Title: "HITL timeout must not auto-approve", Status: "implemented", Module: "pkg/approval"},
		{ID: "INT-PROVIDER", Spec: "AESP-0015", Title: "Provider registry capability advertisement", Status: "implemented", Module: "pkg/providerregistry"},
		{ID: "INT-RUNTIME", Spec: "AESP-0015", Title: "Runtime registry runtime.yaml discovery", Status: "implemented", Module: "pkg/runtimeregistry"},
		{ID: "INT-TOOLS", Spec: "AESP-0015", Title: "Unified tool router + invocation records", Status: "implemented", Module: "pkg/toolrouter"},
		{ID: "INT-PLAN", Spec: "AESP-0015", Title: "Versioned plan artifacts", Status: "implemented", Module: "pkg/planner"},
		{ID: "INT-MCP", Spec: "AESP-0015", Title: "MCP-aligned tool server/client + fixtures", Status: "implemented", Module: "pkg/mcp"},
		{ID: "INT-A2A", Spec: "AESP-0015", Title: "A2A peer registry + task fixtures", Status: "implemented", Module: "pkg/a2a"},
		{ID: "INV-01", Spec: "INV", Title: "Provider ≠ Runtime separate registries", Status: "implemented", Module: "pkg/providerregistry+pkg/runtimeregistry"},
		{ID: "INV-03", Spec: "INV", Title: "Capability-based routing + failover", Status: "implemented", Module: "pkg/router"},
		{ID: "INV-04", Spec: "INV", Title: "Unified memory", Status: "implemented", Module: "pkg/memory"},
		{ID: "INV-05", Spec: "INV", Title: "Context envelope", Status: "implemented", Module: "pkg/contextenv+pkg/agentos"},
		{ID: "INV-07", Spec: "INV", Title: "Credential broker handles", Status: "implemented", Module: "pkg/credentials"},
		{ID: "INV-09", Spec: "INV", Title: "Dynamic runtime registry", Status: "implemented", Module: "pkg/runtimeregistry"},
		{ID: "INV-10", Spec: "INV", Title: "Execution tree / replay journal", Status: "implemented", Module: "pkg/replay+pkg/kernel"},
		{ID: "INV-11", Spec: "INV", Title: "Host interface (in-process + HTTP)", Status: "implemented", Module: "pkg/host+pkg/httpapi"},
	}
}

// ModuleMappings collects SpecMapping from packages.
func ModuleMappings() []types.SpecMapping {
	return []types.SpecMapping{
		kernel.SpecMapping(),
		host.SpecMapping(),
		eventbus.SpecMapping(),
		contextenv.SpecMapping(),
		providerregistry.SpecMapping(),
		runtimeregistry.SpecMapping(),
		capability.SpecMapping(),
		router.SpecMapping(),
		memory.SpecMapping(),
		knowledge.SpecMapping(),
		policy.SpecMapping(),
		approval.SpecMapping(),
		artifact.SpecMapping(),
		toolrouter.SpecMapping(),
		credentials.SpecMapping(),
		orchestrator.SpecMapping(),
		planner.SpecMapping(),
		executor.SpecMapping(),
		reviewer.SpecMapping(),
		agentregistry.SpecMapping(),
		evaluation.SpecMapping(),
		replay.SpecMapping(),
		deploy.SpecMapping(),
		remediation.SpecMapping(),
		docgen.SpecMapping(),
		mcp.SpecMapping(),
		a2a.SpecMapping(),
		missionload.SpecMapping(),
		agentos.SpecMapping(),
	}
}

// Report returns a human-readable conformance summary.
func Report() string {
	cat := Catalog()
	counts := map[string]int{}
	for _, m := range cat {
		counts[m.Status]++
	}
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	s := "AESP Reference Implementation — Conformance Enumeration\n"
	for _, k := range keys {
		s += fmt.Sprintf("  %s: %d\n", k, counts[k])
	}
	s += "\nItems:\n"
	for _, m := range cat {
		s += fmt.Sprintf("  [%s] %s (%s) module=%s — %s\n", m.Status, m.ID, m.Spec, m.Module, m.Title)
	}
	s += fmt.Sprintf("\nModule mappings: %d\n", len(ModuleMappings()))
	return s
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/conformance",
		AESPSpecs:  []string{"AESP-0000", "CONFORMANCE"},
		Invariants: []string{"INV-08"},
		Status:     "implemented",
	}
}
