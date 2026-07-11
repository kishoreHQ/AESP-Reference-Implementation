# AESP Reference Implementation — Host-Neutral Agent OS

Production-oriented **reference implementation** of the [Autonomous Engineering Specification (AESP)](https://github.com/kishoreHQ/AESP).

This is **not** a chat product. It is vendor-neutral, host-neutral **AI middleware**:

```
Host (Platform UI · Local CLI · External orchestrator · API)
        ↓ Host Interface
Agent Runtime Kernel  →  Planning · Execution · Memory · Knowledge · Policy
        ↓
Provider plugins · Runtime plugins · Tool layer
```

**Models sit below the runtime. Hosts sit above it.**

## Invariants

INV-01 Provider ≠ Runtime · INV-02 Plugins · INV-03 Capability routing ·  
INV-04 Unified memory · INV-05 Context envelope · INV-06 Unified tools ·  
INV-07 Unified credentials · INV-08 AESP is the contract · INV-09 Runtime registry ·  
INV-10 Auditable by construction · INV-11 Host-neutral + P1/P2/P3 profiles

## Quick start

```bash
go test ./...
go run ./cmd/aespd
go run ./cmd/aespd conformance
```

## Layout

| Path | Purpose |
|------|---------|
| `docs/architecture/` | Phase 1 architecture set |
| `pkg/` | Kernel modules (Go) |
| `plugins/` | Provider & runtime plugins |
| `cmd/aespd` | Reference daemon / CLI entry |
| `gates/` | GATE-1 … GATE-5 |
| `gaps/` | Protocol/implementation gaps |
| `evaluations/PROCESS-LOG.md` | Process log |
| `docs/deployment/` | P1 / P2 / P3 guides |
| `docs/security/`, `docs/hardening/` | Threat model, failures, replay |

## Related repos

- Spec: https://github.com/kishoreHQ/AESP  
- Examples: https://github.com/kishoreHQ/AESP-Examples  

## License

MIT
