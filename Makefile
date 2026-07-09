# AESP Reference Implementation — Makefile
# =========================================
# Build automation for the Agent Operating System
#
# Usage:
#   make build          Build all binaries
#   make test           Run unit tests
#   make dev-up         Start local development stack
#   make help           Show all available targets

# ─── Variables ──────────────────────────────────────────────────────────

# Project metadata
PROJECT      := aesp
MODULE       := github.com/kishoreHQ/aesp
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE  ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS      := -X $(MODULE)/pkg/api.Version=$(VERSION) \
                -X $(MODULE)/pkg/api.GitCommit=$(GIT_COMMIT) \
                -X $(MODULE)/pkg/api.BuildDate=$(BUILD_DATE)

# Build configuration
GO           := go
GOFLAGS     ?=
GOTESTFLAGS  := -v -race -count=1
BUILD_DIR    := ./bin
CMD_DIR      := ./cmd
COVERAGE_DIR := ./coverage

# Binary names
DAEMON_BIN   := aespd
CLI_BIN      := aesp-cli
AGENT_BIN    := aesp-agent

# Docker configuration
DOCKER_REGISTRY  ?= ghcr.io/kishorehq
DOCKER_IMAGE     ?= aesp
DOCKER_TAG       ?= $(VERSION)
COMPOSE_FILE     ?= deployments/docker/docker-compose.dev.yml

# Proto configuration
PROTO_DIR        := api/proto
PROTO_GEN_DIR    := pkg/api

# Linting and formatting
GOLANGCI_VERSION := v1.60
GOLANGCI_LINT    := golangci-lint

# ─── Colors ─────────────────────────────────────────────────────────────

BLUE  := \033[34m
GREEN := \033[32m
RED   := \033[31m
RESET := \033[0m

# ─── Default Target ─────────────────────────────────────────────────────

.DEFAULT_GOAL := help

# ─── Help ───────────────────────────────────────────────────────────────

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)AESP Reference Implementation — Build System$(RESET)"
	@echo "Version: $(VERSION) | Commit: $(GIT_COMMIT)"
	@echo ""
	@echo "$(GREEN)Usage:$(RESET)"
	@echo "  make $(GREEN)<target>$(RESET)"
	@echo ""
	@echo "$(GREEN)Build Targets:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'

# ─── Build Targets ──────────────────────────────────────────────────────

.PHONY: build
build: build-aespd build-cli build-agent ## Build all binaries
	@echo "$(GREEN)All binaries built successfully in $(BUILD_DIR)/$(RESET)"

.PHONY: build-aespd
build-aespd: ## Build the AESP daemon
	@echo "$(BLUE)Building $(DAEMON_BIN)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(DAEMON_BIN) $(CMD_DIR)/$(DAEMON_BIN)
	@echo "$(GREEN)✓ $(BUILD_DIR)/$(DAEMON_BIN)$(RESET)"

.PHONY: build-cli
build-cli: ## Build the CLI tool
	@echo "$(BLUE)Building $(CLI_BIN)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(CLI_BIN) $(CMD_DIR)/$(CLI_BIN)
	@echo "$(GREEN)✓ $(BUILD_DIR)/$(CLI_BIN)$(RESET)"

.PHONY: build-agent
build-agent: ## Build the standalone agent runner
	@echo "$(BLUE)Building $(AGENT_BIN)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(AGENT_BIN) $(CMD_DIR)/$(AGENT_BIN)
	@echo "$(GREEN)✓ $(BUILD_DIR)/$(AGENT_BIN)$(RESET)"

.PHONY: build-all
build-all: ## Build for all supported platforms
	@echo "$(BLUE)Cross-compiling for all platforms...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(DAEMON_BIN)-linux-amd64 $(CMD_DIR)/$(DAEMON_BIN)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(CLI_BIN)-linux-amd64 $(CMD_DIR)/$(CLI_BIN)
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(DAEMON_BIN)-linux-arm64 $(CMD_DIR)/$(DAEMON_BIN)
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(DAEMON_BIN)-darwin-amd64 $(CMD_DIR)/$(DAEMON_BIN)
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(DAEMON_BIN)-darwin-arm64 $(CMD_DIR)/$(DAEMON_BIN)
	@echo "$(GREEN)✓ Cross-compilation complete$(RESET)"

.PHONY: install
install: build ## Install binaries to $GOPATH/bin
	@echo "$(BLUE)Installing binaries...$(RESET)"
	$(GO) install -ldflags "$(LDFLAGS)" $(CMD_DIR)/$(DAEMON_BIN)
	$(GO) install -ldflags "$(LDFLAGS)" $(CMD_DIR)/$(CLI_BIN)
	$(GO) install -ldflags "$(LDFLAGS)" $(CMD_DIR)/$(AGENT_BIN)
	@echo "$(GREEN)✓ Installed to $(shell go env GOPATH)/bin$(RESET)"

.PHONY: clean
clean: ## Remove build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -rf $(BUILD_DIR) $(COVERAGE_DIR)
	@echo "$(GREEN)✓ Cleaned$(RESET)"

# ─── Development Targets ────────────────────────────────────────────────

.PHONY: dev-up
dev-up: ## Start local development stack (Docker Compose)
	@echo "$(BLUE)Starting development stack...$(RESET)"
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "$(GREEN)✓ Development stack is running$(RESET)"
	@echo "  API:      http://localhost:8080"
	@echo "  gRPC:     localhost:50051"
	@echo "  Postgres: localhost:5432"
	@echo "  Redis:    localhost:6379"
	@echo "  Grafana:  http://localhost:3000"

.PHONY: dev-down
dev-down: ## Stop local development stack
	@echo "$(BLUE)Stopping development stack...$(RESET)"
	docker compose -f $(COMPOSE_FILE) down
	@echo "$(GREEN)✓ Development stack stopped$(RESET)"

.PHONY: dev-logs
dev-logs: ## Show development stack logs
	docker compose -f $(COMPOSE_FILE) logs -f

.PHONY: dev-reset
dev-reset: dev-down ## Reset development stack (removes volumes)
	docker compose -f $(COMPOSE_FILE) down -v
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "$(GREEN)✓ Development stack reset$(RESET)"

# ─── Test Targets ───────────────────────────────────────────────────────

.PHONY: test
test: ## Run unit tests
	@echo "$(BLUE)Running unit tests...$(RESET)"
	$(GO) test $(GOTESTFLAGS) ./pkg/... ./src/...
	@echo "$(GREEN)✓ Unit tests passed$(RESET)"

.PHONY: test-short
test-short: ## Run unit tests (short mode)
	@echo "$(BLUE)Running short tests...$(RESET)"
	$(GO) test -short ./pkg/... ./src/...

.PHONY: test-integration
test-integration: dev-up ## Run integration tests (requires Docker)
	@echo "$(BLUE)Running integration tests...$(RESET)"
	$(GO) test -v -tags=integration ./test/integration/...
	@echo "$(GREEN)✓ Integration tests passed$(RESET)"

.PHONY: test-e2e
test-e2e: dev-up ## Run end-to-end tests (requires Docker)
	@echo "$(BLUE)Running E2E tests...$(RESET)"
	$(GO) test -v -tags=e2e ./test/e2e/...
	@echo "$(GREEN)✓ E2E tests passed$(RESET)"

.PHONY: test-pkg
test-pkg: ## Run tests for a specific package (use PKG=./pkg/agent)
	@echo "$(BLUE)Running tests for $(PKG)...$(RESET)"
	$(GO) test $(GOTESTFLAGS) $(PKG)

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	$(GO) test -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./pkg/... ./src/...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(GREEN)✓ Coverage report: $(COVERAGE_DIR)/coverage.html$(RESET)"

.PHONY: test-coverage-summary
test-coverage-summary: ## Show coverage summary
	@mkdir -p $(COVERAGE_DIR)
	$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./pkg/... ./src/... > /dev/null
	$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1

# ─── Code Quality ───────────────────────────────────────────────────────

.PHONY: lint
lint: ## Run linters
	@echo "$(BLUE)Running linters...$(RESET)"
	@if ! command -v $(GOLANGCI_LINT) &> /dev/null; then \
		echo "$(RED)golangci-lint not found. Installing...$(RESET)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
			sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION); \
	fi
	$(GOLANGCI_LINT) run ./...
	@echo "$(GREEN)✓ Linting passed$(RESET)"

.PHONY: fmt
fmt: ## Format all Go files
	@echo "$(BLUE)Formatting code...$(RESET)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ Formatted$(RESET)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	$(GO) vet ./pkg/... ./src/... ./cmd/...
	@echo "$(GREEN)✓ Vet passed$(RESET)"

.PHONY: generate
generate: ## Run code generation (protobuf, mocks, etc.)
	@echo "$(BLUE)Running code generation...$(RESET)"
	$(GO) generate ./...
	@echo "$(GREEN)✓ Code generation complete$(RESET)"

.PHONY: proto
proto: ## Generate Protocol Buffer code
	@echo "$(BLUE)Generating protobuf code...$(RESET)"
	@mkdir -p $(PROTO_GEN_DIR)
	@for file in $(PROTO_DIR)/*.proto; do \
		protoc --go_out=$(PROTO_GEN_DIR) \
			--go-grpc_out=$(PROTO_GEN_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_opt=paths=source_relative \
			-I $(PROTO_DIR) \
			"$$file"; \
	done
	@echo "$(GREEN)✓ Protobuf code generated$(RESET)"

.PHONY: mocks
mocks: ## Generate test mocks
	@echo "$(BLUE)Generating mocks...$(RESET)"
	@if ! command -v mockery &> /dev/null; then \
		echo "$(RED)mockery not found. Install with: go install github.com/vektra/mockery/v2@latest$(RESET)"; \
		exit 1; \
	fi
	mockery --all --output pkg/testutil/mocks
	@echo "$(GREEN)✓ Mocks generated$(RESET)"

# ─── Dependency Management ──────────────────────────────────────────────

.PHONY: deps
deps: ## Install development dependencies
	@echo "$(BLUE)Installing dependencies...$(RESET)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies installed$(RESET)"

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "$(BLUE)Updating dependencies...$(RESET)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(RESET)"

.PHONY: verify
verify: ## Verify dependencies
	@echo "$(BLUE)Verifying dependencies...$(RESET)"
	$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies verified$(RESET)"

# ─── Docker Targets ─────────────────────────────────────────────────────

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(RESET)"
	docker build \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest \
		-f deployments/docker/Dockerfile \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		.
	@echo "$(GREEN)✓ Docker image built$(RESET)"

.PHONY: docker-push
docker-push: docker-build ## Push Docker image to registry
	@echo "$(BLUE)Pushing Docker image...$(RESET)"
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest
	@echo "$(GREEN)✓ Docker image pushed$(RESET)"

.PHONY: docker-run
docker-run: ## Run Docker container locally
	@echo "$(BLUE)Running Docker container...$(RESET)"
	docker run -it --rm \
		-p 8080:8080 \
		-p 50051:50051 \
		-v $(PWD)/config:/etc/aesp \
		$(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) \
		aespd serve --config /etc/aesp/config.yaml

# ─── Release Targets ────────────────────────────────────────────────────

.PHONY: release-snapshot
release-snapshot: ## Build release snapshot (no publishing)
	@echo "$(BLUE)Building release snapshot...$(RESET)"
	@if ! command -v goreleaser &> /dev/null; then \
		echo "$(RED)goreleaser not found. Install from https://goreleaser.com/install/$(RESET)"; \
		exit 1; \
	fi
	goreleaser release --snapshot --clean
	@echo "$(GREEN)✓ Snapshot built$(RESET)"

.PHONY: release
release: ## Create a release (requires GITHUB_TOKEN)
	@echo "$(BLUE)Creating release...$(RESET)"
	goreleaser release --clean
	@echo "$(GREEN)✓ Release created$(RESET)"

# ─── Utility Targets ────────────────────────────────────────────────────

.PHONY: version
version: ## Show current version
	@echo "$(BLUE)AESP Reference Implementation$(RESET)"
	@echo "  Version:    $(VERSION)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Build Date: $(BUILD_DATE)"
	@echo "  Go Version: $(shell $(GO) version)"

.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "$(GREEN)✓ All checks passed$(RESET)"

.PHONY: ci
ci: deps generate check test-coverage ## CI pipeline (deps, generate, checks, coverage)
	@echo "$(GREEN)✓ CI pipeline complete$(RESET)"

# ─── Pre-commit Hook ────────────────────────────────────────────────────

.PHONY: pre-commit
pre-commit: fmt vet test-short lint ## Run pre-commit checks
	@echo "$(GREEN)✓ Pre-commit checks passed$(RESET)"

# ─── Security ───────────────────────────────────────────────────────────

.PHONY: security-scan
security-scan: ## Run security scanner (govulncheck)
	@echo "$(BLUE)Running security scan...$(RESET)"
	@if ! command -v govulncheck &> /dev/null; then \
		echo "$(RED)govulncheck not found. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest$(RESET)"; \
		exit 1; \
	fi
	govulncheck ./...
	@echo "$(GREEN)✓ Security scan complete$(RESET)"

# ─── Benchmarks ─────────────────────────────────────────────────────────

.PHONY: bench
bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	$(GO) test -bench=. -benchmem ./pkg/... ./src/... 2>/dev/null | \
		grep -E "^(Benchmark|PASS|FAIL|ok)"
	@echo "$(GREEN)✓ Benchmarks complete$(RESET)"
