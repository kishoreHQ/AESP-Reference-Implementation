# Source Code Organization (`src/`)

> **Note**: This project follows Go project conventions. The `src/` directory is reserved for application-specific source code that is not intended for external consumption as library packages.

---

## Overview

The `src/` directory contains the main application source code for the AESP Agent Operating System. This includes the core daemon, service implementations, and internal utilities that are not exposed as public APIs.

In Go projects, `src/` is used alongside `pkg/` (public library code) and `cmd/` (executable entry points) to create a clear separation of concerns.

## Directory Structure

```
src/
├── daemon/                 # Core daemon implementation
│   ├── server.go          # Main server orchestration
│   ├── config.go          # Configuration loading and validation
│   ├── lifecycle.go       # Startup, shutdown, and signal handling
│   └── health.go          # Health check endpoints
│
├── agent/                  # Agent runtime implementation
│   ├── runtime.go         # Agent execution runtime
│   ├── sandbox.go         # Sandbox/isolation for agent execution
│   ├── context.go         # Agent execution context management
│   ├── capabilities.go    # Capability registration and discovery
│   └── hooks.go           # Lifecycle hooks (pre-start, post-stop, etc.)
│
├── swarm/                  # Swarm orchestration implementation
│   ├── orchestrator.go    # Swarm lifecycle management
│   ├── messaging.go       # Inter-agent messaging system
│   ├── discovery.go       # Agent discovery and registration
│   ├── consensus.go       # Distributed consensus for decisions
│   └── loadbalancer.go    # Work distribution across agents
│
├── workflow/               # Workflow engine implementation
│   ├── executor.go        # DAG execution engine
│   ├── scheduler.go       # Task scheduling and dependency resolution
│   ├── checkpoint.go      # Workflow state checkpointing
│   ├── retry.go           # Retry logic with backoff strategies
│   └── parser.go          # Workflow definition parsing (YAML/JSON)
│
├── memory/                 # Memory service implementation
│   ├── shortterm.go       # Short-term/session memory (Redis)
│   ├── longterm.go        # Long-term persistent memory (PostgreSQL)
│   ├── vector.go          # Vector storage for semantic search (pgvector)
│   ├── objectstore.go     # Object storage for artifacts (S3/MinIO)
│   └── checkpoint.go      # Agent state checkpoint/restore
│
├── plugin/                 # Plugin system implementation
│   ├── loader.go          # Dynamic plugin loading
│   ├── registry.go        # Plugin registry and discovery
│   ├── sandbox.go         # Plugin sandboxing (gVisor/WASM)
│   ├── lifecycle.go       # Plugin lifecycle management
│   └── security.go        # Plugin permission and security model
│
├── mcp/                    # MCP gateway implementation
│   ├── client.go          # MCP client implementation
│   ├── server.go          # MCP server implementation
│   ├── tools.go           # Tool registration and discovery
│   ├── context.go         # Context management for MCP
│   └── translator.go      # Protocol translation (AESP <> MCP)
│
├── router/                 # Model router implementation
│   ├── provider.go        # LLM provider abstraction
│   ├── routing.go         # Routing strategies
│   ├── fallback.go        # Fallback and circuit breaker logic
│   ├── cache.go           # Response caching
│   ├── ratelimit.go       # Rate limiting per provider
│   └── costtracker.go     # Cost tracking and optimization
│
├── observability/          # Observability implementation
│   ├── traces.go          # OpenTelemetry trace configuration
│   ├── metrics.go         # Prometheus metrics registration
│   ├── logging.go         # Structured logging setup
│   ├── events.go          # Event bus implementation
│   └── exporters.go       # OTLP, Prometheus, file exporters
│
├── api/                    # Internal API layer
│   ├── middleware/        # HTTP/gRPC middleware
│   │   ├── auth.go       # Authentication middleware
│   │   ├── logging.go    # Request logging
│   │   ├── metrics.go    # Request metrics
│   │   └── recovery.go   # Panic recovery
│   ├── handlers/          # HTTP request handlers
│   └── grpc/             # gRPC service implementations
│
├── auth/                   # Authentication and authorization
│   ├── authenticator.go   # Authentication strategies
│   ├── authorizer.go      # Authorization (RBAC/ABAC)
│   ├── rbac.go           # Role-based access control
│   ├── token.go          # JWT token management
│   └── policy.go         # OPA policy integration
│
└── utils/                  # Internal utilities
    ├── crypto.go          # Cryptographic helpers
    ├── validation.go      # Input validation
    ├── errors.go          # Error handling utilities
    └── backoff.go         # Retry/backoff algorithms
```

## Design Principles

1. **Package Cohesion**: Each directory represents a single, cohesive concept with a well-defined responsibility.

2. **Interface-Driven**: Components depend on interfaces, not concrete implementations. Interfaces are defined in the consuming packages (Go idiom).

3. **Dependency Direction**: Dependencies flow inward. `src/` packages may depend on `pkg/` packages, but not vice versa:
   ```
   cmd/  -->  src/  -->  pkg/
   ```

4. **No Circular Dependencies**: The build will fail if circular imports are introduced. This is enforced by CI.

5. **Internal Visibility**: All packages under `src/` are treated as internal. They should not be imported by external projects. Public APIs live in `pkg/`.

## Key Components

### Daemon (`src/daemon/`)

The main server that orchestrates all services. It handles:
- Configuration loading from files, environment variables, and flags
- Service initialization and dependency injection
- Graceful startup and shutdown
- Health check and readiness probes
- Signal handling (SIGTERM, SIGINT, SIGHUP for config reload)

### Agent Runtime (`src/agent/`)

The core agent execution engine. Key files:
- `runtime.go`: Manages goroutine-per-agent execution model
- `sandbox.go`: Provides isolation boundaries for untrusted agent code
- `context.go`: Propagates execution context (deadlines, cancellation, values)

### Workflow Executor (`src/workflow/`)

The DAG execution engine. Key design decisions:
- Topological sort for dependency resolution
- Concurrent execution of independent tasks
- Checkpointing after each task for fault tolerance
- Event-driven status updates

## Testing

Each package in `src/` has comprehensive tests:

```
src/agent/
├── runtime.go
├── runtime_test.go       # Unit tests
├── runtime_fuzz_test.go  # Fuzz tests
├── runtime_bench_test.go # Benchmarks
└── testdata/             # Test fixtures
```

Run tests for a specific package:
```bash
go test ./src/agent/...
go test ./src/swarm/...
go test ./src/workflow/...
```

## Contributing

When adding new code to `src/`:

1. Place it in the appropriate package based on responsibility
2. Follow the interface-driven design pattern
3. Add comprehensive tests with `_test.go` files
4. Ensure no circular dependencies are introduced
5. Keep packages focused — if a package grows too large, consider splitting
6. Document public-facing types and functions with GoDoc comments

## See Also

- [`pkg/README.md`](../pkg/README.md) — Public library packages
- [`cmd/README.md`](../cmd/README.md) — Executable entry points
- [`docs/architecture.md`](../docs/architecture.md) — System architecture overview
