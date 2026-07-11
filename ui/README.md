# Mission Control UI

Host Interface client for the AESP Agent OS.

**This package lives in the monorepo** at `ui/`. Prefer starting from the repo root:

```bash
# from AESP-Reference-Implementation/
./scripts/dev.sh              # kernel :8080 + UI :5173
./scripts/dev-ui-only.sh      # UI with mocks only
```

## Standalone (from this folder)

```bash
npm install
VITE_USE_MOCKS=1 npm run dev -- --host 127.0.0.1 --port 5173
```

Live against kernel (kernel must already be on :8080):

```bash
VITE_USE_MOCKS=0 npm run dev -- --host 127.0.0.1 --port 5173
```

Vite proxies `/api` → `http://127.0.0.1:8080` (see `vite.config.ts`).

## Design

Control Hub / Cherenkov neon ops console. Mission Spine, dashboard KPIs, HITL approvals.
