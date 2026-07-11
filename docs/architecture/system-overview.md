# System Overview — Host-Neutral Agent Operating System

**Status:** Normative for this reference implementation  
**AESP:** AESP-0000–0015 (suite spine)  
**Invariants:** INV-01 … INV-11

## 1. Purpose

This reference implementation is a **vendor-neutral, host-neutral AI middleware platform**
(an Agent Operating System). Hosts sit above the kernel; models and runtimes sit below it.

```
Host Layer:  Platform UI · Local CLI/TUI · External orchestrators · Headless API
                              ↓  Host Interface (INV-11)
                    Agent Runtime Kernel
           {Planning | Execution | Memory | Knowledge | Policy}
                              ↓
              Provider Router  ·  Runtime Registry  ·  Tool Layer
                              ↓
        Capability-compatible providers & plugin runtimes (INV-01, INV-02, INV-03)
```

**Models sit BELOW the runtime, never above it.**

## 2. Architectural invariants (non-negotiable)

| ID | Invariant | Primary docs |
|----|-----------|--------------|
| INV-01 | Provider ≠ Runtime | provider-model, runtime-model |
| INV-02 | Everything is a plugin | runtime-model, tool-model, provider-model |
| INV-03 | Capability-based routing (never model-name routing) | provider-model, capability |
| INV-04 | Unified memory | memory-model |
| INV-05 | Unified context envelope | context-envelope |
| INV-06 | Unified tool layer | tool-model |
| INV-07 | Unified credentials | credential-model |
| INV-08 | AESP is the contract | this doc, conformance |
| INV-09 | Dynamic runtime registry | runtime-model |
| INV-10 | Auditable by construction | replay-and-audit, event-model |
| INV-11 | Host-neutral core + deployment profiles | host-interface, deployment/* |

## 3. Layered model (maps to AESP)

| Layer | Responsibility | AESP |
|-------|----------------|------|
| Host Interface | Mission submit, events, approvals, artifacts | 0014, 0015, Host |
| Control Plane | Policy, budgets, HITL, session authority | 0013, 0014, 0001, 0005 |
| Agent Runtime Kernel | Loop: accept → context → plan → act → verify → persist | AGENT-RUNTIME, 0001–0005 |
| Compute Plane | Provider inference + sandboxed tool execution | 0015, MCP |
| Memory & Knowledge | Unified memory + KG | 0004, 0006 |
| Evidence Plane | Traces, eval, deploy proofs, replay | 0011, 0010, 0009 |
| Security cross-cut | Trust labels, authz, isolation | 0013 |

## 4. Control loop (production)

1. **Intent** — WorkUnit created (AESP-0001), authorized by role (0002).
2. **Plan/Orchestrate** — Workflow graph (0005) dispatches over (0003).
3. **Context** — Memory (0004) + KG (0006) assemble Context Envelope (INV-05).
4. **Route** — Capability Engine → Routing Policy → Provider + Runtime (INV-03).
5. **Act** — Tools (0015/MCP) under policy; peers via A2A where applicable.
6. **Verify** — Tests (0010), validators, human gates (0014).
7. **Persist** — Memory/KG updates with trust labels; artifacts with digests.
8. **Observe** — Telemetry (0011) correlates WorkUnit → sessions → traces.
9. **Remediate** — Playbooks (0012) with HITL as needed.
10. **Govern** — Security (0013) + Constitution (0000).

## 5. Deployment profiles (INV-11)

| Profile | Mode | Minimum backends |
|---------|------|------------------|
| **P1 Platform** | Server / multi-tenant | durable store, multi-provider, HITL UI host |
| **P2 Local-first** | Single machine, offline-capable | file/SQLite memory, local providers, zero cloud credentials |
| **P3 Embedded** | Library/sidecar in external host | Host Interface SDK only |

A workflow authored under one profile MUST run under the others, degrading only by
**declared capability** (e.g., no vision model available).

## 6. What the kernel does NOT own

- Product UI branding or chat UX
- Specific model vendor SDKs (those are **provider plugins**)
- Specific coding-agent CLIs (those are **runtime plugins**)
- Per-runtime memory silos (forbidden by INV-04)

## 7. AESP requirement families realized

MEM · WF · KG · CG · DOC · DEP · TEST · OBS · REM · SEC · HITL · INT  
(See suite `specification/ARCHITECTURE.md` and `CONFORMANCE.md`.)

## 8. Target multi-repo topology (ADR)

See [ADR-0003 target topology](../adr/0003-target-repo-topology.md):

`Kernel · Providers · Runtimes · Tools · Memory · Workflow · Evaluation · Host UI · SDK · Starter-Agents`

## 9. Normative language

RFC 2119: MUST / SHOULD / MAY. Conformance claims use suite profiles in AESP CONFORMANCE.md.
