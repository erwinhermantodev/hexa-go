# hexa-go documentation

`hexa-go` scaffolds Go services in a hexagonal (ports and adapters) layout and keeps extending them: after the
first `generate`, `hexa-go add ...` writes new models, services, handlers and modules and wires them into
`main.go` and the route table for you.

## Contents

1. [Install](#install)
2. [Quick start](#quick-start)
3. [Commands](#commands)
4. [Naming rules](#naming-rules)
5. [Field syntax](#field-syntax)
6. [What gets generated](#what-gets-generated)
7. [How wiring works](#how-wiring-works)
8. [Generated project reference](#generated-project-reference)
9. [Known limitations](#known-limitations)
10. [Working on hexa-go itself](#working-on-hexa-go-itself)

## Install

Requires Go 1.24 or newer.

```bash
go install github.com/erwinhermantodev/hexa-go@latest
```

From a clone:

```bash
git clone https://github.com/erwinhermantodev/hexa-go.git
cd hexa-go
go install .
```

## Quick start

```bash
hexa-go generate my-api --module github.com/me/my-api --author "Me" --description "Orders API"
cd my-api
go mod tidy
hexa-go add model Product -f "Name:string::required" -f "Price:float64::required,gt=0"
docker-compose up -d postgres
go run main.go
curl localhost:8080/api/v1/health
```

`go mod tidy` resolves the latest dependency versions and sets the `go` directive in the new project's `go.mod`.
Any option you leave off `generate` is asked for interactively.

## Commands

### `hexa-go generate [project-name]`

Creates `./project-name`. Fails if the directory exists and is not empty, and removes the directory again if
generation fails part-way.

| Flag | Meaning |
| --- | --- |
| `-m, --module` | Go module path, for example `github.com/me/my-api` |
| `-a, --author` | Author shown in the generated README |
| `-d, --description` | One-line description |
| `-i, --interactive` | Prompt for extra models and services to generate up front |
| `--minimal` | Skip the default `User` model and `Auth` service |
| `--force` | Generate into an existing non-empty directory |

### `hexa-go add model [name]`

Adds a model with repository, service, handler and tests, and registers all of it. Run it in the project root
(the directory with `go.mod`).

| Flag | Meaning |
| --- | --- |
| `-f, --fields` | A field, repeatable. See [Field syntax](#field-syntax). Without any `-f` you are prompted. |
| `--no-repo` | Model only (also skips service and handler, which depend on it) |
| `--no-service` | Model and repository (also skips the handler) |
| `--no-handler` | Skip the HTTP handler and routes |
| `--force` | Overwrite existing files and register again |

### `hexa-go add service [name]` and `hexa-go add handler [name]`

Add an empty service (`service/<name>.go`) or handler (`transport/http/handler/<name>_handler.go`) and register it
in `main.go`. A trailing `Service` or `Handler` in the name is dropped, so `NotificationService` and
`Notification` are the same. The handler is created with a `Health` method as a starting point.

### `hexa-go add module [name]`

Adds a feature module with its model, repository, service and handler together in `internal/modules/<name>/`
(one Go package), and registers it. Takes `-f` and `--force` like `add model`. Module names that would clash
with packages the generated project already imports (`utils`, `routes`, `handler`, `service`, `model` and so on)
are rejected.

Use `add model` for the layered layout (`model/`, `repository/`, `service/`) and `add module` when you prefer to
keep one feature in one package.

### Failure behavior

- Commands exit non-zero and print `❌ <reason>`.
- A command that fails after writing files removes the files it created and restores `main.go` and `routes.go`.
- Existing files are never overwritten without `--force`.
- Running the same command twice does not duplicate registrations.

## Naming rules

Names may be given in any style: `order-item`, `order_item`, `OrderItem`, `orderItem` and `Order Item` are all
the same name. Allowed: letters and digits, separated by `-`, `_` or spaces, starting with a letter. Go keywords
are rejected.

| Given `order-item` | Used for |
| --- | --- |
| `OrderItem` | Go types and exported names (`OrderItemHandler`, `NewOrderItemService`) |
| `orderItem` | local variables in `main.go` (`orderItemRepo`) |
| `order_item` | file names (`model/order_item.go`) |
| `order-items` | URL path (`/api/v1/order-items`) |
| `orderitem` | package name for modules |

Plurals follow English rules (`Category` becomes `/categories` and `GetAllCategories`, `Address` becomes
`/addresses`) with a short list of irregular nouns (`person`, `child`, `status`, ...).

## Field syntax

```
Name:Type[:GormOptions[:Validation]]
```

| Slot | Meaning |
| --- | --- |
| `Name` | Go field name; the first letter is upper-cased |
| `Type` | `string`, `int`, `uint`, `int64`, `float64`, `bool`, `time.Time`, `gorm.DeletedAt`, a pointer (`*time.Time`), a slice (`[]string`) or a custom type name |
| `GormOptions` | Options for the `gorm` tag separated by `;`. A value is written with `=` because `:` separates the slots: `default=false` becomes `default:false`. |
| `Validation` | go-playground/validator tags separated by `,`: `required`, `email`, `min=N`, `max=N`, `gt=N`, `gte=N`, `lt=N`, `lte=N` and the rest of that library |

Every field also gets a `json` tag in snake case (`IsPublished` becomes `is_published`).

```bash
hexa-go add model Post \
  -f "Title:string::required,min=5,max=200" \
  -f "Slug:string:unique:required" \
  -f "Published:bool:default=false" \
  -f "PublishedAt:*time.Time" \
  -f "AuthorID:uint:index:required"
```

`ID`, `CreatedAt`, `UpdatedAt` and `DeletedAt` (soft delete) are added automatically unless you define them.
In create and update requests those four fields are left out, so clients cannot set them; a field hidden from
responses with `json:"-"` is still accepted in requests.

Mistakes are reported before anything is written, for example:

```
❌ invalid field "Slug:string::unique": "unique" is a GORM option; put it in the third slot (Name:Type:unique)
```

When prompted for fields interactively, type one spec per line using the same syntax, and an empty line to finish.

## What gets generated

`hexa-go generate my-api` produces:

```
my-api/
├── main.go                         config, database, validator, wiring, server, graceful shutdown
├── go.mod
├── Makefile  Dockerfile  docker-compose.yml  .gitignore  .env.example
├── .github/workflows/ci.yml        build and test on push and pull request
├── configs/config.yaml
├── locales/en.json  id.json
├── docs/docs.go                    placeholder for `swag init`
├── model/user.go                   entity, create request, response
├── repository/user.go              GORM repository (+ user_test.go)
├── service/user.go  auth.go        business logic (+ user_test.go); auth.go is an empty stub
├── transport/
│   ├── http/handler/user_handler.go  (+ user_handler_test.go)
│   ├── http/routes/routes.go       Router struct, /swagger, /api/v1/health, model routes
│   └── grpc/server.go  run.go      gRPC skeleton, not started by main.go
├── utils/                          config, database, jwt, password, validator, messages, codes, response
└── migrations/                     empty, for golang-migrate
```

`--minimal` leaves out `user.go`, its tests, and `auth.go`. Modules add `internal/modules/<name>/`.

### Dependency direction

```
HTTP handler  →  service  →  repository interface  →  GORM adapter  →  database
(transport)      (core)       (port)                    (adapter)
```

Handlers depend on service interfaces, services on repository interfaces, and `main.go` is the only place that
builds the concrete objects.

## How wiring works

`main.go` and `routes.go` contain marker comments:

| Marker | File | Receives |
| --- | --- | --- |
| `// [REPOS-INIT]` | main.go | `xRepo := repository.NewXRepository(db)` |
| `// [SERVICES-INIT]` | main.go | `xService := service.NewXService(xRepo)` |
| `// [HANDLERS-INIT]` | main.go | `xHandler := handler.NewXHandler(xService, validator)` |
| `// [ROUTER-HANDLERS-INIT]` | main.go | `router.XHandler = xHandler` |
| `// [HANDLER-FIELDS-EXPORTED]` | routes.go | the `XHandler` field on `Router` |
| `// [ROUTES-INIT]` | routes.go | the route group for the model |
| `// [IMPORTS]` | both | (imports are added through the AST) |

New code is inserted on the line above a marker and the file is re-formatted with `gofmt`. Imports are added
without duplicates. Leave the markers in place; without them `hexa-go add` fails and rolls back. You can freely
edit everything else, including the code that was already inserted.

## Generated project reference

### Configuration

Values come from `configs/config.yaml`, then environment variables. Every key can be set as `APP_<KEY>`
(`APP_SERVER_PORT=9000`, `APP_DATABASE_DRIVER=sqlite`, `APP_LOG_LEVEL=debug`). These shorter names also work:
`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`, `SERVER_PORT`, `GRPC_PORT`.
`.env.example` is only a reference: nothing loads a `.env` file.

### Database

PostgreSQL by default; `APP_DATABASE_DRIVER` selects `mysql` or `sqlite` (SQLite uses `<database_name>.db`).
`utils.InitDB` opens the connection. No tables are created automatically; see [Known limitations](#known-limitations).

### REST API

Routes are registered under `/api/v1`. Each model with a handler gets create, list, get, update and delete
routes; every response uses `{"status": <code>, "message": "...", "data": ...}`. Create requests are validated.

### Swagger

`GET /swagger/*` is registered. Run `make swag-init` (needs `go install github.com/swaggo/swag/cmd/swag@latest`)
to generate the spec into `docs/`. Handler annotations are not generated, so the spec starts out empty.

### Tests

Each model gets a repository test (in-memory SQLite, needs cgo), a service test with a mocked repository, and a
handler scaffold. `make test` runs them; the CI workflow runs build and test with the Go version from `go.mod`.

### Docker

`docker-compose up -d` starts the app and PostgreSQL. The Dockerfile is a multi-stage build onto Alpine.

### Shutdown

On `SIGINT` or `SIGTERM` the server stops accepting requests and waits up to 10 seconds for in-flight ones.

## Known limitations

These are gaps in what the generator produces today; see [ROADMAP.md](ROADMAP.md).

- **No schema creation.** The generated app never calls `AutoMigrate`; create tables yourself or use migrations.
- **No auth endpoints.** `utils/jwt.go` and `utils/password.go` exist, and `service/auth.go` is an empty stub.
  The default `User` service stores the password exactly as received; hash it with `utils.HashPassword`.
- **`.env` is not loaded**, and `APP_LOG_LEVEL` is read into the config but not applied to the logger.
- **Update requests are not validated**, only create requests.
- **gRPC is a skeleton:** no `.proto` files, and `main.go` does not start the server.
- **Locales are not used by handlers**, which return English messages. `utils.LoadMessages` is available.
- Plurals use English rules and cannot be overridden.

## Working on hexa-go itself

```
main.go, cmd/            Cobra commands (generate, add model|service|handler|module)
internal/config/         ProjectConfig, ModelConfig, FieldConfig
internal/naming/         Pascal/camel/snake/kebab/plural helpers and name validation
internal/prompts/        interactive input
internal/utils/          field parser, go.mod reader, AST/text injection helpers
internal/generator/      generators, rollback, and the embedded templates
  templates/base/        files for a new project
  templates/core/        layered model, repository, service, handler and tests
  templates/modules/     co-located module files
```

Templates are Go `text/template` files embedded in the binary, with the helpers `camel`, `snake`, `kebab`,
`plural`, `pascal`, `reqTag`, `writable` and `contains`. To add a generated file to every new project, add the
template under `templates/base/` and an entry to `baseFiles` in `internal/generator/project.go`.

```bash
go test ./... -short      # fast, offline
go test ./...             # also builds and vets generated projects; needs network and a Go toolchain
```

The slow test generates the default and `--minimal` projects with a model, service, handler and module, then runs
`go mod tidy`, `go build ./...` and `go vet ./...` on them.
