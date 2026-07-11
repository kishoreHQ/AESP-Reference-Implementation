# Threat Model — Agent OS Kernel

**Status:** GATE-4 artifact  
**Date:** 2026-07-11  
**Boundaries:** From Phase 1 trust-model & host-interface

## Assets

| Asset | Impact if compromised |
|-------|----------------------|
| Credentials / secret handles | Provider & tool abuse |
| Memory store | Poisoned context → mis-action |
| Policy engine | Privilege escalation |
| Artifact digests | Supply-chain substitution |
| HITL decisions | Unauthorized production change |
| Execution journal | Audit tampering |

## Trust boundaries & threats

| Boundary | Threats | Controls |
|----------|---------|----------|
| Host ↔ Kernel | Spoofed missions, event injection | Authenticated Host Interface; tenant scoping |
| Control ↔ Compute | Self-approval, credential theft | Separation of planes; no self-HITL; scoped handles |
| Kernel ↔ Provider plugin | Prompt exfil, cost abuse | Budget caps; redaction; capability routing |
| Kernel ↔ Runtime plugin | Sandbox escape | Process/container sandbox; capability allowlists |
| Kernel ↔ Tool/MCP | Tool poisoning, confused deputy | Pre-authz; invocation records; trust labels on results |
| Agent ↔ Agent (A2A) | Peer impersonation | Peer identity; task contracts |
| Memory write paths | Memory poison / IPI | Trust labels; untrusted cannot authorize privilege |
| Retrieval → Prompt | Indirect prompt injection | Isolate retrieved content; label `retrieved`/`untrusted` |

## STRIDE summary

- **Spoofing:** principal & peer identity required  
- **Tampering:** content-addressed artifacts; append-only journals  
- **Repudiation:** audit events with actors  
- **Info disclosure:** secret redaction; classification on envelope  
- **DoS:** budgets, max steps, circuit breakers (remediation)  
- **Elevation:** policy fail-closed; HITL for destructive/admin  

## Residual risks

- Full MCP/A2A golden attack suites: GAP-001/002  
- Multi-tenant hard isolation on shared hardware: profile-dependent  
