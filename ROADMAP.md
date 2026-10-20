# Roadmap

What `hexa-go` generates today is described in [DOCUMENTATION.md](DOCUMENTATION.md). This file tracks what is
missing. Items are ordered by how much they matter to someone running the generated project.

## Done

- Tables created at startup with `db.AutoMigrate` for every model and module (`[MIGRATE]` marker)
- Generated Go files are gofmt-formatted
- Automated dependency wiring through marker comments, with rollback on failure
- Generated tests: repository (SQLite), service (mocked repository), handler scaffold
- Graceful shutdown on `SIGINT` and `SIGTERM`
- PostgreSQL, MySQL and SQLite through `database_driver`
- Swagger route and `make swag-init`
- Multi-stage Dockerfile, docker-compose with PostgreSQL, GitHub Actions workflow
- One field syntax with validation, and one naming scheme (case, files, URLs, plurals)
- Generator smoke test: generated projects are built and vetted in CI-style tests

## Next: make the generated app run end to end

1. **Load `.env`.** Read `.env` at startup, or drop `.env.example` and document `docker-compose` `env_file`.
2. **Apply `log_level`** to the zerolog logger, and emit JSON in production.
3. **Auth endpoints.** Register, login, refresh and logout on top of `utils/jwt.go`; hash passwords in the `User`
   service (`utils.HashPassword` exists but is never called).
4. **Validate update requests**, not only create requests.

## Then: depth

- **Swagger annotations** on generated handlers and a `@title` block in `main.go`, so `swag init` produces a real spec
- **gRPC:** generate `.proto` files per model and start the server from `main.go`, or remove the skeleton
- **Observability:** OpenTelemetry tracing and a Prometheus `/metrics` endpoint
- **Real handler and repository tests:** fill in the scaffold with a mocked service and example payloads
- **Typed configuration:** a validated `Config` struct with development, staging and production profiles
- **`healthcheck`** in `docker-compose.yml`
- **Use the message catalogs:** return localized messages from handlers via `utils.LoadMessages`
- **Relationships:** `belongs-to` and `has-many` field syntax with the matching GORM tags and preloads

## Ideas, not committed

- MongoDB or another non-SQL adapter behind the repository interfaces
- A `hexa-go migrate` subcommand in the generated project
- Pluralization overrides in the command line (`--plural`)
- Custom template directories for teams with their own conventions (would need a trust model before running any
  user-supplied code)
