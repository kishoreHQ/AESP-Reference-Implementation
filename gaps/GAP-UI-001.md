# GAP-UI-001 — WebSocket event stream with monotonic seq

**Status:** CLOSED (implemented this change)  
**Impact:** Realtime UI, approval ≤1s tree refresh, Replay, reconnect recovery.

## Evidence

Audit found no `/api/v1/events` WebSocket; EventBridge expected `seq` and `since=`.

## Implementation

- `pkg/eventbus`: expose sequence numbers on publish  
- `pkg/httpapi`: `GET /api/v1/events` upgrades to WebSocket; fan-out with mapping; `since` filter  
- UI EventBridge uses `/api/v1/events?since=`
