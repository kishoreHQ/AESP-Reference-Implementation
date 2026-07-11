# Host Interface (INV-11)

**AESP:** 0014, 0015, 0011

## Principle

The kernel MUST NOT assume any specific host, UI, or environment.
All host interaction goes through this interface.

## Operations

| Operation | Direction | Purpose |
|-----------|-----------|---------|
| SubmitMission | Host → Kernel | Create mission / WorkUnit |
| CancelMission | Host → Kernel | Abort |
| SubscribeEvents | Host ⇄ Kernel | Stream mission events |
| RequestApproval / ResolveApproval | Kernel ⇄ Host | HITL |
| GetArtifact | Host → Kernel | Retrieve by digest/name |
| GetExecutionTree | Host → Kernel | Agents, costs, timeline, failures |
| QueryMemory / QueryKG | Host → Kernel | Optional read APIs |
| Health | Host → Kernel | Liveness |

## Deployment profiles

| Profile | Host examples (non-normative) |
|---------|-------------------------------|
| P1 Platform | Web mission control, multi-user server |
| P2 Local-first | CLI/TUI on laptop |
| P3 Embedded | Third-party orchestrator embeds SDK |

## SDK surface

`pkg/host` defines interfaces; language SDKs wrap gRPC/HTTP/in-process bindings.
