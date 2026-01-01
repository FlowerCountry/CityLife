# Repository Guidelines

## Project Structure & Module Organization
- `cmd/citylife/` holds the entry point for the Gin-based API server.
- `internal/` contains game logic and API layers (e.g., `internal/api`, `internal/action`, `internal/world`, `internal/user`, `internal/save`).
- `tests/` hosts integration scripts such as `tests/run_test.sh`.
- `docs/`, `locale/`, and `bin/` store documentation, localization files, and build outputs.

## Build, Test, and Development Commands
Use the Makefile for common tasks:

- `make deps`: tidy Go module dependencies.
- `make build`: build the binary to `bin/citylife`.
- `make run`: build and run the server on the default port (8080).
- `make test`: run unit tests (`go test ./...`).
- `make fmt`: format Go code (`go fmt ./...`).
- `make vet`: run `go vet` checks.
- `make build-all`: cross-platform builds.

You can also run directly, e.g. `./bin/citylife --port 3000 --debug`.

## Coding Style & Naming Conventions
- Follow Go conventions: tabs for indentation, `CamelCase` for exported symbols, `camelCase` for internal ones.
- Keep functions small and single-purpose; avoid deep nesting.
- Format before committing with `make fmt`; keep comments concise (Chinese is acceptable).
- Test files should use `*_test.go` naming.

## Testing Guidelines
- Unit tests: `make test`.
- Integration tests: `./tests/run_test.sh` (requires `curl` and `jq`, starts a local server on port 18081).
- Add new tests for gameplay actions or API handlers you touch.

## Commit & Pull Request Guidelines
- Commit messages are concise single-line subjects; optional type prefixes like `feat:` or `chore:` are common.
- Chinese or English subjects are acceptable; include issue references when relevant.
- PRs should include: a short summary, tests run, and any API or save-data changes.

## Runtime & Data Notes
- Default server port is 8080; flags include `--port` and `--debug`.
- Save data is managed by the path manager; keep changes backward compatible.
