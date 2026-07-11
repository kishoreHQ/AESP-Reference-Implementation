# P1 — Platform Deployment Profile

**INV-11** · Server / multi-tenant Agent OS + host UI

## Bootstrap path

1. Deploy kernel binary `aespd` with durable store (Postgres) + object store.
2. Configure credential broker with org secrets (not in env of runtimes).
3. Register provider plugins and runtime plugins (`runtime.yaml` discovery).
4. Attach Host UI (Mission Control class) via Host Interface only.
5. Enable OBS exporters; HITL callbacks to host.

## Surfaces

mission → workflow → execution tree → agents → artifacts → costs → timeline →
logs → knowledge updates → evaluations → failures → replay

## Checklist

- [ ] Multi-tenant isolation  
- [ ] HITL never auto-approves on timeout  
- [ ] Capability routing only  
- [ ] Audit journal durable  
