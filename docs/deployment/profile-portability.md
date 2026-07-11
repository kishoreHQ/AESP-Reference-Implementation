# Profile Portability Proof

**Claim:** Phase 3 example `01-single-agent` mission.yaml runs conceptually unchanged under P1/P2/P3.

| Profile | How mission is submitted | Degradation |
|---------|--------------------------|-------------|
| P1 | Host UI → Host Interface | none if caps available |
| P2 | CLI → same mission file | if no `vision`, skip vision steps only |
| P3 | External orchestrator SDK | same |

Required capabilities in the mission file drive routing; no host-specific fields.
