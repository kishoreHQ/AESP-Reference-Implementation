# Artifact Model

**AESP:** 0007, 0009, 0010, 0015 (plan artifacts)

## Identity

Artifacts are content-addressed by **digest** (e.g., sha256).  
Mutable names MAY alias digests; production promotion MUST pin digests.

## Lifecycle

```
draft → validated → approved → published → deprecated → archived
                ↘ rejected
```

## Provenance (MUST for production)

- producing agent / runtime / provider metadata
- WorkUnit + session ids
- input digests
- tool invocation ids
- trust label
- signature or attestation when available

## Plan artifacts

Versioned plans (INT-REQ-075/076): goal, steps, assumptions, successCriteria, revision.
