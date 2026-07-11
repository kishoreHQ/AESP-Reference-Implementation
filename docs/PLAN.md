# AESP AGENT OS — FINAL PROGRAM PLAN (v1.0)

> **Single source of truth for what remains.**  
> Consolidates everything shipped through UI-GATE-7 and sequences open work (ADT-01…12) into two executable phases plus a post-program roadmap.  
> Sections 4–5 are execution prompts — hand them to the lead agent verbatim.  
> All prior rules apply: master prompt loop, GAP protocol, self-evaluation, swarm contract, separate reviewer sign-off, gates in `gates/`, process log in `evaluations/`.

---

## 1. PROGRAM STATE (as of commit 311d7ba; monorepo continues past this tip)

**Shipped and gated:**

- **Kernel:** full agent loop, capability routing, HITL, memory + trust labels, artifact store, event journal, conformance (`aespd conformance`), 10 examples, P1/P2/P3 profiles — **GATES 1–5**
- **Mission Control UI:** Missions, Mission Detail, Approvals, Fleet, Memory/KG, Artifacts, Evaluations, Replay, Settings — **UI-GATES 1–5**
- **Command Deck:** K1–K7 (connections/probe, adapters, sessions, board, routines, goals/journal, analytics), Connect wizard, Agent Rail, live sessions, Control Room, Kanban, cron routines, Brain rail — **UI-GATES 6–7**
- **Also shipped after 311d7ba (do not reopen):** monorepo merge, WebSocket `/api/v1/events` + `seq` (GAP-UI-001), SPA serve from `aespd` (GAP-UI-002)

**Known limits (define Phase 8):**

| Limit | Implication |
|-------|-------------|
| Named-CLI handshake is version + echo | Not a full agent loop yet |
| Model-switch events partially demo-path | Must come from real routing decisions |
| Board drag deferred | Action buttons shipped (a11y-first) |
| PTY adapter unsandboxed by design | Badge + policy tiers still needed |

**Document set (authoritative, in order):**

1. Master Execution Prompt  
2. `docs/ui/UI-SPEC.md`  
3. UI Production Build Prompt  
4. `docs/ui/FEATURE-PACK-CONNECTIONS.md`  
5. `docs/LANDSCAPE.md` (ADT-01…12)  
6. **This plan** (`docs/PLAN.md`)

---

## 2. PRINCIPLES (unchanged, re-stated once)

- Build **on top** of shipped work — no reopening passed gates; destructive refactors require an ADR.  
- **INV-01…11** hold: vendor names only in registry data; everything is a plugin; UI binds to `/api/v1` only; missing endpoints are GAPs closed kernel-side; every new capability lands in the conformance table.  
- Prefer additive requirement IDs; never rewrite existing requirement semantics.

---

## 3. HUMAN CHECKPOINTS (operator, not agents)

- **HC-1 (before Phase 8 starts):** 10-minute smoke of the shipped deck: `make dev` → Connect → generic-pty handshake → rail chat → board claim → routine fire. Agent-declared PASS + one human smoke test is the standard.  
- **HC-2 (GATE-8):** personally run the GATE-8 demo (§4) on MacBook (P2) — a named CLI full loop + one Telegram approval.  
- **HC-3 (GATE-9):** review the security posture summary (modes, scopes, signatures, tiers) before calling the program done.

---

## 4. UI-PHASE 8 — REAL AGENTS & REACH (execution prompt)

**Objective:** cross the line from “can connect agents” to “actually drives them daily.” Three parallel tracks; ⛓ marks hard ordering inside a track.

### Track A — Full agent loops (lead track)

- ⛓ **A1 / ADT-03 spike (timeboxed: one gate-task, max 2 days-equivalent):** run `rivet-dev/sandbox-agent` locally; drive Claude Code + OpenCode via its HTTP API; map its universal session schema onto our event journal. Exit with an ADR: **PASS** → K2 becomes one `sandbox-agent` runtime-host plugin + thin per-CLI manifests (bespoke adapters only for uncovered CLIs); **FAIL** → proceed bespoke per Feature Pack, ADR records why.  
- ⛓ **A2:** implement the chosen path until Claude Code, Codex CLI, and OpenCode each run a **REAL full agent loop** as sessions: streamed steps, structured tool calls, stop/kill working, costs/tokens captured where the CLI reports them (**ADT-04:** adapters translate INTO the AESP event model, never the reverse).  
- **A3:** real model-switch events — when session provider binding changes mid-session via routing, the journaled event comes from the **actual routing decision**, not the demo path.  
- **A4 / ADT-10:** heartbeat monitor per runtime/session; missed beats flip rail status to `error` with last-known state; policy may trigger remediation. (Do early in the track — a rail showing stale “working” is worse than no rail.)

### Track B — Reach (parallel from day 1)

- **B1 / ADT-11:** Telegram channel plugin: approval cards with inline Approve/Reject wired to `/approvals/:id/decision`; mission completion/failure notifications. Registry-driven channel-plugin contract (`channel.yaml`), zero kernel vendor names. Config in Settings with a “send test message” button.

### Track C — Smart routing (parallel, kernel-only)

- **C1 / ADT-01:** cost tiers (`free-local` / `free-hosted` / `budget` / `standard` / `premium`) on providers/models; tier chips in Fleet and session/mission meta; `free-first-coding` example policy (free/local first, auto-escalate on failed verification); “cost avoided” stat in Mission meta.  
- **C2 / ADT-02:** planner complexity score as routing input; score + chosen tier in the routing journal so replay explains model choice.  
- **C3 / ADT-07:** agent modes Full / Assist / Observe — registry field, live-switchable from Control Room and rail, mode changes journaled; Assist converts every external action into an approval; Observe journals without executing.

### GATE-8 (demo-shaped, verified by reviewer AND HC-2)

1. A real named CLI (Claude Code or OpenCode) completes a full task loop, fully visible in Mission Control: streamed steps, ≥1 structured tool call, working Stop — on desktop and 375px.  
2. One approval decided end-to-end from Telegram, reflected in the tree ≤5s.  
3. A `free-first-coding` mission routes to a local/free model, fails verification once, escalates to a premium tier — the whole path visible in Replay with complexity score and tier reasons.  
4. An agent in Assist mode generates an approval for an external write; the same action in Observe mode journals without executing.  
5. Heartbeat kill-test: kill an adapter process; rail flips to `error` ≤10s.  
6. Conformance table updated (ADT-01/02/03/04/07/10/11); process log complete; ADR for the spike committed; reviewer sign-off.

---

## 5. UI-PHASE 9 — GOVERNANCE & HARDENING TAIL (execution prompt)

- **D1 / ADT-05:** Connect wizard “import keys from local agent configs” — one-confirmation import into the credential broker; never displayed; provenance recorded.  
- **D2 / ADT-06:** sandbox tiers formalized in the trust model: `micro-vm` / `container` / `process-pty` (badged unsandboxed). `runtime.yaml` declares tier; policy can require minimum tier per tool scope (external writes ≥ container). Micro-vm implementation may wrap the rivet OSS substrate if the A1 ADR recommended it; otherwise tier exists in policy with container as the strongest available.  
- **D3 / ADT-08:** `memory_read` / `memory_write` glob scopes in agent manifests, kernel-enforced; violations are policy events surfaced in the Control Room.  
- **D4 / ADT-09:** Ed25519-signed plugin manifests; `require_signed` policy with allowed-signer list; manifest hash auto-logged to the journal on load; hash-chain check added to `aespd conformance`.  
- **D5 / ADT-12 (last, optional, off by default):** OpenAI-compatible `/v1/chat/completions` ingress mapping onto sessions; P2-local only unless auth configured.

### GATE-9 (final program gate)

- Signed-manifest rejection test passes (unsigned plugin refused under `require_signed`)  
- A memory-scope violation is blocked and journaled  
- Sandbox tier policy blocks an external write from a PTY-tier runtime  
- Credential import round-trip works without ever displaying a key  
- Full conformance run green  
- Security posture summary written for HC-3  
- Reviewer sign-off  

---

## 6. POST-PROGRAM ROADMAP (explicitly out of scope until GATE-9)

In priority order, each requiring its own mini-plan before starting:

1. Repo topology migration per Phase-5 ADR (Hermes-Kernel / -Providers / -Runtimes / …) — only when the monorepo hurts  
2. Additional channel adapters beyond Telegram  
3. Skills/mission marketplace (openfang “Hands” equivalent)  
4. Deeper P3 embeds (OpenClaw gateway adapter, `mountMissionControl` consumers)  
5. Evaluation benchmarking against public agent benchmarks  

Items in `docs/LANDSCAPE.md` §3 remain **not-adopted** until re-prioritized.

---

## 7. PARALLELIZATION MAP

| Week | Work |
|------|------|
| **1** | HC-1 → A1 spike ⟂ B1 Telegram ⟂ C1 tiers |
| **2** | A2 loops ⟂ C2/C3 ⟂ A4 heartbeat → A3 → **GATE-8 + HC-2** |
| **3** | D1–D4 (D2 informed by A1 ADR) → D5 → **GATE-9 + HC-3** |

Assign one author-agent per track, one shared reviewer; conformance-checker runs at each gate.

---

## 8. DEFINITION OF PROGRAM DONE

**GATE-9 passed + HC-3 signed.**

Concretely: you open Mission Control on your phone, watch Claude Code and OpenCode work real tasks routed free-first with premium escalation, approve the risky steps from Telegram, replay any run to see why every model was chosen, and every plugin that loaded was signed and every memory write was in scope.

Everything after that is **roadmap**, not program.

---

## 9. QUICK REFERENCE — ARTIFACTS

| Artifact | Path |
|----------|------|
| This plan | `docs/PLAN.md` |
| Landscape / ADTs | `docs/LANDSCAPE.md` |
| UI audit | `docs/ui/AUDIT.md` |
| Host contract | `docs/ui/CONTRACT.md` |
| Connections pack | `docs/ui/FEATURE-PACK-CONNECTIONS.md` |
| Conformance | `docs/ui/CONFORMANCE.md` |
| Process log | `evaluations/PROCESS-LOG.md` |
| Gaps | `gaps/GAP-*.md` |
| Gates | `gates/GATE-*.md`, `gates/UI-GATE-*.md` |
