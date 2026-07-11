# Unified Memory Model (INV-04)

**AESP:** 0004, 0006, 0013

## Rule

One memory subsystem for the OS. Every runtime reads and writes it.
**No per-runtime memory silos.**

## Stores (logical)

| Store | Purpose |
|-------|---------|
| Working | Short-horizon task state |
| Session | Mission/session continuity |
| Semantic / vector | Retrieval |
| Knowledge graph | Structured entities/relations (0006) |
| Artifacts | Content-addressed blobs |
| Evaluations | Scores, promotions |

## Trust

Every write carries a trust label (trust-model.md). Retrieval injects labels into the Context Envelope.

## Backends by profile

| Profile | Typical backend |
|---------|-----------------|
| P1 | Postgres + object store + vector |
| P2 | SQLite + filesystem |
| P3 | Host-provided or embedded SQLite |
