# Pull Request

## Summary

<!-- Provide a brief summary of the changes in this PR. What problem does it solve? -->

## Related Issues

<!-- Link to related issues. Use "Fixes #123" or "Closes #456" to auto-close issues on merge. -->

Fixes #
Closes #

## Type of Change

<!-- Mark the relevant option with an "x" (e.g., [x]) -->

- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring (no functional changes)
- [ ] Test improvements
- [ ] CI/CD changes
- [ ] Dependency update

## Component

<!-- Mark the relevant component(s) with an "x" -->

- [ ] Agent Kernel (`pkg/agent/`, `src/agent/`)
- [ ] Swarm Manager (`pkg/swarm/`, `src/swarm/`)
- [ ] Workflow Engine (`pkg/workflow/`, `src/workflow/`)
- [ ] Memory Service (`pkg/memory/`, `src/memory/`)
- [ ] Plugin Manager (`pkg/plugin/`, `src/plugin/`)
- [ ] MCP Gateway (`pkg/mcp/`, `src/mcp/`)
- [ ] Model Router (`pkg/router/`, `src/router/`)
- [ ] Observability (`pkg/observability/`, `src/observability/`)
- [ ] API / SDK (`pkg/api/`)
- [ ] CLI (`cmd/aesp-cli/`)
- [ ] Daemon (`cmd/aespd/`)
- [ ] Documentation (`docs/`)
- [ ] Build / CI (`.github/`, `Makefile`)
- [ ] Other: <!-- specify -->

## Changes

<!-- Describe the changes in detail. What files were modified and why? -->

### Files Changed

<!-- List the key files changed (or use the GitHub auto-generated list) -->

## Testing

<!-- Describe the testing performed. Include test commands and results. -->

### Test Commands Run

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Linting
make lint
```

### Test Results

<!-- Paste test output or summarize results -->

- [ ] All unit tests pass (`make test`)
- [ ] Integration tests pass (`make test-integration`) — if applicable
- [ ] Linting passes (`make lint`)
- [ ] Code coverage maintained or improved

## Documentation

<!-- Mark documentation status -->

- [ ] Documentation updated (if user-facing changes)
- [ ] Architecture Decision Record added (if architectural change)
- [ ] API documentation updated (if API changes)
- [ ] README updated (if significant feature)
- [ ] No documentation changes needed

## Checklist

### Code Quality

- [ ] Code follows project coding standards (see [CONTRIBUTING.md](../CONTRIBUTING.md))
- [ ] Code is formatted (`make fmt`)
- [ ] No new linting warnings (`make lint`)
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Error handling is comprehensive
- [ ] Logging is appropriate and structured

### Testing

- [ ] Tests added for new functionality
- [ ] Tests cover edge cases
- [ ] Existing tests continue to pass
- [ ] Benchmarks added for performance-critical code (if applicable)

### Compatibility

- [ ] Changes are backward compatible
- [ ] Breaking changes are documented with migration path
- [ ] API changes follow versioning policy
- [ ] Database migrations are included (if schema changes)

### Security

- [ ] No secrets or credentials committed
- [ ] Input validation is appropriate
- [ ] Authorization checks are in place
- [ ] No new security vulnerabilities introduced

## Screenshots / Logs

<!-- If applicable, add screenshots, logs, or traces to help reviewers understand the changes -->

## Deployment Notes

<!-- Any special considerations for deployment? Database migrations? Configuration changes? -->

## Review Notes

<!-- Any specific areas you'd like reviewers to focus on? Questions or concerns? -->

---

By submitting this pull request, I confirm that my contribution is made under the terms of the MIT License and that I have read and agree to the [Contributing Guidelines](../CONTRIBUTING.md).
