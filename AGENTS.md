# Repository Guidelines

## Project Structure & Module Organization
- Module: `github.com/choral-io/gommerce-server-aio` (builds on `../gommerce-server-core`).
- Entrypoints: `cmd/server` (HTTP/gRPC app with fx), `cmd/dbseed` (database seeding).
- Services: `server/` with versioned APIs in `server/v1` and `server/v1beta`.
- Data layer: `data/models`, `data/repos`, `data/repos/pgsql`.
- Config and assets: YAML in `config/` (e.g., `app-local.yaml`); static files embedded from `static/`.
- Atlas schema: PostgreSQL schema files in `data/schema/` (e.g., `schema.pg.hcl`), managed via Atlas CLI.

## Build, Test, and Development Commands
- Go: 1.25+.
- Build: `go build ./...`. App binary: `CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/server ./cmd/server`.
- Run (local): `GOMMERCE_CONFIG_PATH=./config/app-local.yaml go run ./cmd/server`.
- Test: `go test ./... -v -cover`. Focused example: `go test -run Users -v ./server/...`.
- Vet/format: `go vet ./...`; `go fmt ./...`; `go mod tidy` before PRs.
- Schema and seed: `atlas schema inspect -u $DATA_SOURCE_URL > data/schema/schema.pg.hcl && atlas schema apply -u $DATA_SOURCE_URL --to "file://data/schema"`; `go run ./cmd/dbseed`.
- JWT keys: `openssl genrsa 2048 | tee >(openssl rsa -pubout 2>/dev/null)`.

## Coding Style & Naming Conventions
- Go style: `gofmt`. Package names are short lowercase nouns; exported identifiers use `CamelCase`.
- Indentation: default 4 spaces; YAML/JSON/Atlas HCL use 2 (see `.editorconfig`).
- Final newline: ensure files end with a newline (per `.editorconfig`).
- Errors: wrap with `%w`; define sentinels as `var ErrX = errors.New(...)`.
- APIs: accept `ctx context.Context` first; avoid long parameter lists.

## Testing Guidelines
- Framework: native `testing`. Files `*_test.go`; functions `TestXxx(t *testing.T)`.
- Prefer deterministic, table-driven tests and subtests (`t.Run`).
- Coverage: target repositories and service handlers. Run `go test -cover ./...`.

## Commit & Pull Request Guidelines
- Conventional commits: `feat:`, `fix:`, `refactor:`, `chore:`.
- Branch names: short kebab-case (e.g., `feat/server-auth`).
- PRs: clear description, motivation, linked issues. Include tests and note config/API changes. Ensure `go build ./...` and `go test ./...` pass.

## Security & Configuration Tips
- Environment: `GOMMERCE_ENVIRONMENT` (defaults to `development`). App loads `.env*` and `GOMMERCE_CONFIG_PATH`.
- Secrets: never commit private keys or `.env`. Rotate credentials seeded by `cmd/dbseed`.
- Endpoints configurable: DB, Redis, NATS, OTEL via env/config. Keep production values out of VCS.

## Architecture Overview
- Composition: the app wires core packages (`config`, `logging`, `otel`, `secure`, `server`) via `fx` in `cmd/server/main.go`.
- Interceptors and handlers come from core; versioned services live under `server/v1` and `server/v1beta`.
- Static files are served from embedded `static/`.
