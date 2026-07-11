# P3 — Embedded Deployment Profile

**INV-11** · Kernel as library/sidecar inside external orchestrator

## Host Interface (SDK surface)

| Method | Purpose |
|--------|---------|
| SubmitMission | Start work |
| CancelMission | Abort |
| SubscribeEvents | Stream progress |
| ResolveApproval | HITL |
| GetArtifact | Fetch digests |
| GetExecutionTree | Audit view |
| Health | Liveness |

Go interface: `pkg/host.Interface`  
In-process: embed `pkg/kernel.Kernel`  
Sidecar: expose gRPC/HTTP adapters (future package `pkg/transport`)

## Rule

External hosts MUST NOT reach into provider SDKs or memory engines directly;
all access via Host Interface / published query APIs.
