# Deployment Documentation

> **Purpose**: This directory contains all deployment configurations, infrastructure definitions, and operational guides for the AESP Agent Operating System.

---

## Overview

The AESP Reference Implementation supports multiple deployment modes to accommodate different environments and scale requirements:

| Deployment Mode | Use Case | Complexity | Scalability |
|---------------|----------|-----------|-------------|
| **Single Binary** | Local development, testing | Low | Single node |
| **Docker Compose** | Small teams, demos | Low-Medium | Single host |
| **Kubernetes** | Production, enterprise | Medium | Horizontal |
| **Helm Chart** | Cloud-native deployments | Medium | Horizontal |
| **Terraform** | Infrastructure provisioning | Medium | Cloud-specific |

## Deployment Modes

### 1. Single Binary (Development)

The simplest deployment — a single Go binary with embedded SQLite:

```bash
# Build
go build -o aespd ./cmd/aespd

# Run with minimal config
./aespd serve --model.provider=openai --model.api-key=$OPENAI_API_KEY
```

**Requirements**: None beyond the binary itself and an LLM API key.

**Pros**: Zero dependencies, fastest startup, easiest debugging  
**Cons**: No persistence, no high availability, limited scalability

### 2. Docker Compose (Small Teams)

A complete stack with all dependencies managed via Docker Compose:

```bash
cd deployments/docker

# Start the full stack
docker compose -f docker-compose.dev.yml up -d

# View logs
docker compose logs -f

# Stop
docker compose down
```

**Services included**:
- AESP daemon (auto-rebuild on code changes)
- PostgreSQL 15
- Redis 7
- NATS
- Grafana + Prometheus + Jaeger (observability stack)

### 3. Kubernetes (Production)

Production deployment using Kubernetes manifests or Helm charts:

```bash
# Using Helm
helm repo add aesp https://charts.aesp.dev
helm install aesp aesp/aesp \
  --namespace aesp \
  --create-namespace \
  --values values-production.yaml

# Using raw manifests
kubectl apply -k deployments/kubernetes/overlays/production
```

## Directory Structure

```
deployments/
├── README.md                     # This file
│
├── docker/                       # Docker configurations
│   ├── Dockerfile                # Main application Dockerfile
│   ├── Dockerfile.dev            # Development Dockerfile (hot reload)
│   ├── docker-compose.dev.yml    # Development stack
│   ├── docker-compose.prod.yml   # Production stack
│   └── docker-compose.test.yml   # Test stack
│
├── kubernetes/                   # Kubernetes manifests
│   ├── base/                     # Base kustomization
│   │   ├── kustomization.yaml
│   │   ├── namespace.yaml
│   │   ├── daemon-deployment.yaml
│   │   ├── daemon-service.yaml
│   │   ├── configmap.yaml
│   │   └── secrets.yaml
│   │
│   └── overlays/                 # Environment overlays
│       ├── development/
│       │   ├── kustomization.yaml
│       │   ├── replica-patch.yaml
│       │   └── resource-patch.yaml
│       ├── staging/
│       │   ├── kustomization.yaml
│       │   ├── replica-patch.yaml
│       │   └── ingress-patch.yaml
│       └── production/
│           ├── kustomization.yaml
│           ├── replica-patch.yaml
│           ├── resource-patch.yaml
│           ├── hpa-patch.yaml
│           └── pdb-patch.yaml
│
├── helm/                         # Helm chart
│   └── aesp/
│       ├── Chart.yaml
│       ├── values.yaml           # Default values
│       ├── values-production.yaml
│       └── templates/
│           ├── _helpers.tpl
│           ├── deployment.yaml
│           ├── service.yaml
│           ├── ingress.yaml
│           ├── configmap.yaml
│           ├── secrets.yaml
│           ├── hpa.yaml
│           ├── pdb.yaml
│           └── serviceaccount.yaml
│
├── terraform/                    # Infrastructure as Code
│   ├── modules/                  # Reusable modules
│   │   ├── aks/                  # Azure AKS
│   │   ├── eks/                  # AWS EKS
│   │   ├── gke/                  # Google GKE
│   │   ├── database/             # Managed databases
│   │   ├── cache/                # Managed cache
│   │   └── networking/           # VPC, subnets, load balancers
│   │
│   └── environments/             # Environment configurations
│       ├── development/
│       ├── staging/
│       └── production/
│
└── scripts/                      # Deployment scripts
    ├── install.sh               # Quick install script
    ├── upgrade.sh               # Upgrade script
    ├── backup.sh                # Backup script
    └── health-check.sh          # Health check script
```

## Docker Configuration

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o aespd ./cmd/aespd

# Runtime stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /build/aespd /app/aespd
COPY --from=builder /build/config/defaults.yaml /etc/aesp/config.yaml

EXPOSE 8080 50051 9090

USER nonroot:nonroot

ENTRYPOINT ["/app/aespd"]
CMD ["serve", "--config", "/etc/aesp/config.yaml"]
```

**Image characteristics**:
- Based on distroless for minimal attack surface
- Non-root user execution
- Multi-stage build for small image size (~30MB)
- Health check endpoint at `/health`

### Docker Compose Stacks

#### Development Stack

```yaml
# deployments/docker/docker-compose.dev.yml
version: "3.8"

services:
  aespd:
    build:
      context: ../..
      dockerfile: deployments/docker/Dockerfile.dev
    ports:
      - "8080:8080"
      - "50051:50051"
    environment:
      - AESP_MODEL_PROVIDER=${AESP_MODEL_PROVIDER}
      - AESP_MODEL_API_KEY=${AESP_MODEL_API_KEY}
      - AESP_MEMORY_LONG_TERM_POSTGRESQL_HOST=postgres
      - AESP_MEMORY_SHORT_TERM_REDIS_ADDRESS=redis:6379
      - AESP_SWARM_COMMUNICATION_NATS_URL=nats://nats:4222
    volumes:
      - ../..:/app
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      nats:
        condition: service_started
    develop:
      watch:
        - action: rebuild
          path: ../..
          target: /app

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: aesp
      POSTGRES_USER: aesp
      POSTGRES_PASSWORD: aesp
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U aesp"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  nats:
    image: nats:2.10-alpine
    ports:
      - "4222:4222"
      - "8222:8222"  # Monitoring

  # Observability Stack
  prometheus:
    image: prom/prometheus:v2.50
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus

  grafana:
    image: grafana/grafana:10.4
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./grafana/datasources:/etc/grafana/provisioning/datasources

  jaeger:
    image: jaegertracing/all-in-one:1.55
    ports:
      - "16686:16686"  # UI
      - "4317:4317"    # OTLP gRPC

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:
```

## Kubernetes Deployment

### Base Manifests

The base kustomization provides the core resources:

```yaml
# deployments/kubernetes/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: aesp

resources:
  - namespace.yaml
  - configmap.yaml
  - secrets.yaml
  - daemon-deployment.yaml
  - daemon-service.yaml

commonLabels:
  app.kubernetes.io/name: aesp
  app.kubernetes.io/managed-by: kustomize
```

### Production Overlay

The production overlay adds scaling, HA, and security configurations:

```yaml
# deployments/kubernetes/overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namePrefix: prod-

resources:
  - ../../base
  - ingress.yaml

patches:
  - path: replica-patch.yaml
  - path: resource-patch.yaml
  - path: hpa-patch.yaml
  - path: pdb-patch.yaml

configMapGenerator:
  - name: aesp-config
    behavior: merge
    files:
      - config/production.yaml
```

```yaml
# deployments/kubernetes/overlays/production/replica-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aesp-daemon
spec:
  replicas: 3
```

```yaml
# deployments/kubernetes/overlays/production/hpa-patch.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: aesp-daemon
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: aesp-daemon
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

## Helm Chart

### Installation

```bash
# Add the Helm repository
helm repo add aesp https://charts.aesp.dev
helm repo update

# Install with default values
helm install aesp aesp/aesp --namespace aesp --create-namespace

# Install with custom values
helm install aesp aesp/aesp \
  --namespace aesp \
  --create-namespace \
  --values my-values.yaml

# Upgrade
helm upgrade aesp aesp/aesp --values my-values.yaml

# Uninstall
helm uninstall aesp --namespace aesp
```

### Key Values

```yaml
# values-production.yaml
replicaCount: 3

image:
  repository: ghcr.io/kishorehq/aesp
  tag: "v0.1.0"
  pullPolicy: IfNotPresent

resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "2000m"

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
  targetCPUUtilizationPercentage: 70

model:
  provider: "openai"
  apiKey:
    existingSecret: "aesp-secrets"
    existingSecretKey: "openai-api-key"

persistence:
  postgresql:
    enabled: true
    host: "postgres.example.com"
    existingSecret: "aesp-db-credentials"
  redis:
    enabled: true
    host: "redis.example.com"

service:
  type: ClusterIP
  ports:
    http: 8080
    grpc: 50051
    metrics: 9090

ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt"
  hosts:
    - host: aesp.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: aesp-tls
      hosts:
        - aesp.example.com

observability:
  prometheus:
    enabled: true
  grafana:
    enabled: true
  jaeger:
    enabled: true
```

## Terraform Infrastructure

### AWS EKS Example

```hcl
# deployments/terraform/environments/production/main.tf
module "eks" {
  source = "../../modules/eks"

  cluster_name   = "aesp-production"
  cluster_version = "1.29"
  
  vpc_cidr = "10.0.0.0/16"
  
  node_groups = {
    general = {
      desired_size = 3
      min_size     = 2
      max_size     = 10
      
      instance_types = ["m6i.xlarge"]
      capacity_type  = "ON_DEMAND"
    }
  }
}

module "database" {
  source = "../../modules/database/aws"
  
  engine         = "postgres"
  engine_version = "15.4"
  instance_class = "db.r6g.xlarge"
  
  database_name = "aesp"
  username      = "aesp"
  
  allocated_storage = 100
}

module "cache" {
  source = "../../modules/cache/aws"
  
  engine         = "redis"
  node_type      = "cache.r6g.large"
  num_cache_nodes = 2
}
```

Deploy:
```bash
cd deployments/terraform/environments/production
terraform init
terraform plan
terraform apply
```

## Health Checks and Monitoring

### Health Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Liveness probe — always returns 200 if process is running |
| `GET /ready` | Readiness probe — returns 200 when all dependencies are healthy |
| `GET /metrics` | Prometheus metrics |

### Readiness Checks

The readiness probe verifies:
- Database connectivity
- Cache connectivity
- Message broker connectivity
- LLM provider availability

### Alerts

Recommended alerts:

| Alert | Severity | Condition |
|-------|----------|-----------|
| HighErrorRate | Critical | >5% error rate for 5 minutes |
| HighLatency | Warning | p99 latency >2s for 10 minutes |
| PodRestarting | Warning | >3 restarts in 15 minutes |
| LLMRateLimit | Warning | >80% rate limit utilization |

## Backup and Disaster Recovery

### Database Backups

```bash
# Automated daily backups (PostgreSQL)
pg_dump -h $DB_HOST -U aesp aesp | gzip > aesp-backup-$(date +%Y%m%d).sql.gz

# Restore
 gunzip < aesp-backup-20250115.sql.gz | psql -h $DB_HOST -U aesp aesp
```

### Configuration Backups

```bash
# Export configuration
aesp-cli config export > aesp-config-backup.yaml

# Import configuration
aesp-cli config import --file aesp-config-backup.yaml
```

## Security Considerations

### Network Security

- All inter-service communication uses mTLS in production
- External API access is behind an API gateway with rate limiting
- LLM API keys are stored in Kubernetes secrets or secret managers

### Secrets Management

| Environment | Recommended Approach |
|------------|---------------------|
| Development | `.env` files (not committed) |
| Staging | Kubernetes secrets |
| Production | External secret manager (Vault, AWS Secrets Manager) |

### Pod Security

- Run as non-root user
- Read-only root filesystem
- No privileged containers
- Security context with dropped capabilities

## Operational Runbooks

### Scaling Up

```bash
# Kubernetes: scale up
kubectl scale deployment aesp-daemon --replicas=10 -n aesp

# Or update HPA
kubectl patch hpa aesp-daemon -n aesp --patch '{"spec":{"maxReplicas":30}}'
```

### Rolling Back

```bash
# Helm rollback
helm rollback aesp 2 -n aesp

# Kubernetes rollback
kubectl rollout undo deployment/aesp-daemon -n aesp
```

### Debugging

```bash
# Check pod status
kubectl get pods -n aesp

# View logs
kubectl logs -f deployment/aesp-daemon -n aesp

# Exec into pod
kubectl exec -it deployment/aesp-daemon -n aesp -- /bin/sh

# Port forward for local debugging
kubectl port-forward deployment/aesp-daemon -n aesp 8080:8080
```

## Contributing

When adding deployment configurations:

1. Follow existing directory structure
2. Include health checks in all deployments
3. Use non-root containers
4. Document resource requirements
5. Include scaling guidelines
6. Test in the appropriate environment

## See Also

- [`config/README.md`](../config/README.md) — Configuration system
- [`docs/architecture.md`](../docs/architecture.md) — System architecture
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
