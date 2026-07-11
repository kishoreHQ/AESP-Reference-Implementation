# Event Model

**AESP:** 0011, EVENT-REGISTRY, 0003

## Principles

1. Every meaningful state change is an event.
2. Events are append-only for a mission's audit journal.
3. Events carry correlation keys (WorkUnit, session, trace, tenant).
4. Hosts subscribe via Host Interface; kernel does not assume a UI.

## Envelope (minimum)

```json
{
  "type": "aesp.runtime.step.completed",
  "id": "evt_...",
  "time": "RFC3339",
  "source": "kernel|plugin:<id>",
  "tenant": "...",
  "workUnitId": "...",
  "sessionId": "...",
  "traceId": "...",
  "data": {}
}
```

## Categories

| Prefix | Owner |
|--------|-------|
| `aesp.control.*` | Control plane |
| `aesp.runtime.*` | Agent loop |
| `aesp.tool.*` | Tool invocations |
| `aesp.provider.*` | Provider calls |
| `aesp.memory.*` | Memory writes |
| `aesp.artifact.*` | Artifact lifecycle |
| `aesp.hitl.*` | Approvals |
| `aesp.obs.*` | Observability signals |
| `aesp.rem.*` | Remediation |

Full names: suite `specification/EVENT-REGISTRY.md`.

## Bus interface

`pkg/eventbus` — Publish, Subscribe, ReplayFrom(offset|time).
