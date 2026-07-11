# AESP Examples — Executable Mission Library

Ten profile-portable examples proving the AESP / Agent OS contract.

| # | Directory | Capabilities |
|---|-----------|--------------|
| 01 | [01-single-agent](./01-single-agent/) | coding, tools |
| 02 | [02-multi-agent](./02-multi-agent/) | coding, tools, planning |
| 03 | [03-code-generation](./03-code-generation/) | coding, tools |
| 04 | [04-review-approval](./04-review-approval/) | coding |
| 05 | [05-memory-update](./05-memory-update/) | reasoning, tools |
| 06 | [06-kg-update](./06-kg-update/) | reasoning |
| 07 | [07-remediation](./07-remediation/) | tools, reasoning |
| 08 | [08-hitl](./08-hitl/) | tools |
| 09 | [09-provider-failover](./09-provider-failover/) | coding, tools (≥2 providers) |
| 10 | [10-rollback-retry](./10-rollback-retry/) | tools |

## Rules

1. Declare **capabilities**, never model vendor names.
2. Map to AESP requirement IDs in each README.
3. Ship expected event trace + artifacts.
4. Same mission MUST run under Platform / Local-first / Embedded (INV-11).
