# Public Library Packages (`pkg/`)

> **Purpose**: The `pkg/` directory contains public library code that is intended for external consumption. These packages form the **AESP Go SDK** and may be imported by other projects.

---

## Overview

Packages in `pkg/` define the public APIs, interfaces, and reusable libraries of the AESP Agent Operating System. This code is designed to be imported by:

- External applications building on AESP
- Plugin developers writing Go plugins
- Third-party integrations
- The AESP CLI and other tools in `cmd/`

## Guiding Principle

> **Once a package is in `pkg/`, its public API becomes a compatibility commitment.** Changes to exported APIs must follow semantic versioning and deprecation policies.

## Directory Structure

```
pkg/
├── api/                    # Shared API types and generated code
│   ├── types/             # Common type definitions (Agent, Task, Swarm, etc.)
│   ├── errors/            # Standardized error types and codes
│   ├── events/            # Event type definitions
│   └── version.go         # Version information
│
├── agent/                  # Agent SDK and runtime interfaces
│   ├── agent.go           # Agent type definition and constructors
│   ├── config.go          # Agent configuration types
│   ├── capabilities.go    # Capability types and registry interface
│   ├── interfaces.go      # Core agent interfaces (Runtime, Executor, etc.)
│   └── client.go          # Agent gRPC client for remote interaction
│
├── kernel/                 # Agent Kernel public interfaces
│   ├── kernel.go          # Kernel interface definition
│   ├── config.go          # Kernel configuration
│   ├── lifecycle.go       # Agent lifecycle event types
│   └── options.go         # Kernel option pattern (functional options)
│
├── swarm/                  # Swarm management SDK
│   ├── swarm.go           # Swarm type and interfaces
│   ├── config.go          # Swarm configuration
│   ├── messaging.go       # Message types and messaging interfaces
│   ├── discovery.go       # Discovery service interfaces
│   └── client.go          # Swarm Manager gRPC client
│
├── memory/                 # Memory service SDK
│   ├── memory.go          # Memory service interface
│   ├── types.go           # Memory-related types (Context, Checkpoint, etc.)
│   ├── store.go           # Storage abstraction interfaces
│   ├── vector.go          # Vector search interfaces
│   └── client.go          # Memory Service gRPC client
│
├── workflow/               # Workflow engine SDK
│   ├── workflow.go        # Workflow type and interfaces
│   ├── parser.go          # Workflow definition parser interface
│   ├── executor.go        # Workflow executor interface
│   ├── tasks.go           # Task type definitions
│   └── client.go          # Workflow Engine gRPC client
│
├── plugin/                 # Plugin system SDK
│   ├── plugin.go          # Plugin type and interfaces
│   ├── registry.go        # Plugin registry interface
│   ├── loader.go          # Plugin loader interface
│   ├── sandbox.go         # Sandbox interface
│   └── host.go            # Plugin host interface (for plugin authors)
│
├── mcp/                    # MCP gateway SDK
│   ├── mcp.go             # MCP client/server interfaces
│   ├── tools.go           # Tool definition types
│   ├── context.go         # MCP context types
│   └── types.go           # MCP protocol types
│
├── router/                 # Model router SDK
│   ├── router.go          # Model router interface
│   ├── provider.go        # LLM provider interface
│   ├── config.go          # Provider configuration types
│   ├── strategies.go      # Routing strategy types
│   └── client.go          # Model Router gRPC client
│
├── observability/          # Observability SDK
│   ├── telemetry.go       # Telemetry configuration and setup
│   ├── traces.go          # Tracing helpers
│   ├── metrics.go         # Metrics registration helpers
│   ├── logging.go         # Structured logging utilities
│   └── events.go          # Event publishing/subscription
│
├── crypto/                 # Cryptographic utilities
│   ├── keys.go            # Key generation and management
│   ├── tokens.go          # Secure token generation
│   ├── hash.go            # Hashing utilities
│   └── encrypt.go         # Encryption/decryption helpers
│
└── testutil/               # Testing utilities
    ├── fixtures.go        # Common test fixtures
    ├── mocks/             # Generated mocks for interfaces
    ├── fakes/             # Fake implementations for testing
    └── helpers.go         # Test helper functions
```

## Interface Design Principles

### 1. Consumer-Defined Interfaces

Following Go best practices, interfaces are defined by the packages that consume them, not the packages that implement them:

```go
// pkg/agent/interfaces.go — defined by the agent package
package agent

// Runtime is the execution environment for an agent.
// Implementations may be in-process or remote.
type Runtime interface {
    Execute(ctx context.Context, task Task) (*Result, error)
    Status(ctx context.Context) (Status, error)
    Stop(ctx context.Context) error
}
```

### 2. Functional Options Pattern

Configuration uses the functional options pattern for flexibility:

```go
// pkg/kernel/options.go
package kernel

type Option func(*config)

func WithModelProvider(provider string) Option {
    return func(c *config) {
        c.modelProvider = provider
    }
}

func WithMemoryEndpoint(endpoint string) Option {
    return func(c *config) {
        c.memoryEndpoint = endpoint
    }
}

// Usage:
kernel, err := kernel.New(
    kernel.WithModelProvider("openai"),
    kernel.WithMemoryEndpoint("localhost:50051"),
)
```

### 3. Context-First API

All operations accept `context.Context` as the first parameter:

```go
// Correct
func (c *Client) GetAgent(ctx context.Context, id string) (*Agent, error)

// Incorrect — missing context
func (c *Client) GetAgent(id string) (*Agent, error)
```

### 4. Error Types

Standardized errors allow callers to handle specific error conditions:

```go
// pkg/api/errors/errors.go
package errors

var (
    ErrAgentNotFound     = errors.New("agent not found")
    ErrSwarmNotFound     = errors.New("swarm not found")
    ErrTaskFailed        = errors.New("task execution failed")
    ErrProviderUnavailable = errors.New("model provider unavailable")
    ErrRateLimited       = errors.New("rate limited")
    ErrUnauthorized      = errors.New("unauthorized")
    ErrTimeout           = errors.New("operation timed out")
)

// With details:
func NewAgentNotFound(id string) error {
    return fmt.Errorf("%w: %s", ErrAgentNotFound, id)
}
```

## Version Compatibility

Packages in `pkg/` follow semantic versioning:

| API Change | Version Impact |
|-----------|---------------|
| Add new exported symbol | Minor (`v0.X.0`) |
| Deprecate exported symbol | Minor (`v0.X.0`) |
| Remove exported symbol | Major (`vX.0.0`) |
| Change exported symbol behavior | Major (`vX.0.0`) |
| Fix bug (no API change) | Patch (`v0.0.X`) |

## Generated Code

The `pkg/` directory contains generated code from Protocol Buffer definitions:

```
pkg/api/
├── types/
│   └── agent.pb.go       # Generated from api/proto/agent.proto
│   └── task.pb.go        # Generated from api/proto/task.proto
│   └── swarm.pb.go       # Generated from api/proto/swarm.proto
```

Generated code is committed to the repository to ensure reproducible builds. Regenerate with:

```bash
make generate
```

## Mocks

Interface mocks are generated using `mockery` or `gomock` and live in `pkg/testutil/mocks/`:

```go
// pkg/testutil/mocks/mock_agent_runtime.go
package mocks

type MockRuntime struct {
    mock.Mock
}

func (m *MockRuntime) Execute(ctx context.Context, task agent.Task) (*agent.Result, error) {
    args := m.Called(ctx, task)
    return args.Get(0).(*agent.Result), args.Error(1)
}
```

Generate mocks:
```bash
make mocks
```

## Usage Examples

### Creating an Agent Client

```go
import (
    "github.com/kishoreHQ/aesp/pkg/agent"
    "github.com/kishoreHQ/aesp/pkg/kernel"
)

// In-process kernel
k, err := kernel.New(kernel.WithModelProvider("openai"))
if err != nil {
    log.Fatal(err)
}

a, err := k.CreateAgent(ctx, agent.Config{
    Name:         "code-reviewer",
    Capabilities: []string{"code-analysis", "review"},
})
```

### Using a Remote Client

```go
import (
    "github.com/kishoreHQ/aesp/pkg/agent"
)

// Connect to remote kernel
client, err := agent.NewClient("aespd.example.com:50051")
if err != nil {
    log.Fatal(err)
}

a, err := client.GetAgent(ctx, "agent-123")
```

### Publishing Events

```go
import (
    "github.com/kishoreHQ/aesp/pkg/observability"
)

publisher, err := observability.NewEventPublisher(config)
if err != nil {
    log.Fatal(err)
}

event := observability.Event{
    Type: "agent.task.completed",
    Data: map[string]interface{}{
        "agent_id": "agent-123",
        "task_id":  "task-456",
        "duration": "1.23s",
    },
}

if err := publisher.Publish(ctx, event); err != nil {
    log.Printf("failed to publish event: %v", err)
}
```

## Contributing

When adding code to `pkg/`:

1. **Define interfaces first** — Think about what consumers need
2. **Document exported APIs** — Every exported symbol needs GoDoc
3. **Add examples** — Include example functions in test files
4. **Maintain compatibility** — Follow semver for API changes
5. **Generate code** — Run `make generate` after proto changes
6. **Add tests** — Every package needs >80% coverage

## See Also

- [`src/README.md`](../src/README.md) — Internal application source code
- [`cmd/README.md`](../cmd/README.md) — Executable entry points
- [`docs/architecture.md`](../docs/architecture.md) — System architecture
