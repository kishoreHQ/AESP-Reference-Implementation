# Agent Lifecycle

**AESP:** 0001, 0002, 0003, 0013

## States

```
registered → admitted → active → suspended → terminated
                 ↘ rejected
active → degraded → active | terminated
```

| State | Meaning |
|-------|---------|
| registered | Principal exists in agent registry |
| admitted | Policy allows work under current org/tenant |
| active | May receive WorkUnits |
| suspended | Temporary hold (budget, incident, human) |
| degraded | Limited tools/capabilities |
| terminated | No further work; audit retained |
| rejected | Admission denied |

## Transitions MUST be audited

Every transition emits `aesp.agent.lifecycle.*` with principal id, reason, actor, correlation ids.

## Roles

Role templates and permission boundaries: AESP-0002. Kernel stores role bindings;
evaluation is policy-engine concern.
