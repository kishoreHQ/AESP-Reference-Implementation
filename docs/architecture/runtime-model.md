# Runtime Model (INV-01, INV-09)

**AESP:** AGENT-RUNTIME, 0001, 0005, 0015  
**Invariants:** INV-01, INV-02, INV-09

## Definition

A **runtime** is a pluggable agent harness that executes WorkUnits given a Context Envelope.
Examples of *categories* (not kernel hardcodes): coding CLI harness, multi-agent SDK runner,
notebook executor. Each ships as a **runtime plugin**.

## Discovery: `runtime.yaml`

```yaml
apiVersion: aesp.runtime/v1
kind: RuntimePlugin
metadata:
  id: example.generic-loop
  version: 1.0.0
spec:
  capabilitiesIn: [tools, streaming]
  capabilitiesOut: [coding, planning]
  sandbox: process
  entrypoint: ./bin/runtime
  configSchema: ./schema.json
```

Adding a runtime MUST require **zero kernel source changes** (INV-09).

## Registry lifecycle

discover → validate manifest → register → load → health → unload → deprecate

## Pairing rule (INV-01)

Any runtime MUST be pairable with any **capability-compatible** provider.
The kernel pairs them via Capability Engine + Routing Policy — not hardcoded couples.

## Isolation

Runtime plugins receive:

- Context Envelope (INV-05)
- Tool specs (INV-06) — not raw credentials for tools
- Scoped provider handle (not raw API keys; INV-07)
- Budget and policy obligations
