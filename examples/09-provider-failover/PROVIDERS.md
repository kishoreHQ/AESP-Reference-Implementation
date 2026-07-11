# Provider failover demonstration

This example uses **abstract** providers only:

| Plugin ID | Capabilities | Role |
|-----------|--------------|------|
| `provider.mock-remote` | coding, tools, vision, reasoning | Primary (simulated unhealthy) |
| `provider.mock-local` | coding, tools, local, reasoning | Failover |

Routing is **capability-based** (`coding` + `tools`). Kernel never branches on vendor model names.
Effective provider id MUST appear on the completion / route event (AESP-0015).
