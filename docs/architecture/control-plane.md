# Control Plane

**AESP:** 0001, 0002, 0005, 0013, 0014, 0015  
**Invariants:** INV-03, INV-07, INV-10, INV-11

## 1. Owns

- Admission of missions / WorkUnits
- Role and permission evaluation
- Budget and quota enforcement
- Routing **policy** (capability → provider/runtime selection rules)
- HITL task lifecycle and approval authority
- Session authority and revocation
- Policy-as-code evaluation
- Freeze windows and blast-radius caps (with remediation)

## 2. Does not own

- Model token sampling
- Sandboxed tool process execution
- Vendor SDK transport details

## 3. Components

| Component | Interface package | AESP |
|-----------|-------------------|------|
| Orchestrator | `pkg/orchestrator` | 0005 |
| Policy engine | `pkg/policy` | 0013 |
| Approval workflow | `pkg/approval` | 0014 |
| Capability engine | `pkg/capability` | INV-03, 0015 |
| Router (policy-driven) | `pkg/router` | INV-03 |
| Credential broker | `pkg/credentials` | INV-07 |
| Host callbacks | `pkg/host` | INV-11 |

## 4. Authority rule

Any action with production side effects MUST pass control-plane authorization before
compute-plane execution. Fail closed.

## 5. Events

`aesp.control.mission.accepted` · `aesp.control.policy.denied` · `aesp.control.budget.exceeded` ·
`aesp.control.approval.requested` · `aesp.control.approval.resolved` · `aesp.control.session.revoked`
