# FEATURE PACK — LIVE CONNECTIONS & COMMAND DECK (v1.0)

Addendum to the UI Production Build. Extends scope with UI-PHASE 6–7 and kernel work K1–K7.

## Intent

Live cockpit for external agents as **Runtime Plugins** (`runtime.yaml`, INV-09) and model sources as **Provider Plugins** (INV-01/02). Credentials only via INV-07. Sessions journaled like missions (INV-10). **Zero vendor special-casing in UI or kernel core** — registry only.

## Kernel (K1–K7)

| ID | Surface |
|----|---------|
| K1 | Connection probe + register + handshake |
| K2 | Runtime adapters (generic-pty first, then named CLIs) |
| K3 | Interactive sessions + messages + stream events |
| K4 | Task board (kanban pull model) |
| K5 | Routines (cron) |
| K6 | Goals + journal (memory-backed) |
| K7 | Analytics rollups from event journal |

## UI

Connections wizard · Agent Rail · Live Session · Control Room · Board · Routines · Brain rail · Analytics

## Gates

- **UI-GATE-6:** Connect & See  
- **UI-GATE-7:** Command Deck  

## Security

Adapters sandboxed (PTY badged unsandboxed). Session memory default trust `agent`. Stop control always visible. Credentials write-only.
