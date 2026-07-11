# Configuration

| File | Purpose |
|------|---------|
| [default.yaml](./default.yaml) | Documents how memory, session, providers, tools, and HITL are wired |

## How configuration is applied

The reference Agent OS assembles defaults in code (`pkg/agentos.New`).  
`default.yaml` is the **human-readable contract** of those defaults so you can
see memory scopes, trust labels, session events, and routing without reading Go.

| Override | Effect |
|----------|--------|
| `AESP_WORKSPACE=/path` | Workspace root for file tools + local data |
| `make demo WORKSPACE=/path` | Same via Makefile |
| Mission `budget` / `requiredCapabilities` | Per-mission (mission.yaml) |

## Inspect

```bash
make show-config      # runtime dump + default.yaml
make show-memory      # memory subsystem cheat sheet
make show-session     # session / event lifecycle
make show-providers   # provider + runtime plugins
```

## Test by concern

```bash
make test-memory
make test-session
make test-routing
make test-hitl
make demo-memory
make demo-session
```
