# Memory Write Path Trust Label Rules

| Write path | Default label | May elevate privilege? |
|------------|---------------|------------------------|
| Kernel system notes | `system` | Yes (within policy) |
| Human-approved content | `verified` | Yes if policy allows |
| Agent working notes | `agent` | Conditional |
| RAG / web retrieval | `retrieved` | No by default |
| Raw tool/MCP output | `untrusted` | **No** |
| Detector flagged | `poison-suspect` | **No**; quarantine |

**Rule:** Every memory write path MUST set a trust label (INV-04).  
Missing label → write rejected (see `pkg/memory`).
