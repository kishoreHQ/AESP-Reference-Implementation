# Migration Path — Current 3 Repos → Target Topology

## Current

| Repo | Role |
|------|------|
| AESP | Protocol suite (contract) |
| AESP-Reference-Implementation | Kernel + plugins monorepo start |
| AESP-Examples | Mission library |

## Target units (ADR-0003)

Kernel · Providers · Runtimes · Tools · Memory · Workflow · Evaluation · Host UI · SDK · Starter-Agents

## Steps

1. Keep AESP suite as normative source of truth (INV-08).  
2. Grow packages under this repo until plugin boundaries stabilize.  
3. Extract `plugins/providers` → Providers repo when third-party plugins exist.  
4. Extract Host UI to separate product repo; depends only on SDK.  
5. Examples remain profile-portable missions consuming SDK + Kernel releases.

No breaking AESP requirement renumbering during extraction.
