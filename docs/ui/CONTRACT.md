# Host Interface Contract (UI ↔ Kernel)

**UI binds only to `/api/v1/*` and WebSocket `/api/v1/events`.**  
Envelope: `{ "data": T, "error": { "code", "message", "remediation" } | null }`.

## Event name mapping

| Kernel bus type (`aesp.*`) | UI HostEvent type |
|----------------------------|-------------------|
| `aesp.control.mission.accepted` | `mission.updated` |
| `aesp.control.mission.cancelled` | `mission.updated` |
| `aesp.runtime.completed` / `failed` | `node.updated` |
| `aesp.tool.invoked` | `log.append` |
| `aesp.hitl.approval.requested` | `approval.created` |
| `aesp.hitl.approval.resolved` | `approval.resolved` |
| `aesp.artifact.created` | `artifact.created` |
| `aesp.memory.write` | `memory.written` |
| `aesp.control.route.selected` | `mission.updated` |
| `aesp.provider.completed` | `log.append` |

WS payload:

```json
{
  "seq": 42,
  "type": "mission.updated",
  "ts": "RFC3339",
  "missionId": "mis_…",
  "data": { "rawType": "aesp.control.mission.accepted", ... }
}
```

Reconnect: `ws://host/api/v1/events?since=<lastSeq>&mission=<optional>`.

## Delta from UI-SPEC §6

| Spec | Kernel reality | Disposition |
|------|----------------|------------|
| REST paths under `/api/v1` | Matched | OK |
| Cursor pagination | Not required for v1 list sizes | deferred |
| `seq` on all events | Added in production build | OK after GAP-UI-001 |
