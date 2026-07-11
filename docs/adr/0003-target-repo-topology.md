# ADR-0003: Target Multi-Repo Topology

| Metadata | Value |
|----------|-------|
| **Status** | Accepted |
| **Date** | 2026-07-11 |
| **Tags** | topology, monorepo-migration, INV-02 |

## Decision

Evolve toward these deployable units (may start as packages inside this repo):

| Unit | Responsibility |
|------|----------------|
| Kernel | Host-neutral core runtime |
| Providers | Provider plugins only |
| Runtimes | Runtime plugins only |
| Tools | Tool plugins / MCP adapters |
| Memory | Memory + KG backends |
| Workflow | Orchestration engine packaging |
| Evaluation | Eval + conformance harnesses |
| Host UI | Optional platform host (not kernel) |
| SDK | Host Interface client libraries |
| Starter-Agents | Example agents |

## Consequences

- Kernel stays free of vendor names (INV-02).
- Migration path documented in docs/deployment/migration-topology.md.
