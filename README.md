# hexa-go

A CLI that scaffolds Go services in a hexagonal (ports and adapters) layout, then keeps extending them:
`hexa-go add model Product ...` writes the model, repository, service, handler and tests, and wires them into
`main.go` and the routes.

```bash
go install github.com/erwinhermantodev/hexa-go@latest   # Go 1.24+

hexa-go generate my-api --module github.com/me/my-api
cd my-api && go mod tidy
hexa-go add model Product -f "Name:string::required" -f "Price:float64::required,gt=0"
docker-compose up -d postgres
go run main.go
```

## What you get

- **Layered structure:** `model/`, `repository/`, `service/`, `transport/http/` with interfaces between layers
- **REST with Echo:** CRUD routes per model, validation on create, one JSON response envelope, Swagger route
- **GORM:** PostgreSQL, MySQL or SQLite, selected by configuration
- **Wiring done for you:** new code is inserted at marker comments, imports included, and a failed command rolls
  back cleanly
- **Tests:** repository, service and handler tests for every model
- **Project extras:** Dockerfile, docker-compose with PostgreSQL, GitHub Actions workflow, Makefile, config
  loading, graceful shutdown, JWT and bcrypt helpers

## Commands

| Command | Description |
| --- | --- |
| `generate [name]` | Create a project (`--minimal` skips the default `User` model) |
| `add model [name]` | Model with repository, service, handler and tests (`-f` for fields) |
| `add module [name]` | The same pieces in one package under `internal/modules/` |
| `add service [name]` | Empty service, registered in `main.go` |
| `add handler [name]` | Empty handler, registered in `main.go` |

Fields use `Name:Type[:GormOptions[:Validation]]`, for example `-f "Slug:string:unique:required"`.
Existing files are never overwritten without `--force`.

## Documentation

- [DOCUMENTATION.md](DOCUMENTATION.md): full reference (commands, naming, field syntax, generated layout, wiring,
  known limitations)
- [ROADMAP.md](ROADMAP.md): what is planned
- [docs/index.html](docs/index.html): the same overview as a web page

## Development

```bash
go test ./... -short   # offline unit tests
go test ./...          # also builds and vets generated projects (needs network)
```

## License

MIT
