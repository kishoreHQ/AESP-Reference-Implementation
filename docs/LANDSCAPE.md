# Landscape & adoption decisions (ADT-01…12)

Industry patterns observed while building the AESP Agent OS and Mission Control.  
**ADT** = Adoption Decision Ticket — tracked here; implemented under `docs/PLAN.md` Phase 8–9.

## §1 — Adopt (in program Phases 8–9)

| ID | Title | Phase | Status |
|----|-------|-------|--------|
| ADT-01 | Cost tiers on providers/models + free-first policy | 8 / C1 | **open** |
| ADT-02 | Planner complexity score → routing journal | 8 / C2 | **open** |
| ADT-03 | sandbox-agent / rivet-style runtime-host spike | 8 / A1 | **open** |
| ADT-04 | Adapters translate INTO AESP events (never reverse) | 8 / A2 | **open** (partial: stream kinds exist) |
| ADT-05 | Import keys from local agent configs | 9 / D1 | **open** |
| ADT-06 | Sandbox tiers micro-vm / container / process-pty | 9 / D2 | **open** (PTY badge only) |
| ADT-07 | Agent modes Full / Assist / Observe | 8 / C3 | **open** |
| ADT-08 | Memory read/write glob scopes | 9 / D3 | **open** |
| ADT-09 | Ed25519-signed plugin manifests | 9 / D4 | **open** |
| ADT-10 | Heartbeat per runtime/session | 8 / A4 | **open** |
| ADT-11 | Telegram channel plugin for HITL | 8 / B1 | **open** |
| ADT-12 | Optional OpenAI-compatible ingress | 9 / D5 | **open** (off by default) |

## §2 — Already adopted (shipped)

| Pattern | Where |
|---------|--------|
| Provider ≠ Runtime | INV-01, separate registries |
| Capability routing | INV-03, `pkg/router` |
| Unified memory + trust labels | INV-04, `pkg/memory` |
| Host Interface only for UI | INV-11, `/api/v1` |
| Runtime.yaml discovery | INV-09, adapters + plugins |
| Event journal + seq | INV-10, bus + WS |
| HITL no auto-approve | AESP-0014 |
| Registry-driven Fleet / Agent Rail | UI-FLT-01 |

## §3 — Not adopted (out of program until GATE-9)

| Pattern | Why not now |
|---------|-------------|
| Hardcoded per-vendor UI panels | Violates INV-02 / registry thesis |
| Second memory store for “brain” | Violates INV-04 |
| Forking AESP semantics in product | INV-08 |
| Full multi-tenant SaaS control plane | Roadmap / ops product |
| Skills marketplace | Post-program §6 item 3 |
| Replacing monorepo with many repos | Only when monorepo hurts |

## §4 — Gap files

Phase 8–9 work files `gaps/GAP-ADT-NN.md` when implementing, or updates this table’s Status column.
