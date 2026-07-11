# Compute Plane

**AESP:** 0015, 0013, AGENT-RUNTIME  
**Invariants:** INV-01, INV-02, INV-09

## 1. Owns

- Provider plugin invocations (model completions)
- Runtime plugin execution (agent harness implementations)
- Sandboxed tool process execution
- Resource metering feedback to control plane

## 2. Isolation requirements

| Boundary | Mechanism |
|----------|-----------|
| Tool process | OS sandbox / container / WASM as configured |
| Runtime plugin | Capability allowlist + network egress policy |
| Provider plugin | Credential injection only via credential broker |
| Untrusted tool output | Trust label `untrusted` until validated |

## 3. Separation from control

Compute MUST NOT:

- Grant itself elevated credentials
- Auto-approve HITL tasks
- Mutate policy documents
- Bypass tool authorization

## 4. Plugin surfaces

- **Provider plugins** — see provider-model.md
- **Runtime plugins** — see runtime-model.md (`runtime.yaml` manifests)
- **Tool plugins** — see tool-model.md
