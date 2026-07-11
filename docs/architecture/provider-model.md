# Provider Model (INV-01)

**AESP:** 0015  
**Invariants:** INV-01, INV-02, INV-03, INV-07

## Definition

A **provider** supplies model inference capabilities.  
A **runtime** executes agent work. They are separate registries and lifecycles.

Core kernel code MUST contain **zero vendor names**. Vendor support = **provider plugins**.

## Capability advertisement

Providers declare capabilities, not just model strings:

`reasoning` · `coding` · `tools` · `streaming` · `thinking` · `files` · `vision` ·
`computer-use` · `embeddings` · `local` · `batch` · …

## Routing flow (INV-03)

```
Intent → Capability Engine → Routing Policy → Provider Selection → Model Selection
       → Runtime Selection → Validation → Memory → Artifacts
```

Never `if model == <vendor-string>` in kernel.

## Provider plugin contract

| Method | Purpose |
|--------|---------|
| Describe() | capabilities, limits, cost hints |
| Complete(ctx, request) | chat/completion |
| Embed(ctx, request) | optional |
| Health(ctx) | readiness |

Credentials arrive only via credential broker (INV-07).

## Fallback

Fallback chains are policy objects. Effective provider/model MUST be recorded on every call
(provider invocation metadata).
