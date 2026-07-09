# Command-Line Tools (`cmd/`)

> **Purpose**: The `cmd/` directory contains the main executable entry points for the AESP Agent Operating System. Each subdirectory compiles into a separate binary.

---

## Overview

Following Go project conventions, `cmd/` houses all executable programs in this repository. Each subdirectory represents a single command-line tool with its own `main.go` file.

This separation ensures:
- **Clear entry points**: Each binary has exactly one location
- **Independent versioning**: Different tools can be versioned separately
- **Minimal binary size**: Each binary only includes necessary dependencies
- **Clean imports**: No accidental imports of `main` packages

## Directory Structure

```
cmd/
├── aespd/                  # Main daemon (control plane server)
│   ├── main.go            # Entry point
│   ├── serve.go           # Server command
│   ├── migrate.go         # Database migration command
│   ├── config.go          # Configuration management
│   └── version.go         # Version information
│
├── aesp-cli/              # CLI tool for interacting with AESP
│   ├── main.go            # Entry point
│   ├── agent.go           # Agent management commands
│   ├── swarm.go           # Swarm management commands
│   ├── workflow.go        # Workflow management commands
│   ├── plugin.go          # Plugin management commands
│   ├── logs.go            # Log streaming commands
│   ├── config.go          # CLI configuration
│   └── version.go         # Version information
│
└── aesp-agent/            # Standalone agent runner
    ├── main.go            # Entry point
    ├── run.go             # Run a single agent
    ├── serve.go           # Start an agent worker
    └── version.go         # Version information
```

## Tools

### `aespd` — The AESP Daemon

The main server process that runs the AESP Agent Operating System. This is the control plane that manages agents, swarms, workflows, and all supporting services.

**Usage:**
```bash
# Start the daemon with default configuration
aespd serve

# Start with custom config file
aespd serve --config /etc/aesp/config.yaml

# Run database migrations
aespd migrate up

# Check version
aespd version

# Show help
aespd --help
```

**Configuration:**

The daemon loads configuration from (in order of precedence):
1. Command-line flags (highest priority)
2. Environment variables (`AESP_*`)
3. Configuration file (`--config` or default paths)
4. Built-in defaults (lowest priority)

```yaml
# Example configuration
server:
  http:
    address: ":8080"
  grpc:
    address: ":50051"
  websocket:
    address: ":8081"

model:
  provider: "openai"
  api_key: "${OPENAI_API_KEY}"
  default_model: "gpt-4o"

memory:
  type: "postgresql"
  connection_string: "${DATABASE_URL}"

observability:
  tracing:
    enabled: true
    exporter: "otlp"
    endpoint: "http://localhost:4317"
  metrics:
    enabled: true
    port: 9090
```

### `aesp-cli` — Command-Line Interface

The primary user-facing tool for interacting with an AESP deployment. It provides a comprehensive set of commands for managing the entire system.

**Usage:**
```bash
# Agent management
aesp-cli agent create --name "code-reviewer" --capability "code-analysis"
aesp-cli agent list
aesp-cli agent get --id agent-123
aesp-cli agent delete --id agent-123
aesp-cli agent logs --id agent-123 --follow

# Swarm management
aesp-cli swarm create --name "dev-team" --template "software-dev"
aesp-cli swarm list
aesp-cli swarm get --id swarm-456
aesp-cli swarm deploy --id swarm-456 --task "review-pr-789"
aesp-cli swarm delete --id swarm-456

# Workflow management
aesp-cli workflow create --file workflow.yaml
aesp-cli workflow list
aesp-cli workflow run --id workflow-789
aesp-cli workflow status --id workflow-789
aesp-cli workflow logs --id workflow-789 --follow

# Plugin management
aesp-cli plugin list
aesp-cli plugin install --name "github-integration"
aesp-cli plugin uninstall --name "github-integration"
aesp-cli plugin info --name "github-integration"

# Configuration
aesp-cli config set --key "model.provider" --value "anthropic"
aesp-cli config get --key "model.provider"
aesp-cli config view

# Context (connection to daemon)
aesp-cli context set --name "production" --endpoint "https://aesp.prod.internal:50051"
aesp-cli context use --name "production"
aesp-cli context list
```

**Global Flags:**
```
--endpoint string   AESP daemon endpoint (default: http://localhost:8080)
--token string      Authentication token
--namespace string  Namespace for resource isolation (default: default)
--output string     Output format: json, yaml, table (default: table)
--verbose           Enable verbose output
```

### `aesp-agent` — Standalone Agent Runner

A lightweight tool for running a single agent without the full control plane. Useful for development, testing, and edge deployments.

**Usage:**
```bash
# Run an agent with a task file
aesp-agent run --config agent.yaml --task task.yaml

# Start an agent worker that connects to a control plane
aesp-agent serve --endpoint "https://aesp.internal:50051" --token "worker-token"

# Run inline task
aesp-agent run --name "summarizer" --capability "text-summarization" --input "Text to summarize..."
```

**Agent Configuration:**
```yaml
# agent.yaml
name: "code-reviewer"
description: "Reviews code changes for quality"
capabilities:
  - "code-analysis"
  - "review"
  - "suggestions"

model:
  provider: "openai"
  model: "gpt-4o"

memory:
  type: "local"  # Uses local SQLite for standalone mode
  path: "./agent.db"

plugins:
  - name: "git-tools"
    config:
      repo_path: "."
```

## Building

All binaries are built using the Makefile:

```bash
# Build all binaries
make build

# Build specific binary
make build-aespd
make build-cli
make build-agent

# Build for specific platform
GOOS=linux GOARCH=amd64 make build
GOOS=darwin GOARCH=arm64 make build

# Build with version info
VERSION=0.1.0 GIT_COMMIT=abc123 make build
```

## Installation

### From Source

```bash
git clone https://github.com/kishoreHQ/AESP-Reference-Implementation.git
cd AESP-Reference-Implementation
make install
```

This installs binaries to `$GOPATH/bin` or `$HOME/go/bin`.

### Docker

```bash
# Pull the image
docker pull ghcr.io/kishorehq/aesp:latest

# Run the daemon
docker run -p 8080:8080 -p 50051:50051 ghcr.io/kishorehq/aesp:latest aespd serve

# Run CLI via Docker
docker run --rm ghcr.io/kishorehq/aesp:latest aesp-cli version
```

## Shell Completion

The CLI supports shell completion for bash, zsh, and fish:

```bash
# Bash
source <(aesp-cli completion bash)

# Zsh
source <(aesp-cli completion zsh)

# Fish
aesp-cli completion fish | source
```

## Contributing

When adding a new command-line tool:

1. Create a new directory under `cmd/` with the tool name
2. Implement `main.go` with proper flag parsing and subcommands
3. Use the shared logging and configuration packages from `pkg/`
4. Add the binary to the Makefile build targets
5. Update this README with documentation
6. Add integration tests in `test/integration/cli/`

## See Also

- [`src/README.md`](../src/README.md) — Internal source code organization
- [`pkg/README.md`](../pkg/README.md) — Public library packages
- [`Makefile`](../Makefile) — Build automation
