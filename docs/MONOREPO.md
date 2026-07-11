# Monorepo merge notes

**Date:** 2026-07-11  
**Decision:** Product code for Agent OS lives in **one repo**. AESP protocol stays separate.

## Merged into this repository

| Source | Destination |
|--------|-------------|
| `kishoreHQ/hermes-mission-control-ui` | `ui/` |
| `kishoreHQ/AESP-Examples` (missions) | `examples/` (already present; re-synced) |
| This repo kernel | `cmd/`, `pkg/`, `plugins/` |

## Not merged

| Repo | Reason |
|------|--------|
| `kishoreHQ/AESP` | Vendor-neutral protocol suite (INV-08). Implementation must not own the spec. |

## How to run

```bash
./scripts/dev.sh     # or: make dev
```

## Historical remotes

Standalone clones may remain on GitHub for discoverability; **active development is this monorepo.**  
Prefer opening PRs against `AESP-Reference-Implementation`.
