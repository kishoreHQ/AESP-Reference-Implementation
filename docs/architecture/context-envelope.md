# Context Envelope (INV-05)

**AESP:** 0004, 0006, 0013, 0015

## Definition

Every runtime receives a structured **Context Envelope**. "Prompt" is one field — not the whole context.

## Fields (normative set)

| Field | Content |
|-------|---------|
| workspace | Paths, mounts, VCS ref |
| mission | Goal, constraints, success criteria |
| memory | Selected working/semantic memories + trust labels |
| knowledge | KG subgraph / facts |
| artifacts | Referenced digests |
| policies | Obligations, denies, approval requirements |
| preferences | Non-security preferences |
| credentials | **Handles only** (never raw secrets in logs) |
| tools | Authorized tool specs |
| budget | Remaining tokens/cost/time/steps |
| security | Tenant, classification, trust posture |
| prompt | Rendered instruction view for the model |
| correlation | WorkUnit, session, trace ids |

## Assembly

Control plane + memory + KG + policy assemble the envelope before compute starts.
Runtimes MUST NOT silently expand tool or credential scope beyond the envelope.
