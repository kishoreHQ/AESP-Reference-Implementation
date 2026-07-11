# Session Lifecycle

**AESP:** 0001, 0005, 0007–0012 session IRIs, 0011

## Session kinds

| Kind | Purpose |
|------|---------|
| mission | Top-level host-submitted mission |
| workflow | Durable orchestration instance |
| runtime | Single runtime plugin execution span |
| tool | Optional per-tool long session |
| codegen / docs / deploy / test / remediate | Domain sessions per suite |

## Lifecycle

```
created → open → (checkpointed)* → closing → closed
open → aborted
```

Sessions MUST carry: session IRI, WorkUnit id, tenant, budget snapshot, started_at, ended_at.

## Correlation

Production systems MUST pivot WorkUnit → sessions → traces → HITL tasks → artifact digests.
