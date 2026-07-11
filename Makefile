# AESP Reference Implementation — Makefile
# =========================================
# Host-neutral Agent OS (AESP-aligned)
#
# Quick start:
#   make help
#   make show-config    # how memory, session, providers are wired
#   make test           # unit tests
#   make smoke          # build + tests + demo + all examples
#   make demo-memory    # memory-update example (trust labels)
#   make demo-session   # single-agent mission (session events)
#
# Override workspace:
#   make demo WORKSPACE=/tmp/my-aesp

# ─── Variables ──────────────────────────────────────────────────────────

PROJECT      := aesp
MODULE       := github.com/kishoreHQ/AESP-Reference-Implementation
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

GO           := go
GOFLAGS     ?=
# Default tests: no -race (faster local); use make test-race for race detector
GOTESTFLAGS ?= -count=1
BUILD_DIR    := ./bin
CMD_DIR      := ./cmd/aespd
DAEMON_BIN   := aespd
COVERAGE_DIR := ./coverage
EXAMPLES_DIR := ./examples
CONFIG_DIR   := ./config

# Runtime / demo defaults (P2 local-first)
WORKSPACE   ?= $(CURDIR)/.aesp-workspace
SERVE_ADDR  ?= :8080
EXAMPLE     ?= 01-single-agent
MISSION     ?= $(EXAMPLES_DIR)/$(EXAMPLE)/mission.yaml

BLUE  := \033[34m
GREEN := \033[32m
YELLOW:= \033[33m
RED   := \033[31m
BOLD  := \033[1m
RESET := \033[0m

.DEFAULT_GOAL := help

# ─── Help ───────────────────────────────────────────────────────────────

.PHONY: help
help: ## Show this help
	@echo "$(BOLD)$(BLUE)AESP Agent OS — Makefile$(RESET)"
	@echo "Version: $(VERSION) | Commit: $(GIT_COMMIT)"
	@echo ""
	@echo "$(GREEN)Most useful:$(RESET)"
	@echo "  make $(GREEN)show-config$(RESET)     Print memory / session / provider wiring"
	@echo "  make $(GREEN)test$(RESET)            Run all package tests"
	@echo "  make $(GREEN)smoke$(RESET)           Full local smoke (test + demo + examples)"
	@echo "  make $(GREEN)demo-memory$(RESET)     Run memory-update example"
	@echo "  make $(GREEN)demo-session$(RESET)    Run single-agent (session events)"
	@echo "  make $(GREEN)demo-failover$(RESET)   Provider failover demo"
	@echo "  make $(GREEN)serve$(RESET)           HTTP Host Interface on $(SERVE_ADDR)"
	@echo ""
	@echo "$(GREEN)All targets:$(RESET)"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-22s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Config knobs:$(RESET) WORKSPACE=$(WORKSPACE)  EXAMPLE=$(EXAMPLE)  SERVE_ADDR=$(SERVE_ADDR)"

# ─── Config & inspect ───────────────────────────────────────────────────

.PHONY: show-config
show-config: build ## Print how memory, session, tools, providers are configured
	@echo "$(BLUE)── Agent OS configuration (defaults) ──$(RESET)"
	@mkdir -p $(WORKSPACE)
	@AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) config
	@echo ""
	@echo "$(BLUE)── config/default.yaml ──$(RESET)"
	@sed 's/^/  /' $(CONFIG_DIR)/default.yaml 2>/dev/null || true
	@echo ""
	@echo "$(YELLOW)Docs:$(RESET) docs/architecture/memory-model.md · session-lifecycle.md · context-envelope.md"

.PHONY: show-memory
show-memory: ## Explain memory subsystem (INV-04) and how to test it
	@echo "$(BOLD)Unified Memory (INV-04 / AESP-0004)$(RESET)"
	@echo ""
	@echo "  Backend (default):  in-process memory.Memory  (P2-friendly)"
	@echo "  Scopes:             working | session | semantic"
	@echo "  Trust labels:       system | verified | agent | retrieved | untrusted | poison-suspect"
	@echo "  Rule:               every write MUST set a trust label"
	@echo "  Privilege:          untrusted/retrieved cannot authorize privileged actions"
	@echo ""
	@echo "  Package:  pkg/memory"
	@echo "  Tools:    memory.write · memory.read  (pkg/tools/builtin)"
	@echo "  Example:  examples/05-memory-update/mission.yaml"
	@echo ""
	@echo "$(GREEN)Try:$(RESET)"
	@echo "  make test-memory"
	@echo "  make demo-memory"
	@echo "  make test-pkg PKG=./pkg/memory"

.PHONY: show-session
show-session: ## Explain session / mission lifecycle and how to test it
	@echo "$(BOLD)Session & mission lifecycle (AESP-0001 / 0005 / INV-10)$(RESET)"
	@echo ""
	@echo "  Mission submit → aesp.control.mission.accepted"
	@echo "  Session open   → aesp.session.opened"
	@echo "  Route          → aesp.control.route.selected"
	@echo "  Provider/runtime events → journal + execution tree"
	@echo "  Correlation:   WorkUnit id · session · trace · artifacts · HITL"
	@echo ""
	@echo "  Kernel:     pkg/kernel  (Host Interface)"
	@echo "  Loop:       pkg/agentos (full OS assembly)"
	@echo "  Journal:    pkg/replay + pkg/eventbus"
	@echo "  Example:    examples/01-single-agent/mission.yaml"
	@echo ""
	@echo "$(GREEN)Try:$(RESET)"
	@echo "  make demo-session"
	@echo "  make test-session"
	@echo "  make events EXAMPLE=01-single-agent"

.PHONY: show-providers
show-providers: ## List abstract provider/runtime plugins (no vendor hardcodes)
	@echo "$(BOLD)Providers & runtimes (INV-01 / INV-03)$(RESET)"
	@echo ""
	@echo "  provider.mock-remote  priority=100  caps: coding,tools,vision,reasoning,planning"
	@echo "  provider.mock-local   priority=10   caps: coding,tools,local,reasoning,planning"
	@echo "  runtime.generic-loop  capabilitiesIn: tools,streaming,reasoning,coding"
	@echo ""
	@echo "  Routing: capability match + health; never model-name ifs in kernel"
	@echo "  Failover demo: make demo-failover"
	@echo ""
	@echo "  Manifest: plugins/runtimes/generic-loop/runtime.yaml"

# ─── Build ──────────────────────────────────────────────────────────────

.PHONY: build
build: ## Build aespd into ./bin
	@mkdir -p $(BUILD_DIR) $(WORKSPACE)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(DAEMON_BIN) $(CMD_DIR)
	@echo "$(GREEN)✓ $(BUILD_DIR)/$(DAEMON_BIN)$(RESET)"

.PHONY: install
install: ## Install aespd to $$(go env GOPATH)/bin
	$(GO) install $(CMD_DIR)
	@echo "$(GREEN)✓ installed aespd$(RESET)"

.PHONY: clean
clean: ## Remove bin/, coverage/, local workspace
	rm -rf $(BUILD_DIR) $(COVERAGE_DIR) $(WORKSPACE)
	@echo "$(GREEN)✓ cleaned$(RESET)"

# ─── Tests ──────────────────────────────────────────────────────────────

.PHONY: test
test: ## Run all package tests
	@echo "$(BLUE)Running tests...$(RESET)"
	$(GO) test $(GOTESTFLAGS) ./...
	@echo "$(GREEN)✓ all tests passed$(RESET)"

.PHONY: test-race
test-race: ## Run tests with race detector
	$(GO) test -race -count=1 ./...

.PHONY: test-short
test-short: ## Fast tests only
	$(GO) test -short -count=1 ./...

.PHONY: test-pkg
test-pkg: ## Test one package: make test-pkg PKG=./pkg/memory
	@test -n "$(PKG)" || (echo "$(RED)PKG required, e.g. PKG=./pkg/memory$(RESET)" && exit 2)
	$(GO) test $(GOTESTFLAGS) -v $(PKG)

.PHONY: test-memory
test-memory: ## Unit tests: memory + trust labels + memory example path
	@echo "$(BLUE)Memory package + agentos memory scenarios$(RESET)"
	$(GO) test $(GOTESTFLAGS) -v ./pkg/memory/ ./pkg/policy/
	$(GO) test $(GOTESTFLAGS) -v ./pkg/agentos/ -run 'Memory|AllScenarios/example.05'

.PHONY: test-session
test-session: ## Unit tests: kernel session/mission + agentos single mission
	@echo "$(BLUE)Session / mission / journal$(RESET)"
	$(GO) test $(GOTESTFLAGS) -v ./pkg/kernel/ ./pkg/eventbus/ ./pkg/replay/ ./pkg/host/
	$(GO) test $(GOTESTFLAGS) -v ./pkg/agentos/ -run 'Single|AllScenarios/example.01'

.PHONY: test-routing
test-routing: ## Unit tests: capability routing + failover
	$(GO) test $(GOTESTFLAGS) -v ./pkg/router/ ./pkg/providerregistry/ ./pkg/capability/
	$(GO) test $(GOTESTFLAGS) -v ./pkg/agentos/ -run 'Failover|AllScenarios/example.09'

.PHONY: test-hitl
test-hitl: ## Unit tests: HITL no auto-approve
	$(GO) test $(GOTESTFLAGS) -v ./pkg/approval/
	$(GO) test $(GOTESTFLAGS) -v ./pkg/agentos/ -run 'HITL|AllScenarios/example.08'

.PHONY: test-interop
test-interop: ## MCP + A2A golden fixtures
	$(GO) test $(GOTESTFLAGS) -v ./pkg/mcp/ ./pkg/a2a/

.PHONY: test-coverage
test-coverage: ## Coverage report → coverage/coverage.html
	@mkdir -p $(COVERAGE_DIR)
	$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ $(COVERAGE_DIR)/coverage.html$(RESET)"

# ─── Demos & examples ───────────────────────────────────────────────────

.PHONY: demo
demo: build ## Built-in single mission demo
	@mkdir -p $(WORKSPACE)
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) demo

.PHONY: demo-session
demo-session: build ## Single-agent mission (session + event trace)
	@mkdir -p $(WORKSPACE)
	@echo "$(BLUE)Mission: examples/01-single-agent (session lifecycle)$(RESET)"
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run $(EXAMPLES_DIR)/01-single-agent/mission.yaml

.PHONY: demo-memory
demo-memory: build ## Memory-update example (trust-labeled writes)
	@mkdir -p $(WORKSPACE)
	@echo "$(BLUE)Mission: examples/05-memory-update$(RESET)"
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run $(EXAMPLES_DIR)/05-memory-update/mission.yaml

.PHONY: demo-kg
demo-kg: build ## Knowledge-graph update example
	@mkdir -p $(WORKSPACE)
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run $(EXAMPLES_DIR)/06-kg-update/mission.yaml

.PHONY: demo-hitl
demo-hitl: build ## HITL approval path (no auto-approve on expire)
	@mkdir -p $(WORKSPACE)
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run $(EXAMPLES_DIR)/08-hitl/mission.yaml

.PHONY: demo-failover
demo-failover: build ## Provider failover (≥2 abstract providers)
	@mkdir -p $(WORKSPACE)
	@echo "$(BLUE)Expect: provider.mock-remote fails → provider.mock-local$(RESET)"
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run $(EXAMPLES_DIR)/09-provider-failover/mission.yaml

.PHONY: demo-rollback
demo-rollback: build ## Rollback / retry workflow
	@mkdir -p $(WORKSPACE)
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run $(EXAMPLES_DIR)/10-rollback-retry/mission.yaml

.PHONY: run-example
run-example: build ## Run one example: make run-example EXAMPLE=05-memory-update
	@mkdir -p $(WORKSPACE)
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run-example $(EXAMPLE)

.PHONY: run-mission
run-mission: build ## Run mission YAML: make run-mission MISSION=path/to/mission.yaml
	@mkdir -p $(WORKSPACE)
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run $(MISSION)

.PHONY: examples
examples: build ## Run all 10 portable examples
	@mkdir -p $(WORKSPACE)
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run-all-examples

.PHONY: events
events: build ## Run example and print event types only
	@mkdir -p $(WORKSPACE)
	@echo "$(BLUE)Events for EXAMPLE=$(EXAMPLE)$(RESET)"
	@AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run-example $(EXAMPLE) 2>/dev/null | \
		sed -n '/eventTypes/,/]/p' || AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) run-example $(EXAMPLE)

# ─── Conformance & HTTP ─────────────────────────────────────────────────

.PHONY: conformance
conformance: build ## Print AESP MUST catalog (implemented/stubbed/missing)
	$(BUILD_DIR)/$(DAEMON_BIN) conformance

.PHONY: serve
serve: build ## Start HTTP Host Interface (P1/P3): make serve SERVE_ADDR=:8080
	@mkdir -p $(WORKSPACE)
	@echo "$(GREEN)Host Interface on $(SERVE_ADDR)$(RESET)"
	@echo "  POST /v1/missions   GET /health   GET /v1/missions/{id}/tree"
	AESP_WORKSPACE=$(WORKSPACE) $(BUILD_DIR)/$(DAEMON_BIN) serve $(SERVE_ADDR)

.PHONY: smoke
smoke: test build demo examples conformance ## Full local smoke suite
	@echo ""
	@echo "$(GREEN)$(BOLD)✓ smoke complete$(RESET) — tests, demo, 10 examples, conformance"

.PHONY: check
check: test build ## CI-friendly check (tests + binary)
	@echo "$(GREEN)✓ check passed$(RESET)"

# ─── Docs shortcuts ─────────────────────────────────────────────────────

.PHONY: docs-arch
docs-arch: ## List architecture docs
	@ls -1 docs/architecture/*.md

.PHONY: docs-memory
docs-memory: ## Open memory model doc path
	@echo docs/architecture/memory-model.md
	@echo docs/hardening/memory-trust-rules.md
	@echo docs/architecture/session-lifecycle.md
	@echo docs/architecture/context-envelope.md

# ─── Monorepo (kernel + UI) ─────────────────────────────────────────────

.PHONY: dev dev-ui dev-ui-mocks install-ui build-ui monorepo-help

dev: ## Start kernel + Mission Control UI (./scripts/dev.sh)
	bash scripts/dev.sh

dev-ui: ## UI only (live API proxy; start kernel separately)
	npm --prefix ui run dev -- --host 127.0.0.1 --port 5173

dev-ui-mocks: ## UI only with MSW mocks
	bash scripts/dev-ui-only.sh

install-ui: ## npm install in ui/
	npm --prefix ui install

build-ui: ## Production build of Mission Control UI
	npm --prefix ui run build

monorepo-help: ## Show monorepo layout
	@echo "AESP Agent OS monorepo"
	@echo "  cmd/aespd, pkg/     → kernel"
	@echo "  ui/                 → Mission Control UI"
	@echo "  examples/           → portable missions"
	@echo "  make dev            → kernel + UI"
	@echo "  make demo           → kernel CLI demo"
	@echo "  Spec: https://github.com/kishoreHQ/AESP"
