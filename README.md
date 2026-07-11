# AESP Reference Implementation — Host-Neutral Agent OS

Production-oriented **reference implementation** of the [Autonomous Engineering Specification (AESP)](https://github.com/kishoreHQ/AESP).

This is **not** a chat product. It is vendor-neutral, host-neutral **AI middleware**:

```
Host (Platform UI · Local CLI · External orchestrator · HTTP API)
        ↓ Host Interface
Agent Runtime Kernel  →  Planning · Execution · Memory · Knowledge · Policy
        ↓
Provider plugins · Runtime plugins · Tool layer (MCP-aligned) · A2A peers
```

**Models sit below the runtime. Hosts sit above it.**

## Status

| Area | State |
|------|-------|
| Full agent loop | **Functional** (`pkg/agentos`) |
| All 10 examples | **Runnable** (`aespd run-all-examples`) |
| Conformance catalog | **28/28 implemented** |
| Profiles | P1 HTTP serve · P2 CLI local · P3 embed `agentos.System` |

## Quick start

```bash
go test ./...
go run ./cmd/aespd demo
go run ./cmd/aespd run-all-examples
go run ./cmd/aespd conformance
go run ./cmd/aespd run examples/01-single-agent/mission.yaml
go run ./cmd/aespd serve :8080
```

### HTTP Host Interface (P1/P3)

```bash
# Submit mission
curl -s localhost:8080/v1/missions -d '{
  "id":"wu_http","goal":"demo","requiredCapabilities":["coding","tools"],
  "successCriteria":["example-complete"],"budget":{"maxSteps":10}
}'
# Health
curl -s localhost:8080/health
```

## Agent loop (AESP-aligned)

1. Accept WorkUnit (0001) + capability requirements (INV-03)
2. Plan artifact (0015) → content-addressed store (0007)
3. Assemble Context Envelope (INV-05): memory, tools, policy, budget
4. Route provider + runtime by **capabilities** with failover (INV-01, INV-03)
5. Provider complete + runtime execute (compute plane)
6. Tools via unified router / MCP surface (INV-06)
7. HITL gates never auto-approve on timeout (0014)
8. Verify (0010) · memory write with trust labels (0004) · docgen (0008)
9. Optional deploy session (0009) · remediation playbooks (0012)
10. Execution tree + event journal (INV-10)

## Layout

| Path | Purpose |
|------|---------|
| `pkg/agentos` | Fully wired OS + mission loop |
| `pkg/kernel` | Host Interface core |
| `pkg/httpapi` | HTTP Host Interface |
| `pkg/{provider,runtime}registry` | Separate plugin registries |
| `pkg/router` | Capability routing + failover |
| `pkg/memory`, `knowledge`, `artifact` | Unified memory / KG / digests |
| `pkg/policy`, `approval`, `credentials` | Control plane |
| `pkg/deploy`, `remediation`, `docgen` | Ship / heal / docs |
| `pkg/mcp`, `a2a` | Interop + golden fixtures |
| `plugins/` | Mock local/remote providers, generic runtime |
| `examples/` | Bundled mission YAMLs (portable) |
| `conformance/fixtures/` | MCP + A2A golden fixtures |
| `docs/architecture/` | INV-01…INV-11 architecture set |
| `gates/`, `gaps/` | Program gates and gap register |

## Invariants

INV-01 Provider ≠ Runtime · INV-02 Plugins · INV-03 Capability routing ·  
INV-04 Unified memory · INV-05 Context envelope · INV-06 Unified tools ·  
INV-07 Unified credentials · INV-08 AESP is the contract · INV-09 Runtime registry ·  
INV-10 Auditable by construction · INV-11 Host-neutral + P1/P2/P3 profiles

## Related repos

- Spec: https://github.com/kishoreHQ/AESP  
- Examples: https://github.com/kishoreHQ/AESP-Examples  

## License

MIT
