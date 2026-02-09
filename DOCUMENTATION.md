# hexa-go CLI Documentation

`hexa-go` is a powerful CLI tool designed to scaffold Go projects using the **Hexagonal Architecture** (also known as Ports and Adapters). It provides a robust foundation for building scalable, maintainable, and testable applications with built-in support for REST APIs, gRPC, and PostgreSQL.

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Commands Reference](#commands-reference)
4. [Project Architecture](#project-architecture)
5. [Development Workflow](#development-workflow)
6. [Field Definitions & Validation](#field-definitions-validation)

---

## Installation

### From Source

Ensure you have Go 1.21+ installed.

```bash
git clone https://github.com/erwinhermantodev/hexa-go.git
cd hexa-go
go install
```

---

## Quick Start

Create your first project in minutes:

```bash
# Generate a new project
hexa-go generate my-awesome-api --interactive

# Go to project directory
cd my-awesome-api

# Install dependencies
go mod tidy

# Run the API
go run main.go
```

---

## Commands Reference

### `generate [project-name]`

Generates a new project structure.

**Flags:**

- `-m, --module`: Go module name (e.g., `github.com/user/repo`).
- `-a, --author`: Author name for documentation.
- `-d, --description`: Short project description.
- `-i, --interactive`: Enable interactive mode to define models and services during generation.
- `--minimal`: Generate a minimal project without the default Authentication module.

### `add model [name]`

Adds a new domain model with its repository, service, and HTTP handler.

**Flags:**

- `-f, --field`: Define model fields. Format: `FieldName:Type:Tag:Validation`.
  - Example: `-f "Title:string::required,min=5"`

### `add service [name]`

Adds a standalone business logic service.

### `add handler [name]`

Adds a standalone HTTP handler.

---

## Project Architecture

The generated project follows Hexagonal Architecture:

- `model/`: Domain entities and request/response structures.
- `repository/`: Data access layer (Adapters). Interface-driven to support multiple databases.
- `service/`: Business logic layer (Application Core). Contains the "source of truth" for your application.
- `transport/`: Entry points to the application.
  - `http/`: REST API handlers and routes using Echo.
  - `grpc/`: gRPC server and proto definitions.
- `utils/`: Cross-cutting concerns like JWT, validation, and configuration.

---

## Development Workflow

1. **Define Models**: Use `hexa-go add model` to create your data structures.
2. **Implement Logic**: Open the generated files in `service/` to add your business rules.
3. **Database Queries**: Customize the generated GORM queries in `repository/`.
4. **Wire it Up**:
   - Register your new repository and service in `main.go`.
   - Add your handler's routes in `transport/http/routes/routes.go`.

---

## Field Definitions & Validation

When using the `-f` flag, follow this structure:
`"Name:Type:Tag:Validation"`

### Supported Types

- `string`, `int`, `uint`, `int64`, `float64`, `bool`, `time.Time`

### Validation Tags (via go-playground/validator)

- `required`: Field must not be empty.
- `email`: Must be a valid email.
- `min=N`, `max=N`: Length or value constraints.
- `gt=N`, `gte=N`: Numerical comparisons.

### Example

```bash
hexa-go add model Product \
  -f "Name:string::required" \
  -f "Price:float64::required,gt=0" \
  -f "IsPublished:bool:default:false"
```

---

## 🌟 Advanced Features

### 🔐 Graceful Shutdown

The generated application handles OS signals (`SIGINT`, `SIGTERM`) to shutdown the HTTP server and database connections gracefully, ensuring no requests are lost.

### 📝 Structured Logging

Powered by `rs/zerolog`, the application provides structured JSON logs in production and human-friendly console output in development. Configure log levels via `LOG_LEVEL` environment variable.

### 🏗️ Automated Dependency Wiring

When you use `hexa-go add model`, the CLI automatically:

1. Generates the Model, Repository, Service, and Handler.
2. Registers the Repository, Service, and Handler in `main.go`.
3. Injects the REST API routes in `transport/http/routes/routes.go`.

### 📚 Swagger API Documentation

Integrated with `swaggo/swag`. Generate your API documentation with:

```bash
make swag-init
```

Access the UI at `http://localhost:8080/swagger/index.html`.

### 🗃️ Multi-Database Support

Support for PostgreSQL (default), MySQL, and SQLite. Change the driver via `DATABASE_DRIVER` environment variable.

### 🚀 CI/CD

A GitHub Actions workflow is automatically generated in `.github/workflows/ci.yml` for automated testing and building.
