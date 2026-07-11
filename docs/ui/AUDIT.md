# Mission Control UI — Production Audit (TASK 0)

**Date:** 2026-07-11  
**Repo:** AESP-Reference-Implementation (monorepo)  
**Purpose:** Ground truth before production UI build. Every later task cites this file.

## (a) UI feature × state

| Feature area | Route(s) | Components | State | Notes |
|--------------|----------|------------|-------|-------|
| Dashboard / Overview | `/` | `DashboardPage` | **partial** | KPIs, dispatch chips, activity; not all live-linked |
| Missions home | `/missions` | `MissionsPage` | **partial** | Board + kanban; ⌘K not wired; four states OK |
| Mission Detail | `/missions/:id` | `MissionDetailPage` | **partial** | Tree/logs/meta present; tool blocks limited; node actions incomplete |
| Approval Inbox | `/approvals` | `ApprovalsPage` | **partial** | Previews + decisions; escalate is toast-only; mobile OK |
| Fleet | `/fleet` | `FleetPage` | **partial** | Registry tabs; enable/disable not persisted; latency sparklines mock |
| Memory & KG | `/memory` | `MemoryPage` | **partial** | Search + trust chips; pin/quarantine wired; KG canvas minimal |
| Artifacts | `/artifacts` | `ArtifactsPage` | **partial** | Preview; version diff shallow |
| Evaluations | `/evaluations` | `EvaluationsPage` | **partial** | Table; needs real eval runs from kernel |
| Replay | `/replay/:runId` | `ReplayPage` | **partial** | Step player; tree reconstruction incomplete |
| Settings | `/settings` | `SettingsPage` | **partial** | Creds write-only UI; policies/budgets stubbed live |
| Mission Spine | shell | `MissionSpine` | **partial** | Virtualized; 500-node OK on mock tree |
| Layout / a11y | shell | `AppShell` | **partial** | Breakpoints OK; axe not gated yet |
| Themes | tokens | `tokens.css` | **complete** | Dark+light Cherenkov; hex confined |
| Embed P3 | `embed.tsx` | `mountHermesUI` | **partial** | Exists; rename to `mountMissionControl` |
| Realtime | EventBridge | `EventBridge.ts` | **partial** | Mock ticker OK; live WS path wrong/missing kernel |
| Mocks | MSW | `mocks/` | **partial** | Coverage good; progression ticker limited |

## (b) Endpoint × state (`/api/v1/*` only for UI)

| Endpoint | State | Notes |
|----------|-------|-------|
| `GET /api/v1/health` | **served** | |
| `GET/POST /api/v1/missions` | **served** | POST runs real OS; list is in-memory UI store |
| `GET /api/v1/missions/:id` | **served** | |
| `POST /api/v1/missions/:id/cancel` | **served** | |
| `GET /api/v1/missions/:id/tree` | **stubbed** | Static 3-node tree |
| `GET /api/v1/missions/:id/logs` | **stubbed** | Static 3 lines |
| `GET /api/v1/approvals` | **stubbed** | Empty list unless seeded |
| `POST /api/v1/approvals/:id/decision` | **served** | HITL service resolve |
| `GET /api/v1/registry/{kind}` | **served** | providers/runtimes live; agents/tools partial |
| `GET /api/v1/memory/search` | **served** | From memory store |
| `GET /api/v1/memory/kg` | **stubbed** | Placeholder graph |
| `POST /api/v1/memory/:id/{action}` | **stubbed** | Returns ok only |
| `GET /api/v1/artifacts` | **stubbed** | Empty |
| `GET /api/v1/artifacts/:id/versions` | **stubbed** | Empty |
| `GET /api/v1/evaluations` | **stubbed** | Empty |
| `GET /api/v1/replay/:runId/events` | **served** | From bus replay |
| `GET /api/v1/budgets` | **stubbed** | One placeholder |
| `GET/PUT /api/v1/policies` | **stubbed** | |
| `POST /api/v1/credentials` | **served** | Broker put |
| `GET /api/v1/events` (WS) | **missing** | **GAP-UI-001** |
| Static UI `/*` from aespd | **missing** | **GAP-UI-002** P1 packaging |

Legacy `/v1/*` and `/health` exist — UI must not use them.

## (c) Event type × state

| Event type (UI contract) | State | Notes |
|--------------------------|-------|-------|
| Monotonic `seq` on bus events | **partial** | Bus has internal seq; not exposed on envelope |
| `mission.updated` | **missing** WS | Published as `aesp.*` names on bus |
| `node.updated` | **missing** | |
| `log.append` | **missing** | |
| `approval.created` / `approval.resolved` | **partial** | Bus has hitl events; no WS fanout |
| `artifact.created` | **partial** | Bus events exist |
| `memory.written` | **partial** | Bus events exist |
| `budget.threshold` | **missing** | |
| `eval.completed` | **missing** | |
| WebSocket reconnect `since=seq` | **missing** | **GAP-UI-001** |

Kernel bus types today use `aesp.*` (e.g. `aesp.control.mission.accepted`). CONTRACT.md maps these to UI event names.

## Disposition

| ID | Action |
|----|--------|
| GAP-UI-001 | Implement `/api/v1/events` WebSocket + seq on payloads |
| GAP-UI-002 | Serve `ui/dist` from aespd for P1 |
| Tree/logs/artifacts | Upgrade stubs to use real mission results + journal |
| Approvals seed | Create HITL tasks on approval-gated missions |
| UI polish | Production states, node actions, mock parity ticker |

## Feature pack delta (post UI-PHASE 6–7)

| Endpoint | State |
|----------|-------|
| POST/GET /api/v1/connections/probe | **served** |
| POST/GET /api/v1/connections | **served** |
| POST/GET /api/v1/sessions | **served** |
| POST /api/v1/sessions/:id/message\|stop | **served** |
| GET /api/v1/boards | **served** |
| GET/POST /api/v1/tasks | **served** |
| POST /api/v1/tasks/:id/claim | **served** |
| GET/POST /api/v1/routines | **served** |
| GET /api/v1/goals, journal | **served** |
| GET /api/v1/analytics/agents/:id | **served** |

| UI route | State |
|----------|-------|
| /connect | **complete** |
| /agents/:id | **complete** |
| /sessions/* | **complete** |
| /board | **complete** |
| /routines | **complete** |
| Brain rail (home) | **complete** |
