# AESP Agent OS (monorepo)

Production-oriented **reference implementation** of the [Autonomous Engineering Specification (AESP)](https://github.com/kishoreHQ/AESP) — host-neutral AI middleware plus Mission Control UI and examples **in one repo**.

```
┌─────────────────────────────────────────────────────────┐
│  ui/          Mission Control (React)  →  :5173         │
│  cmd/aespd    Agent OS kernel (Go)     →  :8080         │
│  examples/    Portable mission YAMLs                    │
│  docs/        Architecture, gates, deployment           │
└─────────────────────────────────────────────────────────┘
         Host UI  →  Host Interface (/api/v1)  →  Kernel
```

**Models sit below the runtime. Hosts sit above it.**  
**Spec stays separate:** [kishoreHQ/AESP](https://github.com/kishoreHQ/AESP) (protocol suite — not forked here).

---

## Monorepo layout

| Path | What |
|------|------|
| `cmd/aespd`, `pkg/`, `plugins/` | Agent OS kernel |
| `ui/` | Mission Control UI (Vite + React) |
| `examples/` | 10 profile-portable missions |
| `config/` | Documented defaults (memory, session, routing) |
| `scripts/dev.sh` | Start **kernel + UI** together |
| `docs/` | Architecture, gates, security |

Related standalone repos (optional; product code lives here):

- Spec: https://github.com/kishoreHQ/AESP  
- Legacy UI clone: https://github.com/kishoreHQ/hermes-mission-control-ui (prefer `ui/` in this monorepo)  
- Legacy examples clone: https://github.com/kishoreHQ/AESP-Examples (prefer `examples/` here)

---

## Run everything (recommended)

```bash
cd ~/git/AESP-Reference-Implementation

# one-time
make build                 # Go kernel binary
make install-ui            # npm install in ui/

# daily: kernel :8080 + UI :5173
make dev
# or:  ./scripts/dev.sh
```

Then open:

| Service | URL |
|---------|-----|
| **Mission Control UI** | http://127.0.0.1:5173 |
| **Kernel health** | http://127.0.0.1:8080/api/v1/health |

Hard-refresh the browser after UI changes (`Cmd+Shift+R`).

### UI only (mocks, no kernel)

```bash
make dev-ui-mocks
# or: ./scripts/dev-ui-only.sh
```

### Kernel only (CLI)

```bash
make build
./bin/aespd demo
./bin/aespd run-all-examples
./bin/aespd serve :8080
```

---

## Quick Makefile map

```bash
make help              # all targets
make dev               # monorepo: kernel + UI
make test              # Go package tests
make smoke             # kernel tests + demos + examples
make demo-memory       # memory example
make show-config       # memory / session / provider wiring
make build-ui          # production UI → ui/dist
make serve             # kernel HTTP only :8080
```

Config docs: [`config/default.yaml`](./config/default.yaml)

---

## Status

| Area | State |
|------|-------|
| Full agent loop | Functional (`pkg/agentos`) |
| Host API for UI | `/api/v1/*` envelope + legacy `/v1/*` |
| Mission Control UI | In-repo `ui/` (Control Hub aesthetic) |
| All 10 examples | `examples/` + `aespd run-all-examples` |
| Conformance | `aespd conformance` |
| Profiles | P1 HTTP · P2 local · P3 embed |

---

## Agent loop (AESP-aligned)

1. Accept WorkUnit (0001) + capability requirements (INV-03)  
2. Plan artifact (0015) → content-addressed store (0007)  
3. Assemble Context Envelope (INV-05)  
4. Route provider + runtime by **capabilities** with failover (INV-01, INV-03)  
5. Provider complete + runtime execute  
6. Tools via unified router / MCP (INV-06)  
7. HITL never auto-approves on timeout (0014)  
8. Verify · memory trust labels · docgen  
9. Optional deploy · remediation  
10. Execution tree + event journal (INV-10)  

---

## Architecture

| Doc | Purpose |
|-----|---------|
| [docs/architecture/](./docs/architecture/) | INV-01…INV-11 |
| [docs/deployment/](./docs/deployment/) | P1 / P2 / P3 |
| [ui/README.md](./ui/README.md) | Mission Control UI |
| [examples/README.md](./examples/README.md) | Mission library |

### Invariants (kernel)

INV-01 Provider ≠ Runtime · INV-02 Plugins · INV-03 Capability routing ·  
INV-04 Unified memory · INV-05 Context envelope · INV-06 Unified tools ·  
INV-07 Credentials · INV-08 AESP is the contract · INV-09 Runtime registry ·  
INV-10 Auditable · INV-11 Host-neutral core  

---

## HTTP Host Interface

```bash
# UI contract (envelope { data, error })
curl -s http://127.0.0.1:8080/api/v1/health

# Legacy
curl -s http://127.0.0.1:8080/health
```

Vite proxies `/api` → `:8080` when running the UI in monorepo mode (`VITE_USE_MOCKS=0`).

---

## License

MIT
