# Tool Model (Unified Tool Layer)

**AESP:** 0015, 0013, 0010  
**Invariants:** INV-06

## Principles

1. Tools are defined **once** in the kernel tool registry.
2. All runtimes receive the **same** tool specification (schema + docs + policy hooks).
3. MCP is the preferred **tool access** interop; A2A is for **peer agents** — keep separate.
4. Every invocation produces a tool-invocation record (schema in suite `schemas/tool-invocation.json`).

## Invocation flow

```
Intent → Policy authorize → Sandbox execute → Trust-label result → Persist record → Emit aesp.tool.*
```

## Spec fields (minimum)

name, description, inputSchema, outputSchema, sideEffectClass, requiredCapabilities,
egressClass, approvalRequired (bool|policy-ref)

## Side effect classes

`read` · `write-workspace` · `write-remote` · `admin` · `destructive`
