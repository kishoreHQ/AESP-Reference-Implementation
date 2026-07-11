# ADR-0004: Phase 8 A1 — sandbox-agent / rivet-style runtime host (placeholder)

| Metadata | Value |
|----------|-------|
| **Status** | Proposed (spike not yet run) |
| **Date** | 2026-07-11 |
| **Related** | ADT-03, PLAN.md §4 Track A |

## Context

Named CLI adapters currently handshake with version + echo. Full agent loops need either:

1. **Host path:** one runtime-host plugin (e.g. rivet-dev/sandbox-agent) + thin per-CLI manifests  
2. **Bespoke path:** per-CLI structured adapters parsing stdout  

## Decision (pending spike)

**Not decided.** Phase 8 A1 must:

1. Attempt local `sandbox-agent` (or equivalent) drive of Claude Code + OpenCode  
2. Map session schema → AESP journal events  
3. Exit with **PASS** (adopt host) or **FAIL** (bespoke) recorded by updating this ADR to Accepted  

## Consequences

- Until A1 completes, GATE-8 cannot pass item 1  
- D2 sandbox tiers may wrap host if PASS  

## Exit checklist

- [ ] Spike environment documented  
- [ ] Event mapping table committed  
- [ ] Status → Accepted with PASS or FAIL  
