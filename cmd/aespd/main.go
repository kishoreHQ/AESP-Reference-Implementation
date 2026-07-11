package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/agentos"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/conformance"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/httpapi"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/missionload"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}
	switch os.Args[1] {
	case "conformance":
		fmt.Print(conformance.Report())
	case "run":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: aespd run <mission.yaml>")
			os.Exit(2)
		}
		if err := cmdRun(os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "run-example":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: aespd run-example <01-single-agent|...>")
			os.Exit(2)
		}
		if err := cmdRunExample(os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "run-all-examples":
		if err := cmdRunAllExamples(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "serve":
		addr := ":8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		if err := cmdServe(addr); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "demo":
		if err := cmdDemo(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "config":
		cmdConfig()
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(2)
	}
}

func printHelp() {
	fmt.Print(`aespd — AESP Reference Agent OS

Commands:
  config                       Show memory / session / provider wiring
  demo                         Run built-in single-agent demo
  run <mission.yaml>           Execute a mission file
  run-example <name>           Run bundled example (e.g. 01-single-agent)
  run-all-examples             Run all 10 bundled examples
  serve [addr]                 HTTP Host Interface (default :8080)
  conformance                  Print AESP MUST catalog status
  help                         Show this help

Profiles: P1 via serve, P2 local CLI, P3 embed pkg/agentos.System

Env:
  AESP_WORKSPACE   Workspace root for file tools (default: temp or .aesp-workspace)
`)
}

func workspace() string {
	if w := os.Getenv("AESP_WORKSPACE"); w != "" {
		return w
	}
	return filepath.Join(os.TempDir(), "aesp-workspace")
}

func cmdConfig() {
	ws := workspace()
	sys := agentos.New(agentos.Config{Workspace: ws, AutoApprove: false})
	fmt.Println("AESP Agent OS — runtime configuration")
	fmt.Println("=====================================")
	fmt.Printf("profile:           local-first (P2 defaults; serve = P1 HTTP)\n")
	fmt.Printf("workspace:         %s\n", ws)
	fmt.Printf("workspace env:     AESP_WORKSPACE=%s\n", os.Getenv("AESP_WORKSPACE"))
	fmt.Println()
	fmt.Println("memory (INV-04 / AESP-0004)")
	fmt.Println("  backend:        in-process pkg/memory.Memory")
	fmt.Println("  scopes:         working | session | semantic")
	fmt.Println("  trust labels:   system, verified, agent, retrieved, untrusted, poison-suspect")
	fmt.Println("  write rule:     trust label REQUIRED on every write")
	fmt.Println("  tools:          memory.write, memory.read")
	fmt.Println()
	fmt.Println("knowledge (AESP-0006)")
	fmt.Println("  backend:        in-process pkg/knowledge.Graph")
	fmt.Println("  tools:          kg.upsert")
	fmt.Println()
	fmt.Println("session / mission (AESP-0001, 0005, INV-10)")
	fmt.Println("  admission:      requiredCapabilities non-empty (INV-03)")
	fmt.Println("  host API:       pkg/host.Interface via pkg/kernel")
	fmt.Println("  journal:        pkg/eventbus + pkg/replay")
	fmt.Println("  execution tree: agents, artifacts, costs, timeline, failures")
	fmt.Println("  hitl timeout:   expire WITHOUT auto-approve (AESP-0014)")
	fmt.Println()
	fmt.Println("providers (INV-01) — abstract plugins only")
	for _, p := range sys.Providers.List() {
		d, err := p.Describe(context.Background())
		if err != nil {
			continue
		}
		fmt.Printf("  %-24s local=%v priority=%d caps=%v\n", d.ID, d.Local, d.Priority, d.Capabilities)
	}
	fmt.Println()
	fmt.Println("runtimes (INV-01 / INV-09)")
	for _, r := range sys.Runtimes.List() {
		d, err := r.Describe(context.Background())
		if err != nil {
			continue
		}
		fmt.Printf("  %-24s sandbox=%s capsIn=%v\n", d.ID, d.Sandbox, d.CapabilitiesIn)
	}
	fmt.Println()
	fmt.Println("routing (INV-03): capability-based + health failover (never model-name)")
	fmt.Println("credentials (INV-07): broker handles; raw secrets not logged")
	fmt.Println("mcp/a2a:          pkg/mcp + pkg/a2a (fixtures under conformance/fixtures/)")
	fmt.Println()
	fmt.Println("see also: config/default.yaml  make show-memory  make show-session")
}

func cmdDemo() error {
	sys := agentos.New(agentos.Config{AutoApprove: true, Workspace: workspace()})
	res, err := sys.RunMission(context.Background(), types.Mission{
		ID:              "wu_demo",
		Tenant:          "local",
		Goal:            "prove host-neutral kernel loop",
		RequiredCaps:    []types.Capability{"coding", "tools", "local"},
		SuccessCriteria: []string{"example-complete"},
		Budget:          types.Budget{MaxSteps: 10, MaxTokens: 2000, MaxWallSec: 60},
		CreatedAt:       time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	printResult(res)
	return nil
}

func cmdRun(path string) error {
	m, scenario, err := missionload.Load(path)
	if err != nil {
		return err
	}
	if scenario != "" && m.Labels != nil {
		m.Labels["scenario"] = scenario
	}
	// Infer scenario from path
	if m.Labels == nil {
		m.Labels = map[string]string{}
	}
	base := filepath.Base(filepath.Dir(path))
	if strings.Contains(base, "failover") {
		m.Labels["scenario"] = "failover"
	}
	if strings.Contains(base, "rollback") {
		m.Labels["scenario"] = "rollback"
	}
	sys := agentos.New(agentos.Config{AutoApprove: true, Workspace: workspace()})
	res, err := sys.RunMission(context.Background(), m)
	if err != nil {
		return err
	}
	printResult(res)
	return nil
}

func cmdRunExample(name string) error {
	path, err := findExample(name)
	if err != nil {
		return err
	}
	return cmdRun(path)
}

func cmdRunAllExamples() error {
	names := []string{
		"01-single-agent", "02-multi-agent", "03-code-generation", "04-review-approval",
		"05-memory-update", "06-kg-update", "07-remediation", "08-hitl",
		"09-provider-failover", "10-rollback-retry",
	}
	for _, n := range names {
		fmt.Printf("=== %s ===\n", n)
		if err := cmdRunExample(n); err != nil {
			return fmt.Errorf("%s: %w", n, err)
		}
	}
	fmt.Println("all examples succeeded")
	return nil
}

func findExample(name string) (string, error) {
	candidates := []string{
		filepath.Join("examples", name, "mission.yaml"),
		filepath.Join("..", "AESP-Examples", "examples", name, "mission.yaml"),
		filepath.Join(os.Getenv("HOME"), "git", "AESP-Examples", "examples", name, "mission.yaml"),
	}
	// also allow numeric-only
	if !strings.Contains(name, "-") {
		// no-op
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	return "", fmt.Errorf("example mission not found for %q (tried local examples/ and AESP-Examples)", name)
}

func cmdServe(addr string) error {
	sys := agentos.New(agentos.Config{AutoApprove: false, Workspace: workspace()})
	srv := httpapi.New(sys)
	fmt.Printf("AESP Agent OS Host Interface listening on %s\n", addr)
	fmt.Println("POST /v1/missions  GET /v1/missions/{id}/tree  GET /health  GET /v1/conformance")
	return http.ListenAndServe(addr, srv.Handler())
}

func printResult(res *agentos.MissionResult) {
	fmt.Printf("mission %s status=%s provider=%s runtime=%s\n", res.WorkUnitID, res.Status, res.ProviderID, res.RuntimeID)
	fmt.Printf("artifacts=%d events=%d costUSD=%.4f\n", len(res.Artifacts), len(res.Events), res.CostUSD)
	if len(res.Output) > 200 {
		fmt.Printf("output=%s...\n", res.Output[:200])
	} else {
		fmt.Printf("output=%s\n", res.Output)
	}
	if res.Tree != nil {
		fmt.Printf("executionTree agents=%v failures=%v\n", res.Tree.Agents, res.Tree.Failures)
	}
	b, _ := json.MarshalIndent(map[string]any{
		"status":     res.Status,
		"providerId": res.ProviderID,
		"runtimeId":  res.RuntimeID,
		"artifacts":  res.Artifacts,
		"eventTypes": eventTypes(res),
	}, "", "  ")
	fmt.Println(string(b))
}

func eventTypes(res *agentos.MissionResult) []string {
	out := make([]string, 0, len(res.Events))
	for _, e := range res.Events {
		out = append(out, e.Type)
	}
	return out
}
