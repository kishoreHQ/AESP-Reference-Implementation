# Configuration System

> **Purpose**: This directory contains configuration schemas, defaults, validation rules, and example configurations for the AESP Agent Operating System.

---

## Overview

The AESP Reference Implementation uses a **layered configuration system** that combines multiple sources with well-defined precedence. This allows flexible deployment across environments while maintaining consistency and type safety.

### Configuration Philosophy

1. **Convention over Configuration**: Sensible defaults minimize required configuration
2. **Environment-Specific**: Same binary, different configs per environment
3. **Version Controlled**: Configuration schemas are versioned alongside code
4. **Validatable**: All configuration is validated at startup with clear error messages
5. **Secure**: Secrets are never stored in config files; use environment variables or secret management

## Configuration Sources (Precedence Order)

Configuration is loaded from multiple sources, with higher precedence sources overriding lower ones:

```
1. Command-line flags         (highest priority)
2. Environment variables
3. Configuration file
4. Built-in defaults          (lowest priority)
```

### Example Resolution

```yaml
# config.yaml
model:
  provider: "openai"
  model: "gpt-4o"
```

```bash
# Environment variable overrides config file
export AESP_MODEL_PROVIDER="anthropic"
export AESP_MODEL_MODEL="claude-3-5-sonnet"

# Command-line flag overrides everything
aespd serve --model.provider="openai" --model.model="gpt-4o-mini"
```

## Directory Structure

```
config/
├── README.md               # This file
├── schema.json             # JSON Schema for configuration validation
├── defaults.yaml           # Built-in default values
├── validation.go           # Configuration validation rules (source)
│
├── examples/               # Example configurations
│   ├── minimal.yaml       # Minimal working configuration
│   ├── development.yaml   # Development environment
│   ├── production.yaml    # Production environment
│   ├── kubernetes.yaml    # Kubernetes-specific settings
│   └── multi-provider.yaml # Multiple LLM provider setup
│
└── environments/           # Environment-specific configs (for the project itself)
    ├── dev.yaml
    ├── staging.yaml
    └── production.yaml
```

## Configuration Schema

### Top-Level Structure

```yaml
# Configuration file schema (v1)

# ─── Server Configuration ──────────────────────────────────────────────
server:
  # HTTP REST API
  http:
    address: ":8080"           # Listen address
    read_timeout: "30s"        # Request read timeout
    write_timeout: "30s"       # Response write timeout
    idle_timeout: "120s"       # Keep-alive timeout
    max_header_size: "1MB"     # Maximum header size
    tls:
      enabled: false
      cert_file: ""            # Path to TLS certificate
      key_file: ""             # Path to TLS key
      
  # gRPC API
  grpc:
    address: ":50051"
    max_recv_size: "16MB"
    max_send_size: "16MB"
    tls:
      enabled: false
      cert_file: ""
      key_file: ""
      
  # WebSocket for real-time events
  websocket:
    enabled: true
    address: ":8081"
    read_buffer_size: "4KB"
    write_buffer_size: "4KB"

# ─── Model / LLM Configuration ─────────────────────────────────────────
model:
  # Default provider
  provider: "openai"
  api_key: "${OPENAI_API_KEY}"  # Supports env var substitution
  base_url: ""                  # Custom endpoint (for proxies)
  
  # Default model settings
  model: "gpt-4o"
  temperature: 0.7
  max_tokens: 4096
  timeout: "60s"
  
  # Routing configuration
  routing:
    strategy: "capability"      # capability | cost | latency | fallback
    providers:
      - name: "openai"
        priority: 1
        models:
          - "gpt-4o"
          - "gpt-4o-mini"
        weight: 70
        
      - name: "anthropic"
        priority: 2
        models:
          - "claude-3-5-sonnet"
          - "claude-3-haiku"
        weight: 30
        
  # Fallback configuration
  fallback:
    enabled: true
    max_retries: 3
    retry_delay: "2s"
    fallback_providers:
      - "openai"
      - "anthropic"
      - "ollama"                 # Local fallback

# ─── Memory / Persistence Configuration ────────────────────────────────
memory:
  # Short-term memory (session cache)
  short_term:
    type: "redis"               # redis | memory (in-process)
    redis:
      address: "localhost:6379"
      password: "${REDIS_PASSWORD}"
      db: 0
      pool_size: 10
      
  # Long-term memory (persistent storage)
  long_term:
    type: "postgresql"
    postgresql:
      host: "localhost"
      port: 5432
      database: "aesp"
      username: "aesp"
      password: "${DATABASE_PASSWORD}"
      ssl_mode: "disable"       # disable | require | verify-ca | verify-full
      max_connections: 20
      max_idle: 5
      
  # Vector storage for semantic search
  vector:
    type: "pgvector"            # pgvector | weaviate | pinecone
    pgvector:
      # Uses the same PostgreSQL connection as long_term
      dimension: 1536           # Embedding dimension
      
  # Object storage for artifacts
  object_store:
    type: "minio"               # minio | s3 | gcs | azure | local
    minio:
      endpoint: "localhost:9000"
      access_key: "${MINIO_ACCESS_KEY}"
      secret_key: "${MINIO_SECRET_KEY}"
      bucket: "aesp-artifacts"
      secure: false

# ─── Swarm Configuration ───────────────────────────────────────────────
swarm:
  max_agents: 100               # Maximum agents per swarm
  default_timeout: "5m"         # Default task timeout
  message_buffer: 1000          # Message buffer size per agent
  heartbeat_interval: "10s"     # Agent health check interval
  failure_threshold: 3          # Failed heartbeats before eviction
  
  # Communication
  communication:
    type: "nats"                # nats | redis | in-process
    nats:
      url: "nats://localhost:4222"
      
# ─── Workflow Configuration ────────────────────────────────────────────
workflow:
  max_concurrent: 50            # Max concurrent workflows
  default_timeout: "30m"        # Default workflow timeout
  checkpoint_interval: "30s"    # State checkpoint frequency
  retry:
    max_attempts: 3
    backoff: "exponential"      # exponential | linear | fixed
    initial_delay: "1s"
    max_delay: "5m"

# ─── Plugin Configuration ──────────────────────────────────────────────
plugin:
  registry:
    type: "filesystem"          # filesystem | http | oci
    filesystem:
      path: "/var/lib/aesp/plugins"
      
  sandbox:
    enabled: true
    type: "gvisor"              # gvisor | wasm | none
    
  allowed_origins:              # Allowed plugin sources
    - "github.com/kishoreHQ/aesp-plugins/*"
    - "registry.aesp.dev/official/*"

# ─── MCP (Model Context Protocol) Configuration ─────────────────────────
mcp:
  servers:
    - name: "filesystem"
      command: "npx"
      args: ["-y", "@modelcontextprotocol/server-filesystem", "/workspace"]
      
    - name: "github"
      env:
        GITHUB_TOKEN: "${GITHUB_TOKEN}"
      command: "npx"
      args: ["-y", "@modelcontextprotocol/server-github"]
      
  # Expose AESP as MCP server
  expose_as_server: true
  server_name: "aesp"

# ─── Observability Configuration ───────────────────────────────────────
observability:
  # Structured logging
  logging:
    level: "info"               # debug | info | warn | error
    format: "json"              # json | text
    output: "stdout"            # stdout | stderr | file
    file:
      path: "/var/log/aesp.log"
      max_size: "100MB"
      max_age: "30d"
      max_backups: 5
      
  # Distributed tracing
  tracing:
    enabled: true
    exporter: "otlp"            # otlp | jaeger | zipkin
    otlp:
      endpoint: "http://localhost:4317"
      insecure: true
    sampling_rate: 1.0          # 0.0 - 1.0
    
  # Metrics
  metrics:
    enabled: true
    exporter: "prometheus"      # prometheus | otlp
    prometheus:
      path: "/metrics"
      port: 9090
      
  # Events
  events:
    enabled: true
    buffer_size: 10000
    exporters:
      - type: "websocket"       # Real-time web UI
      - type: "file"            # Audit log
        path: "/var/log/aesp-events.log"

# ─── Authentication & Authorization ────────────────────────────────────
auth:
  enabled: false                # Enable for production
  
  authentication:
    type: "jwt"                 # jwt | api_key | oidc
    jwt:
      secret: "${JWT_SECRET}"
      issuer: "aesp"
      expiry: "24h"
      
  authorization:
    type: "rbac"                # rbac | abac
    rbac:
      policies:
        - role: "admin"
          permissions: ["*"]
        - role: "developer"
          permissions: ["agent:*", "swarm:*", "workflow:*"]
        - role: "viewer"
          permissions: ["agent:read", "swarm:read", "workflow:read"]

# ─── Feature Flags ─────────────────────────────────────────────────────
features:
  vector_search: true
  multi_tenant: false
  plugin_sandboxing: true
  auto_scaling: false
  advanced_routing: true
  
# ─── Development Settings ──────────────────────────────────────────────
development:
  debug: false
  pprof:
    enabled: false
    address: "localhost:6060"
  grpc_reflection: true         # Enable for grpcurl testing
```

## Environment Variable Substitution

Configuration values support environment variable substitution:

```yaml
# Direct substitution
api_key: "${OPENAI_API_KEY}"

# With default value
api_key: "${OPENAI_API_KEY:default-key}"

# Required (fails if not set)
api_key: "${OPENAI_API_KEY:!required}"
```

All environment variables are prefixed with `AESP_` when using the automatic mapping:

```bash
# Maps to model.provider
export AESP_MODEL_PROVIDER=openai

# Maps to memory.long_term.postgresql.host
export AESP_MEMORY_LONG_TERM_POSTGRESQL_HOST=db.example.com

# Maps to server.http.address
export AESP_SERVER_HTTP_ADDRESS=:8080
```

## Validation

Configuration is validated on startup with clear, actionable error messages:

```go
// Example validation errors:
// - "model.provider: required field is empty"
// - "server.http.address: invalid address format: 'not-an-address'"
// - "memory.postgresql.port: must be between 1 and 65535, got 70000"
// - "auth.jwt.secret: must be at least 32 characters for security"
```

## Environment-Specific Examples

### Minimal Configuration

```yaml
# config/examples/minimal.yaml
# The bare minimum to get started

model:
  provider: "openai"
  api_key: "${OPENAI_API_KEY}"

memory:
  short_term:
    type: "memory"  # Use in-process cache
  long_term:
    type: "sqlite"
    sqlite:
      path: "./aesp.db"
```

### Development Configuration

See `config/examples/development.yaml` for a full development setup with hot-reload, debug endpoints, and local dependencies.

### Production Configuration

See `config/examples/production.yaml` for a hardened production configuration with TLS, authentication, external dependencies, and comprehensive observability.

### Kubernetes Configuration

See `config/examples/kubernetes.yaml` for Kubernetes-specific settings including service discovery, ConfigMap/Secret integration, and health probes.

## Configuration Reload

Certain configuration changes can be applied at runtime without restarting:

```bash
# Send SIGHUP to reload configuration
kill -HUP $(pgrep aespd)

# Or use the API
aesp-cli config reload
```

### Reloadable Settings

- Logging level and format
- Model routing weights
- Feature flags
- Retry policies
- Rate limits

### Non-Reloadable Settings

- Server addresses and ports
- Database connections
- Authentication secrets
- TLS certificates

## Best Practices

1. **Never commit secrets**: Use environment variables for API keys, passwords, and tokens
2. **Use environment-specific configs**: Don't use production config in development
3. **Validate in CI**: Run `aespd config validate --file config.yaml` in CI
4. **Version your configs**: Track configuration changes alongside code
5. **Document deviations**: Comment any non-default settings and their rationale

## Contributing

When adding new configuration options:

1. Add the field to the configuration struct in `src/daemon/config.go`
2. Add validation rules
3. Add default values to `config/defaults.yaml`
4. Update `config/schema.json`
5. Update this README with documentation
6. Add an example to `config/examples/`

## See Also

- [`docs/architecture.md`](../docs/architecture.md) — System architecture
- [`deployments/README.md`](../deployments/README.md) — Deployment configurations
