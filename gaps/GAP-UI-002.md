# GAP-UI-002 — Serve Mission Control static UI from aespd (P1)

**Status:** CLOSED (implemented this change)  
**Impact:** Single-process Platform profile; identical `ui/dist` artifact.

## Implementation

- `aespd serve` optionally serves `ui/dist` when present  
- Env `AESP_UI_DIST` overrides path  
- SPA fallback to `index.html`
