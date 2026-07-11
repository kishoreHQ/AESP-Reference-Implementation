# Failure Mode Catalog & Recovery Flows

| Failure | Detection | Recovery | Escalate |
|---------|-----------|----------|----------|
| Provider unhealthy | Health probe / error rate | Failover to next capability-compatible provider; record effective id | If none left → HITL or fail mission |
| Runtime crash | Non-zero exit / timeout | Retry with same envelope (bounded); then alternate runtime | HITL if side effects uncertain |
| Tool denied | Policy Deny | Surface to plan revision or HITL | Always audit |
| Tool timeout | Deadline | Cancel sandbox; mark step failed; retry policy | — |
| Budget exceeded | Metering | Stop loop; emit `aesp.control.budget.exceeded` | Optional HITL to raise budget |
| HITL timeout | Timer | **Expire without approve** | Notify host |
| Memory poison suspect | Detector / policy | Quarantine label; exclude from privileged paths | Security review |
| Artifact digest mismatch | Verify | Reject promotion; rollback deploy session | Incident |
| Orchestration step fail | Step status | Retry → compensate/rollback → fail | Playbook (0012) |
| Event bus backpressure | Queue depth | Drop non-critical telemetry only; never drop audit | Ops |

Every mode MUST leave a journal event (INV-10).
