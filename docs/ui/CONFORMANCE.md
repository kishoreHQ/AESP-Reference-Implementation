# UI-SPEC Conformance Table (GATE-5)

| ID | Title | Result | Notes |
|----|-------|--------|-------|
| UI-ARCH-01 | Host Interface only | **pass** | `/api/v1` + WS |
| UI-ARCH-02 | P1/P2/P3 one codebase | **pass** | `ui/` + embed |
| UI-ARCH-03 | Responsive | **pass** | shell breakpoints |
| UI-ARCH-04 | No vendor hardcodes | **pass** | registry-rendered |
| UI-SIG-01 | Mission Spine | **pass** | virtualized |
| UI-COL-* | Tokens dark/light | **pass** | tokens.css |
| UI-HITL-01 | Inline previews | **pass** | ApprovalsPage |
| UI-HITL-02 | Decision ≤1s | **pass** | invalidate queries + WS |
| UI-FLT-01 | Registry-driven fleet | **pass** | Fleet tabs |
| UI-MEM-01 | Trust + provenance | **pass** | MemoryPage |
| UI-ART-01 | Artifact deep link | **pass** | link to mission |
| UI-RPL-01 | Replay banner | **pass** | ReplayPage |
| UI-SEC-02 | Creds write-only | **pass** | Settings |
| UI-STA-01 | Four UI states | **pass** | loading/empty/error/data |
| UI-RSP-* | Breakpoints / 44px | **pass** | AppShell mobile |
| UI-A11Y-* | Keyboard / focus | **partial** | focus rings; axe not in CI |
| UI-PERF-01 | Bundle ≤300KB gz | **pass** | main ~108KB gz (+mock worker split) |
| UI-PERF-02 | Virtualize logs/tree | **pass** | TanStack Virtual |
| UI-RT-01 | seq reconnect | **pass** | GAP-UI-001 closed |
| UI-API-01 | No invented endpoints | **pass** | gaps filed/closed |
| UI-TEC-06 | mountMissionControl | **pass** | embed.tsx |
| UI-TST-01 | Contract mocks | **partial** | MSW parity; CI diff deferred |

Unresolved / deferred: formal axe CI, Playwright visual golden suite, full 100k log soak.

## Feature pack — Connections & Command Deck

| ID | Result |
|----|--------|
| K1 Connection probe | **pass** |
| K2 Runtime adapters | **pass** (generic-pty + CLI adapters) |
| K3 Sessions | **pass** |
| K4 Task board | **pass** |
| K5 Routines | **pass** |
| K6 Goals/journal | **pass** |
| K7 Analytics | **pass** |
| UI Connections wizard | **pass** |
| UI Agent rail | **pass** |
| UI Live session | **pass** |
| UI Control room | **pass** |
| UI Board | **pass** |
| UI Routines | **pass** |
| UI Brain rail | **pass** |
