# P2 — Local-first Deployment Profile

**INV-11** · Single machine, offline-capable

## Bootstrap path

```bash
git clone https://github.com/kishoreHQ/AESP-Reference-Implementation
cd AESP-Reference-Implementation
go test ./...
go run ./cmd/aespd
go run ./cmd/aespd conformance
```

## Local provider plugins

Abstract categories (plugins, not kernel hardcodes):

- Local inference servers (Ollama-class, MLX-class, llama.cpp-class, LM Studio-class)
- Bundled `provider.mock-local` for zero-dependency demos

## Backends

| Concern | Backend |
|---------|---------|
| Memory | SQLite / filesystem |
| Artifacts | `./.aesp/artifacts` |
| Credentials | OS keychain or local encrypted file; **zero cloud credentials required** |

## Zero cloud credentials checklist

- [ ] No remote provider plugins registered  
- [ ] `provider.mock-local` or local inference only  
- [ ] Credential broker empty of cloud keys  
- [ ] Offline conformance subset passes (`go test ./...`)  
- [ ] Examples 01 and 09 (failover to local) runnable conceptually  
