# Runtime Loop — Agent Harness Semantics

**AESP:** AGENT-RUNTIME.md, 0001, 0004, 0005, 0010, 0011, 0013, 0014, 0015  
**Invariants:** INV-05, INV-06, INV-10

## 1. Canonical agent loop

```
Accept WorkUnit + authz
        ↓
Load Context Envelope (memory, KG, artifacts, policy, tools, budget, security)
        ↓
Plan artifact (versioned; INT-REQ-075/076)
        ↓
Act: Runtime Plugin + Provider Plugin + Tools under policy
     record tool invocations; emit spans
        ↓
Verify: tests / validators / HITL
        ↓
Persist: memory/KG (trust labels), evidence, audit, replay journal
        ↓
Stop when success | budget | escalation | hard fail
```

## 2. Stopping conditions (MUST)

A runtime MUST define for each WorkUnit:

| Condition | Source |
|-----------|--------|
| Success criteria | Plan / task contract |
| Max steps | Quotas (0001) |
| Max tokens / cost | Budget (0015) |
| Max wall time | Session TTL |
| Escalation path | HITL (0014) |
| Failure taxonomy | Structured errors |

Unbounded production loops are **non-conformant**.

## 3. Control vs compute rule

Compute workers MUST NOT self-approve production side effects (0013, 0014).
Approvals are control-plane only.

## 4. Subagent isolation

| Rule | Requirement |
|------|-------------|
| Distinct principal or scoped capability set | 0001, 0002 |
| Context isolation parent/child | 0004 |
| Tool allowlists per subagent | 0013, 0015 |
| Parent monitors child WorkUnits | 0005 |
| No silent inheritance of break-glass credentials | 0013 |

## 5. Correlation keys (every step)

WorkUnit id · workflow instance · task id · session id · trace id · artifact digests · HITL task id

## 6. AESP mapping

| Loop phase | Specs |
|------------|-------|
| Accept | 0001, 0002, 0013 |
| Context | 0004, 0006, INV-05 |
| Plan | 0015 INT-REQ-075+ |
| Act | 0015 tools/providers, runtime plugins |
| Verify | 0010, 0014 |
| Persist | 0004, 0006, 0011, INV-10 |
