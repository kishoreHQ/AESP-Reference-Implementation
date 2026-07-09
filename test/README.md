# Testing Strategy

> **Purpose**: This directory contains all integration and end-to-end tests for the AESP Agent Operating System. Unit tests are co-located with source code (`*_test.go` files).

---

## Overview

The AESP Reference Implementation follows a **comprehensive testing strategy** based on the test pyramid. Testing is not an afterthought — it is a core part of the development workflow.

### Test Pyramid

```
        /\\
       /  \\         E2E Tests
      /----\\        (test/e2e/)
     /      \\
    /--------\\
   /            \\   Integration Tests
  /--------------\\  (test/integration/)
 /                \\
/------------------\\ Unit Tests
                     (*_test.go alongside source)
```

| Level | Location | Scope | Speed | Count |
|-------|----------|-------|-------|-------|
| **Unit** | `pkg/*/*_test.go` | Single function/type | <10ms | Hundreds |
| **Integration** | `test/integration/` | Component interactions | <30s | Dozens |
| **E2E** | `test/e2e/` | Full system scenarios | <5min | Tens |

## Philosophy

1. **Test Behavior, Not Implementation**: Tests verify what code does, not how it does it
2. **Fast Feedback**: Unit tests must complete in under 10 seconds for the entire project
3. **Deterministic**: Tests produce the same results every time (no flaky tests)
4. **Isolated**: Tests don't depend on external state or other tests
5. **Readable**: Tests serve as documentation for expected behavior

## Directory Structure

```
test/
├── README.md                 # This file
├── fixtures/                 # Shared test data and fixtures
│   ├── agents/              # Agent configuration fixtures
│   ├── swarms/              # Swarm definition fixtures
│   ├── workflows/           # Workflow definition fixtures
│   └── responses/           # Mock LLM response fixtures
│
├── helpers/                  # Shared test utilities
│   ├── container.go         # Docker container management
│   ├── database.go          # Database setup/teardown
│   ├── kernel.go            # Test kernel factory
│   ├── assertions.go        # Custom test assertions
│   └── golden.go            # Golden file testing utilities
│
├── integration/              # Integration tests
│   ├── agent_test.go        # Agent lifecycle tests
│   ├── swarm_test.go        # Swarm orchestration tests
│   ├── workflow_test.go     # Workflow execution tests
│   ├── memory_test.go       # Memory service tests
│   ├── plugin_test.go       # Plugin system tests
│   ├── mcp_test.go          # MCP gateway tests
│   ├── router_test.go       # Model router tests
│   └── api_test.go          # API endpoint tests
│
├── e2e/                      # End-to-end tests
│   ├── smoke_test.go        # Smoke / health check tests
│   ├── workflow_test.go     # Full workflow scenarios
│   ├── plugin_test.go       # Plugin lifecycle E2E
│   └── cli_test.go          # CLI tool E2E tests
│
├── benchmark/                # Performance benchmarks
│   ├── kernel_bench.go      # Agent kernel benchmarks
│   ├── swarm_bench.go       # Swarm orchestration benchmarks
│   └── router_bench.go      # Model router benchmarks
│
└── load/                     # Load testing
    ├── k6/                   # k6 load test scripts
    └── locust/               # Locust load test scripts
```

## Unit Tests

Unit tests are co-located with source code in `*_test.go` files. They test individual functions and types in isolation.

### Conventions

```go
// pkg/agent/runtime_test.go
package agent

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRuntime_Execute_Success(t *testing.T) {
    // Arrange
    ctx := context.Background()
    rt := NewTestRuntime(t)
    task := Task{
        Type: "echo",
        Input: map[string]interface{}{
            "message": "hello",
        },
    }
    
    // Act
    result, err := rt.Execute(ctx, task)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, "hello", result.Output["message"])
}

func TestRuntime_Execute_Timeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer cancel()
    
    rt := NewTestRuntime(t)
    task := Task{Type: "slow", Input: nil}
    
    _, err := rt.Execute(ctx, task)
    
    assert.ErrorIs(t, err, context.DeadlineExceeded)
}

// Table-driven test
func TestRuntime_ValidateTask(t *testing.T) {
    tests := []struct {
        name    string
        task    Task
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid task",
            task: Task{Type: "echo", Input: map[string]interface{}{"msg": "hi"}},
            wantErr: false,
        },
        {
            name: "missing type",
            task: Task{Input: map[string]interface{}{"msg": "hi"}},
            wantErr: true,
            errMsg:  "task type is required",
        },
        {
            name: "empty type",
            task: Task{Type: "", Input: map[string]interface{}{}},
            wantErr: true,
            errMsg:  "task type is required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            rt := NewTestRuntime(t)
            err := rt.ValidateTask(tt.task)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
                return
            }
            assert.NoError(t, err)
        })
    }
}
```

### Running Unit Tests

```bash
# All unit tests
make test

# Specific package
make test-pkg PKG=./pkg/agent

# With coverage
make test-coverage

# With race detection
make test-race

# Short mode (skip slow tests)
make test-short
```

## Integration Tests

Integration tests verify that components work together correctly. They use Docker containers for external dependencies.

### Setup

Integration tests use Docker containers managed by the test suite:

```go
// test/integration/agent_test.go
//go:build integration

package integration

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/suite"
    "github.com/kishoreHQ/aesp/test/helpers"
)

type AgentIntegrationSuite struct {
    suite.Suite
    ctx      context.Context
    postgres *helpers.PostgresContainer
    redis    *helpers.RedisContainer
    kernel   *helpers.TestKernel
}

func (s *AgentIntegrationSuite) SetupSuite() {
    s.ctx = context.Background()
    
    // Start PostgreSQL
    s.postgres = helpers.MustStartPostgres(s.T())
    
    // Start Redis
    s.redis = helpers.MustStartRedis(s.T())
    
    // Initialize kernel with test containers
    s.kernel = helpers.NewTestKernel(s.T(), helpers.KernelConfig{
        DatabaseURL: s.postgres.URL(),
        RedisAddr:   s.redis.Addr(),
    })
}

func (s *AgentIntegrationSuite) TearDownSuite() {
    s.postgres.Stop(s.ctx)
    s.redis.Stop(s.ctx)
}

func (s *AgentIntegrationSuite) TestAgent_CreateAndExecute() {
    agent, err := s.kernel.CreateAgent(s.ctx, agent.Config{
        Name: "test-agent",
        Capabilities: []string{"test"},
    })
    s.NoError(err)
    s.NotNil(agent)
    
    result, err := s.kernel.ExecuteTask(s.ctx, agent.ID, agent.Task{
        Type: "echo",
        Input: map[string]interface{}{"message": "hello"},
    })
    s.NoError(err)
    s.Equal("hello", result.Output["message"])
}

func TestAgentIntegration(t *testing.T) {
    suite.Run(t, new(AgentIntegrationSuite))
}
```

### Running Integration Tests

```bash
# Start dependencies and run integration tests
make test-integration

# Requires Docker to be running
# Docker containers are automatically started and stopped
```

### Integration Test Tags

Integration tests are guarded by build tags to prevent accidental execution:

```go
//go:build integration
```

This means they only run when explicitly requested with `-tags=integration`.

## End-to-End Tests

E2E tests verify complete user workflows from the API layer through all components.

### Architecture

```go
// test/e2e/workflow_test.go
//go:build e2e

package e2e

import (
    "testing"
    
    "github.com/stretchr/testify/suite"
    "github.com/kishoreHQ/aesp/test/helpers"
)

type WorkflowE2ESuite struct {
    suite.Suite
    env *helpers.TestEnvironment  // Full stack in containers
}

func (s *WorkflowE2ESuite) SetupSuite() {
    s.env = helpers.MustStartEnvironment(s.T(), helpers.EnvConfig{
        Daemon:   true,
        Postgres: true,
        Redis:    true,
        NATS:     true,
    })
}

func (s *WorkflowE2ESuite) TestSoftwareDevelopmentWorkflow() {
    // Create a swarm for software development
    swarm := s.env.MustCreateSwarm("dev-team", "software-dev")
    
    // Submit a task
    result := s.env.MustSubmitTask(swarm.ID, agent.Task{
        Type: "implement-feature",
        Input: map[string]interface{}{
            "requirements": "Add user authentication to the API",
        },
    })
    
    // Verify result
    s.NotNil(result)
    s.NotEmpty(result.Artifacts)
}
```

### Running E2E Tests

```bash
# Run full E2E test suite (takes several minutes)
make test-e2e

# Requires Docker and the full development stack
```

## Benchmarks

Performance-critical code includes benchmarks:

```go
func BenchmarkKernel_ExecuteTask(b *testing.B) {
    ctx := context.Background()
    kernel := NewTestKernel(b)
    task := Task{Type: "noop", Input: nil}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := kernel.ExecuteTask(ctx, task)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

Run benchmarks:
```bash
make bench
```

## Load Testing

Load tests are in `test/load/` and use industry-standard tools:

### k6

```javascript
// test/load/k6/api-load.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '2m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '2m', target: 200 },
        { duration: '5m', target: 200 },
        { duration: '2m', target: 0 },
    ],
};

export default function () {
    const res = http.post('http://localhost:8080/v1/agents', JSON.stringify({
        name: `load-test-agent-${__VU}-${__ITER}`,
        capabilities: ['test'],
    }), {
        headers: { 'Content-Type': 'application/json' },
    });
    
    check(res, {
        'status is 201': (r) => r.status === 201,
        'response time < 500ms': (r) => r.timings.duration < 500,
    });
    
    sleep(1);
}
```

Run load tests:
```bash
cd test/load/k6
k6 run api-load.js
```

## Test Data Management

### Fixtures

Reusable test data lives in `test/fixtures/`:

```yaml
# test/fixtures/agents/code-reviewer.yaml
name: "code-reviewer"
description: "Reviews code changes"
capabilities:
  - "code-analysis"
  - "review"
  - "suggestions"
model:
  provider: "openai"
  model: "gpt-4o"
```

### Golden Files

Expected outputs for snapshot testing:

```go
func TestWorkflowParser_Parse(t *testing.T) {
    input := helpers.LoadFixture(t, "workflows/simple.yaml")
    result := parser.Parse(input)
    
    helpers.AssertGolden(t, "workflows/simple.json", result)
}
```

Update golden files:
```bash
UPDATE_GOLDEN=1 go test ./test/integration/...
```

## Test Containers

The `test/helpers/` package provides Docker container management:

```go
// Start a PostgreSQL container
postgres := helpers.MustStartPostgres(t, helpers.PostgresConfig{
    Database: "test_aesp",
    Username: "test",
    Password: "test",
})

// Start Redis
redis := helpers.MustStartRedis(t, helpers.RedisConfig{
    Password: "",
})

// Containers are automatically cleaned up via t.Cleanup()
```

## CI Integration

Tests run in CI with the following pipeline:

```yaml
# .github/workflows/ci.yml (excerpt)
test:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    
    - name: Unit Tests
      run: make test
    
    - name: Integration Tests
      run: make test-integration
    
    - name: Coverage
      run: make test-coverage
    
    - name: Upload Coverage
      uses: codecov/codecov-action@v3
```

## Best Practices

1. **Use testify**: `assert` for non-fatal checks, `require` for fatal checks
2. **Table-driven tests**: Use test tables for multiple input/output cases
3. **Parallel tests**: Mark independent tests with `t.Parallel()`
4. **Cleanup**: Use `t.Cleanup()` for resource cleanup
5. **Named subtests**: Use `t.Run("descriptive name", ...)` for clarity
6. **No global state**: Tests must not depend on or modify global state
7. **Deterministic**: Tests must produce the same results every time
8. **Fast**: Unit tests should complete in under 10 seconds total

## Debugging Tests

```bash
# Run specific test
make test-pkg PKG=./pkg/agent TEST=TestRuntime_Execute

# With verbose output
make test-pkg PKG=./pkg/agent GOTESTFLAGS="-v"

# With debugger
dlv test ./pkg/agent -- -test.run TestRuntime_Execute
```

## Contributing

When adding tests:

1. Follow the existing naming conventions (`Test<Type>_<Method>_<Scenario>`)
2. Use the test helpers in `test/helpers/`
3. Add fixtures for reusable test data
4. Update golden files if output format changes
5. Ensure tests are deterministic and don't depend on external services

## See Also

- [`CONTRIBUTING.md`](../CONTRIBUTING.md) — Contribution guidelines
- [`Makefile`](../Makefile) — Test targets
- [`docs/architecture.md`](../docs/architecture.md) — System architecture
