# AGENTS.md - mb-cli Development Guide

## Project Structure

mb-cli follows the standard Go project layout:
```
mb-cli/
├── cmd/mb/main.go         # Application entry point
├── internal/
│   ├── cli/               # CLI command implementations (Cobra)
│   ├── client/            # Metabase API client
│   ├── config/            # Configuration management
│   ├── formatter/         # Output formatting (JSON, table)
│   └── version/           # Version information
├── tests/                 # Test files
└── Makefile               # Build targets
```

## Build/Test Commands

- `make build` - Build binary to bin/mb
- `make test` - Run all tests
- `make test-verbose` - Run tests with verbose output
- `make fmt` - Format code with gofmt
- `make vet` - Static analysis with go vet
- `make lint` - Run golangci-lint
- `make clean` - Remove build artifacts
- `make deps` - Download and tidy dependencies
- `make build-all` - Cross-platform builds

## Configuration

Environment variables (both required):
- `MB_HOST` - Metabase instance URL (e.g. `https://your-metabase-instance.com`)
- `MB_API_KEY` - Metabase API key

## Code Style Guidelines

- Use Go standard formatting (gofmt)
- Package names: lowercase, single word
- Types: PascalCase
- Functions/methods: PascalCase (exported), camelCase (unexported)
- Variables: camelCase
- Constants: PascalCase or ALL_CAPS
- Import ordering: standard library, third-party, local packages
- Use `fmt.Errorf` for error wrapping
- Always defer `resp.Body.Close()` after error check for HTTP responses
- Use table-driven tests with `tests := []struct{}` pattern
- Test files go in the `tests/` directory
- Use `httptest.NewServer` for HTTP client testing

## Smoke Testing

- After implementing API methods or CLI commands, always smoke test **every** method against the real Metabase API (`make build && ./bin/mb-cli <command>`)
- Do not consider a step complete until all endpoints have been verified end-to-end, not just unit tested

## Project Style

- When you generate or update the CHANGELOG.md, be concise
- New additions go in [Unreleased] section
- Don't bump version unless requested
