# Trust Model

**AESP:** 0013, 0004, 0015  
**Invariants:** INV-04, INV-06

## Trust labels (memory & tool results)

| Label | Meaning | May authorize privileged actions? |
|-------|---------|-----------------------------------|
| `system` | Kernel-authored | Yes (within policy) |
| `verified` | Passed validation / human / test | Yes if policy allows |
| `agent` | Agent-written working memory | Conditional |
| `retrieved` | External retrieval (RAG/web) | No by default |
| `untrusted` | Tool/MCP/unvalidated | **No** |
| `poison-suspect` | Detection flagged | **No**; quarantine |

## Boundaries

1. **Host ↔ Kernel** — authenticated Host Interface
2. **Control ↔ Compute** — no self-approval
3. **Kernel ↔ Provider plugin** — credentials scoped, no policy mutation
4. **Kernel ↔ Runtime plugin** — sandboxed capability set
5. **Kernel ↔ Tool / MCP** — pre-authorization + invocation record
6. **Agent ↔ Agent (A2A)** — peer identity + task contract
7. **Memory write paths** — every write gets a trust label (INV-04)

## Memory poison / IPI

Untrusted content MUST NOT be executable as instructions without explicit elevation
(human or policy). See AESP-0013 memory poison requirements.
