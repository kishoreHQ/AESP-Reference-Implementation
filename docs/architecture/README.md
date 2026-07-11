# Architecture Document Set — Hermes Agent OS (Reference Implementation)

**Status:** Phase 1 complete  
**Profile:** Host-neutral Agent OS kernel (AESP-conformant)  
**Date:** 2026-07-11  
**Invariants:** INV-01 … INV-11 (see [system-overview.md](./system-overview.md))

This directory is the **canonical architecture set** for the AESP Reference Implementation.
Documents are implementation-oriented but **vendor-neutral** and **host-neutral**.
Core architecture text MUST NOT name specific model vendors or product hosts except in
explicitly labeled non-normative examples.

## Document index

| # | Document | Realizes |
|---|----------|----------|
| 01 | [system-overview.md](./system-overview.md) | All INV, suite map |
| 02 | [runtime-loop.md](./runtime-loop.md) | Agent harness loop |
| 03 | [control-plane.md](./control-plane.md) | Policy, HITL, budgets |
| 04 | [compute-plane.md](./compute-plane.md) | Inference + sandboxed tools |
| 05 | [agent-lifecycle.md](./agent-lifecycle.md) | Agent principals |
| 06 | [session-lifecycle.md](./session-lifecycle.md) | Sessions & WorkUnits |
| 07 | [event-model.md](./event-model.md) | Event bus, registry |
| 08 | [trust-model.md](./trust-model.md) | Trust labels, boundaries |
| 09 | [policy-model.md](./policy-model.md) | Policy engine |
| 10 | [artifact-model.md](./artifact-model.md) | Digests, provenance |
| 11 | [tool-model.md](./tool-model.md) | Unified tools / MCP |
| 12 | [provider-model.md](./provider-model.md) | INV-01 providers |
| 13 | [runtime-model.md](./runtime-model.md) | INV-01/09 runtimes |
| 14 | [context-envelope.md](./context-envelope.md) | INV-05 |
| 15 | [host-interface.md](./host-interface.md) | INV-11 |
| 16 | [memory-model.md](./memory-model.md) | INV-04 |
| 17 | [credential-model.md](./credential-model.md) | INV-07 |
| 18 | [replay-and-audit.md](./replay-and-audit.md) | INV-10 |

## Gate

Gate artifact: [`../../gates/GATE-1.md`](../../gates/GATE-1.md)
