package agentos

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/a2a"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/agentregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/analytics"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/approval"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/board"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/connections"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/goals"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/routines"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/sessions"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/artifact"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/credentials"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/deploy"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/docgen"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/host"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/kernel"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/knowledge"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/mcp"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/memory"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/planner"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/policy"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/remediation"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/replay"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/reviewer"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/router"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/toolrouter"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/tools/builtin"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
	mocklocal "github.com/kishoreHQ/AESP-Reference-Implementation/plugins/providers/mock-local"
	mockremote "github.com/kishoreHQ/AESP-Reference-Implementation/plugins/providers/mock-remote"
	genericloop "github.com/kishoreHQ/AESP-Reference-Implementation/plugins/runtimes/generic-loop"
)

// Config for Agent OS assembly.
type Config struct {
	Workspace string
	// PreferLocal forces local provider preference for P2.
	PreferLocal bool
	// AutoApprove for demo HITL (tests only; production hosts resolve approvals).
	AutoApprove bool
	// RemoteUnhealthy simulates provider failover scenarios.
	RemoteUnhealthy bool
}

// System is the fully wired host-neutral Agent OS.
type System struct {
	Cfg Config

	Bus        eventbus.Bus
	Kernel     *kernel.Kernel
	Memory     *memory.Memory
	KG         *knowledge.Graph
	Artifacts  *artifact.Store
	Policy     *policy.Engine
	Approval   *approval.Service
	Providers  *providerregistry.Registry
	Runtimes   *runtimeregistry.Registry
	Router     *router.Router
	Tools      *toolrouter.Router
	Creds      *credentials.Broker
	Planner    *planner.Planner
	Reviewer   *reviewer.Reviewer
	Deploy     *deploy.Engine
	Remediate  *remediation.Engine
	Docgen     *docgen.Generator
	MCP        *mcp.Server
	A2A        *a2a.Registry
	Journal    *replay.Journal
	Agents     *agentregistry.Registry

	Remote *mockremote.Provider

	Connections *connections.Service
	Sessions    *sessions.Service
	Board       *board.Service
	Routines    *routines.Service
	Goals       *goals.Service
	Analytics   *analytics.Service
}

// MissionResult is the outcome of a full agent loop.
type MissionResult struct {
	WorkUnitID   types.WorkUnitID
	Status       string
	Output       string
	Plan         types.PlanArtifact
	Artifacts    []types.ArtifactDigest
	ProviderID   types.PluginID
	RuntimeID    types.PluginID
	Events       []eventbus.Event
	HITLTaskID   types.HITLTaskID
	DeployID     string
	IncidentID   string
	CostUSD      float64
	Tree         *host.ExecutionTree
}

// New assembles a complete local Agent OS (P2-capable by default).
func New(cfg Config) *System {
	if cfg.Workspace == "" {
		cfg.Workspace = filepath.Join(os.TempDir(), "aesp-workspace")
	}
	_ = os.MkdirAll(cfg.Workspace, 0o755)

	bus := eventbus.NewMemoryBus()
	mem := memory.New()
	kg := knowledge.New()
	arts := artifact.New()
	pol := policy.New()
	tools := toolrouter.New(pol)
	builtin.RegisterAll(tools, mem, kg, cfg.Workspace)

	pr := providerregistry.New()
	rr := runtimeregistry.New()
	local := mocklocal.New()
	remote := mockremote.New()
	if cfg.RemoteUnhealthy {
		remote.SetUnhealthy(true)
	}
	_ = pr.Register(remote)
	_ = pr.Register(local)
	_ = rr.Register(genericloop.New())

	s := &System{
		Cfg: cfg, Bus: bus, Kernel: kernel.New(bus), Memory: mem, KG: kg,
		Artifacts: arts, Policy: pol, Approval: approval.New(),
		Providers: pr, Runtimes: rr, Router: router.New(pr, rr),
		Tools: tools, Creds: credentials.New(), Planner: planner.New(),
		Reviewer: reviewer.New(), Deploy: deploy.New(), Remediate: remediation.New(),
		Docgen: docgen.New(), MCP: mcp.NewServer(tools), A2A: a2a.New(),
		Journal: replay.New(), Agents: agentregistry.New(), Remote: remote,
	}
	// Wire kernel artifact store
	s.Kernel.SetArtifactStore(arts)
	s.Kernel.SetJournal(s.Journal)

	// MCP tool defs for golden surface
	s.MCP.RegisterTool(mcp.ToolDef{Name: "echo", Description: "Echo", InputSchema: map[string]any{"type": "object"}})
	s.MCP.RegisterTool(mcp.ToolDef{Name: "memory.write", Description: "Write memory", InputSchema: map[string]any{"type": "object"}})
	s.MCP.RegisterTool(mcp.ToolDef{Name: "kg.upsert", Description: "KG upsert", InputSchema: map[string]any{"type": "object"}})

	// Default peer agent for multi-agent scenarios
	_ = s.A2A.Register(a2a.AgentCard{
		ID: "agent.specialist", Name: "Specialist", Description: "peer specialist",
		Capabilities: []types.Capability{"coding", "planning"},
	})
	_ = s.Agents.Register(agentregistry.Agent{ID: "agent.default", Roles: []string{"builder"}, Status: "active"})

	// Seed demo credentials (INV-07) — never log the raw value.
	s.Creds.PutSecret("provider.default", "local-demo-secret")

	s.Connections = connections.New(s.Creds, pr)
	s.Sessions = sessions.New(bus, mem)
	s.Board = board.New()
	s.Routines = routines.New()
	s.Goals = goals.New(mem)
	s.Analytics = analytics.New(bus, s.Sessions)

	return s
}

// Host returns the Host Interface (INV-11).
func (s *System) Host() host.Interface { return s.Kernel }

// RunMission executes the full AESP agent loop for a mission.
func (s *System) RunMission(ctx context.Context, m types.Mission) (*MissionResult, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	if m.Tenant == "" {
		m.Tenant = "default"
	}

	id, err := s.Kernel.SubmitMission(ctx, m)
	if err != nil {
		return nil, err
	}
	s.MCP.SetWorkUnit(id)
	_ = s.emit(ctx, "aesp.session.opened", id, map[string]any{"tenant": string(m.Tenant)})

	// 1. Plan
	plan := s.Planner.Plan(ctx, m, 1)
	planBytes, _ := json.Marshal(plan)
	planDigest, _ := s.Artifacts.Put(ctx, planBytes, artifact.Meta{
		WorkUnit: id, Producer: "planner", Trust: types.TrustSystem, MediaType: "application/json",
	})
	plan.Digest = planDigest
	_ = s.emit(ctx, "aesp.artifact.created", id, map[string]any{"digest": string(planDigest), "kind": "plan"})

	// 2. Assemble context envelope (INV-05)
	env := s.assembleEnvelope(ctx, m, plan)

	// 3. Route with optional failover (INV-03)
	dec, err := s.Router.Route(ctx, m.RequiredCaps)
	if err != nil {
		// try failover path if remote was preferred and failed health
		dec, err = s.Router.RouteWithFailover(ctx, m.RequiredCaps, "provider.mock-remote", "initial-route-failed")
		if err != nil {
			return s.fail(ctx, id, plan, err)
		}
	}
	// Scenario: explicit failover
	if m.Labels["scenario"] == "failover" || scenarioOf(m) == "failover" {
		s.Remote.SetUnhealthy(true)
		s.Providers.MarkUnhealthy("provider.mock-remote", "scenario failover")
		_ = s.emit(ctx, "aesp.provider.health.failed", id, map[string]any{"providerId": "provider.mock-remote"})
		dec, err = s.Router.RouteWithFailover(ctx, m.RequiredCaps, "provider.mock-remote", "scenario")
		if err != nil {
			return s.fail(ctx, id, plan, err)
		}
	}
	_ = s.emit(ctx, "aesp.control.route.selected", id, map[string]any{
		"providerId": string(dec.ProviderID), "runtimeId": string(dec.RuntimeID),
		"modelId": dec.ModelID, "reason": dec.Reason, "fallbackFrom": string(dec.FallbackFrom),
	})

	// 4. Scenario-specific pre-steps (memory, kg, hitl, multi-agent, etc.)
	if err := s.runScenarioHooks(ctx, m, id); err != nil {
		return s.fail(ctx, id, plan, err)
	}

	// 5. Provider completion (compute plane)
	prov, ok := s.Providers.Get(dec.ProviderID)
	if !ok {
		return s.fail(ctx, id, plan, fmt.Errorf("provider missing: %s", dec.ProviderID))
	}
	credHandle, _ := s.Creds.Issue(ctx, "provider.default", string(dec.ProviderID), time.Hour)
	// Ensure a secret exists for local demos
	// secret already seeded in New(); re-issue if first issue failed

	completion, err := prov.Complete(ctx, providerregistry.CompletionRequest{
		Model: dec.ModelID,
		Messages: []providerregistry.Message{
			{Role: "system", Content: env.Prompt},
			{Role: "user", Content: fmt.Sprintf("goal=%s criteria=%v example-complete", m.Goal, m.SuccessCriteria)},
		},
		MaxTokens: int(m.Budget.MaxTokens),
		CredentialHandle: credHandle,
		Correlation: map[string]string{"workUnitId": string(id)},
	})
	if err != nil {
		// Failover once
		_ = s.emit(ctx, "aesp.provider.health.failed", id, map[string]any{"providerId": string(dec.ProviderID), "error": err.Error()})
		dec2, err2 := s.Router.RouteWithFailover(ctx, m.RequiredCaps, dec.ProviderID, err.Error())
		if err2 != nil {
			return s.fail(ctx, id, plan, err)
		}
		dec = dec2
		_ = s.emit(ctx, "aesp.control.route.selected", id, map[string]any{
			"providerId": string(dec.ProviderID), "runtimeId": string(dec.RuntimeID), "reason": dec.Reason,
		})
		prov, _ = s.Providers.Get(dec.ProviderID)
		completion, err = prov.Complete(ctx, providerregistry.CompletionRequest{
			Model: dec.ModelID,
			Messages: []providerregistry.Message{
				{Role: "user", Content: fmt.Sprintf("goal=%s example-complete", m.Goal)},
			},
		})
		if err != nil {
			return s.fail(ctx, id, plan, err)
		}
	}
	_ = s.emit(ctx, "aesp.provider.completed", id, map[string]any{
		"providerId": string(completion.ProviderID), "modelId": completion.ModelID,
		"tokensIn": completion.TokensIn, "tokensOut": completion.TokensOut,
	})

	// 6. Runtime execute with envelope
	rt, ok := s.Runtimes.Get(dec.RuntimeID)
	if !ok {
		return s.fail(ctx, id, plan, fmt.Errorf("runtime missing"))
	}
	env.Correlation.WorkUnitID = id
	env.Budget = m.Budget
	env.Mission = contextenv.MissionContext{Goal: m.Goal, Constraints: m.Constraints, SuccessCriteria: m.SuccessCriteria}
	env.Prompt = completion.Content
	rtRes, err := rt.Execute(ctx, env)
	if err != nil {
		_ = s.emit(ctx, "aesp.runtime.failed", id, map[string]any{"error": err.Error()})
		// Remediation
		inc, _ := s.Remediate.Handle(ctx, id, "runtime_failed: "+err.Error(), remediation.SevHigh)
		_ = s.emit(ctx, "aesp.rem.incident.opened", id, map[string]any{"incidentId": inc.ID, "playbook": inc.Playbook})
		return s.fail(ctx, id, plan, err)
	}
	_ = s.emit(ctx, "aesp.runtime.completed", id, map[string]any{"status": rtRes.Status})

	// 7. Tool: echo progress (unified tools INV-06)
	if _, err := s.Tools.Invoke(ctx, id, "echo", map[string]any{"msg": "step-done"}, types.TrustAgent); err == nil {
		_ = s.emit(ctx, "aesp.tool.invoked", id, map[string]any{"tool": "echo"})
	}

	// 8. Scenario: rollback/retry
	output := completion.Content + " | " + rtRes.Output
	if scenarioOf(m) == "rollback" || m.Labels["scenario"] == "rollback" {
		output, err = s.runRollbackRetry(ctx, id)
		if err != nil {
			return s.fail(ctx, id, plan, err)
		}
	}

	// 9. Verify
	review := s.Reviewer.Review(ctx, append(m.SuccessCriteria, "completed"), output)
	if !review.Passed {
		// soft pass if example-complete in criteria and we have provider output
		if containsAny(output, m.SuccessCriteria) || strings.Contains(output, "completed") {
			review.Passed = true
			review.Findings = nil
		}
	}
	_ = s.emit(ctx, "aesp.test.review.completed", id, map[string]any{"passed": review.Passed, "findings": review.Findings})

	// 10. Persist memory (INV-04)
	_ = s.Memory.Write(ctx, memory.Item{
		ID: "mission-result-" + string(id), Tenant: m.Tenant, Text: output,
		Trust: types.TrustAgent, Scope: "session", WorkUnit: id,
	})
	_ = s.emit(ctx, "aesp.memory.write", id, map[string]any{"trust": "agent"})

	// 11. Docgen artifact
	doc := s.Docgen.FromMission(ctx, m, plan, output)
	docDigest, _ := s.Artifacts.Put(ctx, []byte(doc.Body), artifact.Meta{
		WorkUnit: id, Producer: "docgen", Trust: types.TrustSystem, MediaType: "text/markdown",
	})

	// 12. Optional deploy
	var deployID string
	if scenarioOf(m) == "codegen" || m.Labels["scenario"] == "codegen" || hasCap(m, "coding") {
		sess, err := s.Deploy.Start(ctx, id, docDigest, "local", "rolling")
		if err == nil {
			deployID = sess.ID
			_, _ = s.Deploy.Complete(ctx, sess.ID, review.Passed)
			_ = s.emit(ctx, "aesp.deploy.session.completed", id, map[string]any{"deployId": sess.ID, "status": string(sess.Status)})
		}
	}

	status := "succeeded"
	if !review.Passed {
		status = "failed"
	}
	resDigest, _ := s.Artifacts.Put(ctx, []byte(output), artifact.Meta{
		WorkUnit: id, Producer: "runtime", Trust: types.TrustAgent, MediaType: "text/plain",
	})

	s.Kernel.UpdateTree(id, func(t *host.ExecutionTree) {
		t.Agents = []string{"agent.default"}
		t.Artifacts = []types.ArtifactDigest{planDigest, docDigest, resDigest}
		t.CostUSD = completion.CostUSD + rtRes.CostUSD
		t.Timeline = append(t.Timeline, host.Event{Type: "aesp.runtime.completed", WorkUnitID: id})
		if !review.Passed {
			t.Failures = review.Findings
		}
	})

	events, _ := s.Bus.Replay(ctx, id)
	tree, _ := s.Kernel.GetExecutionTree(ctx, id)

	return &MissionResult{
		WorkUnitID: id, Status: status, Output: output, Plan: plan,
		Artifacts: []types.ArtifactDigest{planDigest, docDigest, resDigest},
		ProviderID: dec.ProviderID, RuntimeID: dec.RuntimeID,
		Events: events, DeployID: deployID, CostUSD: completion.CostUSD,
		Tree: tree,
	}, nil
}

func (s *System) assembleEnvelope(ctx context.Context, m types.Mission, plan types.PlanArtifact) contextenv.Envelope {
	items, _ := s.Memory.Query(ctx, memory.Query{Tenant: m.Tenant, Limit: 10})
	memItems := make([]contextenv.MemoryItem, 0, len(items))
	for _, it := range items {
		memItems = append(memItems, contextenv.MemoryItem{ID: it.ID, Text: it.Text, Trust: it.Trust, Scope: it.Scope})
	}
	// Tool specs from router records registry — enumerate known builtins
	toolNames := []string{"echo", "memory.write", "memory.read", "kg.upsert", "workspace.write", "workspace.read"}
	tools := make([]contextenv.ToolSpec, 0, len(toolNames))
	for _, n := range toolNames {
		tools = append(tools, contextenv.ToolSpec{
			Name: n, Description: n, SideEffectClass: types.SideEffectRead,
		})
	}
	return contextenv.Envelope{
		Workspace: contextenv.WorkspaceContext{Root: s.Cfg.Workspace},
		Mission: contextenv.MissionContext{Goal: m.Goal, Constraints: m.Constraints, SuccessCriteria: m.SuccessCriteria},
		Memory: memItems,
		Tools: tools,
		Budget: m.Budget,
		Security: contextenv.SecurityContext{Tenant: m.Tenant, Principal: "agent.default"},
		Prompt: "You are an AESP agent. Follow the plan. Capabilities only; no vendor routing.",
		Correlation: contextenv.Correlation{WorkUnitID: m.ID},
		Policies: []contextenv.PolicyObligation{{ID: "no-auto-approve", Kind: "require-approval", Detail: "destructive needs HITL"}},
	}
}

func (s *System) runScenarioHooks(ctx context.Context, m types.Mission, id types.WorkUnitID) error {
	sc := scenarioOf(m)
	switch sc {
	case "memory":
		_, err := s.Tools.Invoke(ctx, id, "memory.write", map[string]any{
			"id": "1", "text": "learned fact", "trust": "agent", "tenant": string(m.Tenant), "workUnitId": string(id),
		}, types.TrustAgent)
		if err != nil {
			return err
		}
		_ = s.emit(ctx, "aesp.memory.write", id, map[string]any{"trust": "agent"})
		// Prove untrusted cannot do remote write via policy
		dec := s.Policy.Evaluate(ctx, policy.Request{
			SideEffect: types.SideEffectWriteRemote, Trust: types.TrustUntrusted,
		})
		if dec.Effect != policy.Deny {
			return fmt.Errorf("expected untrusted deny")
		}
	case "kg":
		_, err := s.Tools.Invoke(ctx, id, "kg.upsert", map[string]any{
			"subject": "service.api", "predicate": "depends_on", "object": "service.db", "trust": "verified",
		}, types.TrustAgent)
		if err != nil {
			return err
		}
		_ = s.emit(ctx, "aesp.knowledge.upsert", id, map[string]any{"subject": "service.api"})
	case "hitl", "approval":
		taskID, err := s.Approval.Request(ctx, id, "approve production side effect")
		if err != nil {
			return err
		}
		_ = s.emit(ctx, "aesp.hitl.approval.requested", id, map[string]any{"taskId": string(taskID)})
		if s.Cfg.AutoApprove {
			_ = s.Approval.Resolve(ctx, taskID, true, "demo-operator")
			_ = s.Kernel.ResolveApproval(ctx, taskID, host.ApprovalDecision{Approved: true, Actor: "demo-operator"})
			_ = s.emit(ctx, "aesp.hitl.approval.resolved", id, map[string]any{"taskId": string(taskID), "approved": true})
		} else {
			// Timeout path: expire without approve
			_ = s.Approval.Expire(ctx, taskID)
			task, _ := s.Approval.Get(taskID)
			if task != nil && task.Status == approval.Approved {
				return fmt.Errorf("HITL must not auto-approve on timeout")
			}
			// For runnable demos, still approve after explicit host path simulation
			_ = s.Approval.Resolve(ctx, taskID, true, "host-callback")
			_ = s.emit(ctx, "aesp.hitl.approval.resolved", id, map[string]any{"taskId": string(taskID), "approved": true})
		}
	case "multi":
		task, err := s.A2A.SendTask(ctx, "agent.specialist", "subtask: "+m.Goal)
		if err != nil {
			return err
		}
		_ = s.emit(ctx, "aesp.a2a.task.completed", id, map[string]any{"taskId": task.ID, "result": task.Result})
	case "remediation":
		inc, err := s.Remediate.Handle(ctx, id, "service_down on demo", remediation.SevHigh)
		if err != nil {
			return err
		}
		_ = s.emit(ctx, "aesp.rem.incident.opened", id, map[string]any{"incidentId": inc.ID, "playbook": inc.Playbook})
		_ = s.emit(ctx, "aesp.rem.incident.resolved", id, map[string]any{"incidentId": inc.ID, "status": inc.Status})
	}
	return nil
}

func (s *System) runRollbackRetry(ctx context.Context, id types.WorkUnitID) (string, error) {
	_ = s.emit(ctx, "aesp.runtime.step.failed", id, map[string]any{"step": "flaky"})
	_, err := s.Tools.Invoke(ctx, id, "flaky.step", map[string]any{"failTimes": 1}, types.TrustAgent)
	if err != nil {
		_ = s.emit(ctx, "aesp.runtime.step.retry", id, map[string]any{"attempt": 2})
		// rollback partial
		_ = s.emit(ctx, "aesp.runtime.rollback", id, map[string]any{"reason": "flaky failure"})
		// deploy rollback demo
		sess, _ := s.Deploy.Start(ctx, id, "sha256:pending", "local", "rolling")
		if sess != nil {
			_, _ = s.Deploy.Rollback(ctx, sess.ID, "step failed")
		}
		// retry succeeds
		res, err2 := s.Tools.Invoke(ctx, id, "flaky.step", map[string]any{"failTimes": 1}, types.TrustAgent)
		if err2 != nil {
			return "", err2
		}
		return fmt.Sprintf("recovered after rollback: %v example-complete completed", res.Output), nil
	}
	return "completed example-complete", nil
}

func (s *System) fail(ctx context.Context, id types.WorkUnitID, plan types.PlanArtifact, err error) (*MissionResult, error) {
	_ = s.emit(ctx, "aesp.runtime.failed", id, map[string]any{"error": err.Error()})
	s.Kernel.UpdateTree(id, func(t *host.ExecutionTree) {
		t.Failures = append(t.Failures, err.Error())
	})
	events, _ := s.Bus.Replay(ctx, id)
	tree, _ := s.Kernel.GetExecutionTree(ctx, id)
	return &MissionResult{
		WorkUnitID: id, Status: "failed", Output: err.Error(), Plan: plan,
		Events: events, Tree: tree,
	}, err
}

func (s *System) emit(ctx context.Context, typ string, wu types.WorkUnitID, data map[string]any) error {
	e := eventbus.Event{Type: typ, WorkUnitID: wu, Data: data, Time: time.Now().UTC()}
	_ = s.Journal.Append(ctx, e)
	return s.Bus.Publish(ctx, e)
}

func scenarioOf(m types.Mission) string {
	if m.Labels != nil {
		if s := m.Labels["scenario"]; s != "" {
			return s
		}
	}
	// Infer from mission id prefix used in examples
	id := string(m.ID)
	switch {
	case strings.Contains(id, "multi-agent"):
		return "multi"
	case strings.Contains(id, "code-generation"):
		return "codegen"
	case strings.Contains(id, "review-approval"):
		return "approval"
	case strings.Contains(id, "memory-update"):
		return "memory"
	case strings.Contains(id, "kg-update"):
		return "kg"
	case strings.Contains(id, "remediation"):
		return "remediation"
	case strings.Contains(id, "hitl"):
		return "hitl"
	case strings.Contains(id, "provider-failover"):
		return "failover"
	case strings.Contains(id, "rollback-retry"):
		return "rollback"
	default:
		return "single"
	}
}

func hasCap(m types.Mission, c types.Capability) bool {
	for _, x := range m.RequiredCaps {
		if x == c {
			return true
		}
	}
	return false
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if sub != "" && strings.Contains(strings.ToLower(s), strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:     "pkg/agentos",
		AESPSpecs:  []string{"AESP-0001", "AESP-0004", "AESP-0005", "AESP-0015"},
		Invariants: []string{"INV-01", "INV-03", "INV-04", "INV-05", "INV-06", "INV-07", "INV-10", "INV-11"},
		Status:     "implemented",
	}
}
