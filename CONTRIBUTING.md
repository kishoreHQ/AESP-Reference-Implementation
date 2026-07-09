# Contributing to AESP Reference Implementation

Thank you for your interest in contributing to the AESP Reference Implementation! This project is a community-driven effort to build a production-quality Agent Operating System, and we welcome contributions from developers, researchers, and organizations.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Commit Message Conventions](#commit-message-conventions)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)
- [Release Process](#release-process)
- [Community](#community)

---

## Code of Conduct

This project adheres to a Code of Conduct that we expect all contributors to follow:

- **Be respectful**: Treat everyone with respect. Healthy debate is encouraged, but harassment or discrimination is not tolerated.
- **Be constructive**: Provide constructive feedback and be open to receiving it.
- **Be collaborative**: Work together towards the best possible solution.
- **Be professional**: Maintain professionalism in all interactions.

---

## Getting Started

### Prerequisites

Before contributing, ensure you have:

1. A GitHub account
2. Go 1.23+ installed
3. Docker and Docker Compose installed
4. Read the [Architecture Documentation](docs/architecture.md)
5. Familiarized yourself with the [AESP Specification](https://github.com/kishoreHQ/AESP)

### Setting Up Your Development Environment

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/AESP-Reference-Implementation.git
cd AESP-Reference-Implementation

# 3. Add upstream remote
git remote add upstream https://github.com/kishoreHQ/AESP-Reference-Implementation.git

# 4. Install dependencies
make deps

# 5. Verify your setup
make build test
```

---

## Development Workflow

We follow a **fork-and-branch** workflow with trunk-based development.

### Branching Strategy

```
main          -- production-ready code (always deployable)
  │
  ├── feat/agent-kernel-config    -- feature branches
  ├── fix/memory-leak
  ├── docs/architecture-update
  └── refactor/workflow-engine
```

### Creating a Feature Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git rebase upstream/main

# Create a feature branch
git checkout -b feat/descriptive-branch-name
```

### Branch Naming Conventions

| Prefix | Purpose | Example |
|--------|---------|---------|
| `feat/` | New feature | `feat/agent-kernel-config` |
| `fix/` | Bug fix | `fix/memory-leak` |
| `docs/` | Documentation | `docs/api-reference` |
| `refactor/` | Code refactoring | `refactor/workflow-engine` |
| `test/` | Test additions/improvements | `test/swarm-orchestration` |
| `chore/` | Maintenance tasks | `chore/update-dependencies` |
| `perf/` | Performance improvements | `perf/query-optimization` |

---

## Commit Message Conventions

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `style` | Code style changes (formatting, semicolons, etc.) |
| `refactor` | Code refactoring |
| `perf` | Performance improvements |
| `test` | Adding or updating tests |
| `chore` | Build process, dependencies, tooling |
| `ci` | CI/CD changes |

### Scopes

Common scopes in this project:

- `kernel` — Agent Kernel
- `swarm` — Swarm Manager
- `memory` — Memory Service
- `workflow` — Workflow Engine
- `plugin` — Plugin Manager
- `mcp` — MCP Gateway
- `router` — Model Router
- `obs` — Observability
- `api` — API layer
- `cli` — CLI tool
- `docs` — Documentation
- `deps` — Dependencies

### Examples

```
feat(kernel): add agent lifecycle hooks

Implement pre-start, post-start, pre-stop, and post-stop hooks
for the agent kernel to allow plugins to intercept lifecycle events.

Closes #123
```

```
fix(memory): resolve race condition in state persistence

The state persistence layer had a race condition when multiple
goroutines attempted to update agent state simultaneously.
Added mutex-based synchronization to prevent concurrent writes.

Fixes #456
```

---

## Pull Request Process

1. **Before Starting**: Check existing issues and PRs to avoid duplication. Comment on an issue to claim it.

2. **Create Branch**: Create a feature branch from `main` following our naming conventions.

3. **Make Changes**: Write code following our coding standards and include tests.

4. **Run Checks**: Ensure all checks pass before submitting:
   ```bash
   make lint
   make test
   make test-integration
   ```

5. **Update Documentation**: Update relevant documentation if your changes affect behavior.

6. **Submit PR**: Create a pull request with a clear description following our PR template.

7. **Review**: Address reviewer feedback promptly and professionally.

8. **Merge**: Once approved and CI passes, a maintainer will merge your PR.

### PR Checklist

- [ ] Code follows project coding standards
- [ ] All tests pass (`make test`)
- [ ] Integration tests pass (`make test-integration`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation is updated
- [ ] Commit messages follow conventional commits
- [ ] PR description is complete
- [ ] Linked to related issue(s)

---

## Coding Standards

### Go Code

We follow standard Go conventions with some project-specific additions:

- **Formatting**: Use `gofmt` or `goimports`. Enforced by CI.
- **Linting**: Use `golangci-lint` with our configuration in `.golangci.yml`.
- **Testing**: Aim for >80% code coverage on new code.
- **Documentation**: All exported symbols must have GoDoc comments.
- **Error Handling**: Use wrapped errors with context (`fmt.Errorf("...: %w", err)`).
- **Logging**: Use structured logging via the observability package.
- **Concurrency**: Follow Go concurrency best practices. Document goroutine lifetimes.

### Code Review Criteria

Reviewers evaluate code on:

1. **Correctness**: Does it work as intended? Are edge cases handled?
2. **Test Coverage**: Are there adequate tests? Are boundary conditions tested?
3. **Documentation**: Is the code well-documented? Are complex algorithms explained?
4. **Performance**: Are there obvious performance issues?
5. **Security**: Are there security concerns?
6. **Maintainability**: Is the code readable and maintainable?
7. **Architecture**: Does it fit the overall architecture?

---

## Testing Requirements

### Unit Tests

- Every package must have comprehensive unit tests
- Use table-driven tests where appropriate
- Mock external dependencies
- Target: >80% coverage for new code

### Integration Tests

- Test component interactions
- Use Dockerized dependencies
- Test realistic scenarios

### E2E Tests

- Full system tests
- Run against the complete stack
- Test critical user journeys

### Writing Good Tests

```go
func TestAgentKernel_CreateAgent(t *testing.T) {
    // Table-driven test example
    tests := []struct {
        name    string
        config  agent.Config
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid config",
            config: agent.Config{
                Name: "test-agent",
                Capabilities: []string{"test"},
            },
            wantErr: false,
        },
        {
            name: "missing name",
            config: agent.Config{
                Capabilities: []string{"test"},
            },
            wantErr: true,
            errMsg:  "agent name is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            kernel := newTestKernel(t)
            agent, err := kernel.CreateAgent(context.Background(), tt.config)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
                return
            }
            assert.NoError(t, err)
            assert.NotNil(t, agent)
        })
    }
}
```

---

## Documentation

Good documentation is as important as good code.

### Code Documentation

- All exported functions, types, and constants must have GoDoc comments
- Complex algorithms should have inline comments explaining the approach
- Package documentation in `doc.go` files

### Architecture Documentation

Changes to architecture require updating:

- `docs/architecture.md` — System overview
- `docs/adr/` — New ADRs for significant decisions
- Mermaid diagrams for visual clarity

### User Documentation

- Update README.md if user-facing behavior changes
- Add examples for new features
- Update API documentation

---

## Release Process

We follow [Semantic Versioning](https://semver.org/) (SemVer).

### Version Format

```
v<major>.<minor>.<patch>[-<prerelease>]
```

Examples: `v0.1.0`, `v0.2.0-alpha.1`, `v1.0.0-rc.1`

### Release Steps

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create a release PR
4. After merge, tag the release commit
5. GitHub Actions builds and publishes artifacts
6. Release notes are published on GitHub

---

## Community

- **Discord**: Real-time chat and support
- **GitHub Discussions**: Long-form discussions and Q&A
- **Issue Tracker**: Bug reports and feature requests
- **Mailing List**: Announcements and async discussions

### Getting Help

If you need help:

1. Check the [documentation](https://aesp.dev/docs)
2. Search [GitHub Discussions](https://github.com/kishoreHQ/AESP-Reference-Implementation/discussions)
3. Ask in [Discord](https://discord.gg/aesp)
4. Open a [GitHub Issue](https://github.com/kishoreHQ/AESP-Reference-Implementation/issues)

---

## Attribution

This contributing guide is adapted from best practices in the Kubernetes, Go, and CNCF communities.

---

Thank you for contributing to the future of autonomous engineering!
