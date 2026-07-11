# Provider-failover workflow

**Profile portability:** P1 Platform · P2 Local-first · P3 Embedded (INV-11)  
**Required capabilities (not model names):** `coding`, `tools`  
**AESP mapping:** AESP-0015, INV-03

## Intent

≥2 abstract providers; route by capability; failover records effective provider.

## Mission declaration

See `mission.yaml`. Hosts submit via Host Interface `SubmitMission`.

## Expected event trace

See `expected-events.json`.

## Expected artifacts

See `expected-artifacts.json`.

## Run (conceptual)

```bash
# Against AESP-Reference-Implementation once wired:
aespd run --example 09-provider-failover
# Or: submit mission.yaml via SDK (P3) / CLI (P2) / Mission Control (P1)
```

## Success criteria

- Events match expected types in order (allow extra telemetry).
- No model-name routing in mission file.
- Works under all three profiles degrading only by missing capabilities.
