# Single-agent task

**Profile portability:** P1 Platform · P2 Local-first · P3 Embedded (INV-11)  
**Required capabilities (not model names):** `coding`, `tools`  
**AESP mapping:** AESP-0001, AESP-0005, AESP-0015

## Intent

Template for all examples: one agent completes a WorkUnit.

## Mission declaration

See `mission.yaml`. Hosts submit via Host Interface `SubmitMission`.

## Expected event trace

See `expected-events.json`.

## Expected artifacts

See `expected-artifacts.json`.

## Run (conceptual)

```bash
# Against AESP-Reference-Implementation once wired:
aespd run --example 01-single-agent
# Or: submit mission.yaml via SDK (P3) / CLI (P2) / Mission Control (P1)
```

## Success criteria

- Events match expected types in order (allow extra telemetry).
- No model-name routing in mission file.
- Works under all three profiles degrading only by missing capabilities.
