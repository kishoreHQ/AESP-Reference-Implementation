# Policy Model

**AESP:** 0013, 0002, 0014, 0015

## Policy classes

| Class | Examples |
|-------|----------|
| Admission | Who may submit missions |
| Tool authorization | Tool + args allow/deny |
| Provider routing | Capability constraints, cost caps |
| Data classification | What may enter prompts |
| HITL gates | When human approval is required |
| Egress | Network destinations |
| Freeze | Change freezes, blast radius |

## Evaluation order

1. Hard deny (security)
2. Freeze / maintenance
3. Tenant isolation
4. Role permissions
5. Resource budgets
6. Routing preferences
7. Soft recommendations

Fail closed on missing policy for production side effects.

## Engine interface

`pkg/policy` — Evaluate(DecisionRequest) → Allow|Deny|RequireApproval + obligations.
