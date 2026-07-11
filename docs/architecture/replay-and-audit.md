# Replay and Audit (INV-10)

**AESP:** 0011, 0010, 0005

## Execution tree

Every mission produces an execution tree:

agents · artifacts · costs · timeline · logs · knowledge updates · evaluations · failures · replay data

## Deterministic replay

Given the mission journal + pinned digests + recorded provider/tool results (or recorded seeds),
the orchestrator MUST be able to reconstruct the decision timeline.

True bit-identical model sampling is not assumed; replay reconstructs **control decisions**
and substitutes recorded compute outputs unless live re-execution is explicitly requested.

## Audit requirements

- Who authorized what, when
- Effective policy version
- Effective provider/runtime ids (not assumed defaults)
- Tool args redacted per policy
- HITL decisions immutable once recorded
