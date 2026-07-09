# ADR-0001: Use Modular Architecture

| Metadata | Value |
|----------|-------|
| **Status** | Accepted |
| **Date** | 2025-01-15 |
| **Author(s)** | Kishore Kumar Behera, AESP Engineering Committee |
| **Deciders** | Engineering Committee |
| **Tags** | architecture, modularity, microservices |

---

## Context

The AESP Reference Implementation is a complex system comprising multiple subsystems: agent runtime, swarm orchestration, workflow management, memory services, plugin management, MCP integration, model routing, and observability. We need to decide on an overall structural approach that balances development velocity, operational simplicity, and long-term maintainability.

The primary architectural styles under consideration were:

1. **Monolithic Architecture**: All components in a single deployable unit
2. **Modular Monolith**: Single deployment with clearly separated internal modules
3. **Microservices Architecture**: Independent deployable services
4. **Modular Hybrid**: Core as a modular monolith with select services extractable

## Decision

We will adopt a **Modular Hybrid Architecture** — a modular monolith as the default deployment mode with well-defined boundaries that allow components to be extracted into independent services when operational requirements demand it.

### Key Characteristics

```
┌─────────────────────────────────────────────────────┐
│                  AESP Daemon (aespd)                 │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐  │
│  │  Agent  │ │  Swarm  │ │ Workflow│ │  Memory  │  │
│  │ Kernel  │ │ Manager │ │ Engine  │ │ Service  │  │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬─────┘  │
│       │           │           │            │        │
│       └───────────┴───────────┴────────────┘        │
│                    Internal gRPC Bus                  │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐  │
│  │  Plugin │ │   MCP   │ │  Model  │ │   Obs    │  │
│  │ Manager │ │ Gateway │ │ Router  │ │  Stack   │  │
│  └─────────┘ └─────────┘ └─────────┘ └──────────┘  │
└─────────────────────────────────────────────────────┘
         Can be extracted to standalone services
```

## Consequences

### Positive

1. **Operational Simplicity**: Default deployment is a single binary with minimal moving parts, reducing operational overhead for small-to-medium deployments.

2. **Development Velocity**: Changes across module boundaries don't require multi-service coordination during development.

3. **Type Safety**: Internal communication uses generated gRPC/Protobuf interfaces, providing compile-time contract validation.

4. **Testing Simplicity**: Integration testing is straightforward within a single process; no test containers needed for basic coverage.

5. **Deployment Flexibility**: Organizations can start with a single binary and extract services as they scale.

6. **Performance**: In-process communication avoids network overhead for the default deployment mode.

7. **Consistent Observability**: All modules share a single telemetry pipeline by default.

### Negative

1. **Deployment Granularity**: The default mode doesn't allow independent scaling of individual components.

2. **Technology Lock-in**: All modules must use the same primary language (Go) for in-process compilation.

3. **Blast Radius**: A bug in one module can potentially affect the entire process (mitigated by sandboxing for plugins).

4. **Cognitive Load**: Developers must understand module boundaries to avoid creating tight couplings.

## Mitigations

| Risk | Mitigation |
|------|-----------|
| Independent scaling | Each module has a standalone service variant; extract via configuration |
| Technology diversity | SDKs in Python/TypeScript; plugins can be external processes |
| Blast radius | Plugin sandboxing via gVisor/WebAssembly; circuit breakers between modules |
| Coupling | Automated architecture tests enforce import boundaries |

## Alternatives Considered

### Pure Microservices

**Rejected**: Operational complexity is too high for the target audience. Many users will be individual developers or small teams evaluating the system. The complexity of deploying and managing 8+ services would create a barrier to adoption.

**Reference**: "[Don't start with microservices](https://martinfowler.com/articles/dont-start-monolith.html)" — Martin Fowler

### Pure Monolith

**Rejected**: While simpler to deploy, a pure monolith would make it impossible to scale individual components independently or adopt new technologies for specific services. The observability and AI infrastructure spaces are evolving rapidly; we need the ability to evolve subsystems independently.

### Service Mesh

**Rejected**: Adding a service mesh (Istio, Linkerd) would significantly increase operational complexity. The Modular Hybrid approach achieves similar separation without requiring additional infrastructure.

## Implementation

### Module Boundaries

Each module lives in its own package under `pkg/` with the following rules:

1. **No cross-package imports** except through generated API interfaces
2. **Each module exposes** a gRPC service definition in `api/proto/`
3. **Each module has** its own configuration, tests, and documentation
4. **Internal communication** uses the generated gRPC clients/servers

### Extraction Path

When a module needs to run independently:

1. The module already has a gRPC service definition
2. A standalone `main.go` is added under `cmd/<module>/`
3. Configuration switches from `in-process` to `remote` mode
4. The in-process module is replaced with a gRPC client proxy

### Example: Extracting the Memory Service

```go
// In-process mode (default)
memService := memory.NewService(cfg.Memory)
kernel.WithMemoryService(memService)

// Remote mode (extracted)
memClient := memory.NewGRPCClient(cfg.Memory.Endpoint)
kernel.WithMemoryClient(memClient)
```

## Related Decisions

- [ADR-0002: Language Selection](0002-language-selection.md) — Primary language for module implementation
- Future ADR: Plugin Sandboxing Strategy — How external plugins are isolated
- Future ADR: Service Discovery — How extracted services discover each other

## References

1. [Modular Monolith Architecture](https://shopify.engineering/modular-monolith-rails) — Shopify Engineering
2. [The Modular Monolith: Rails Architecture](https://medium.com/@dan_manges/the-modular-monolith-rails-architecture-fb1023826fc4) — Dan Manges
3. [Building Microservices](https://samnewman.io/books/building_microservices/) — Sam Newman
4. [AESP Specification: Architecture Principles](https://github.com/kishoreHQ/AESP/blob/main/spec/architecture.md)

---

*Last updated: 2025-01-15*
