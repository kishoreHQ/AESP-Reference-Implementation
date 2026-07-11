# Unified Credential Model (INV-07)

**AESP:** 0013, 0015

## Flow

```
Secret store → Credential broker → scoped handle → Provider/Runtime/Tool plugin
```

## Rules

1. Core kernel never hardcodes per-vendor key env var names in business logic.
2. Plugins declare credential **requirements** (schema); broker binds secrets.
3. Handles are short-lived and auditable.
4. Raw secrets MUST NOT appear in logs, prompts (default), or provenance payloads.
5. Break-glass credentials never silently inherit to subagents.
