package templates

const GoModTemplate = `module {{.ModuleName}}

go 1.21

require (
	github.com/labstack/echo/v4 v4.13.4
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
	github.com/go-playground/validator/v10 v10.15.1
	github.com/rs/zerolog v1.31.0
	github.com/swaggo/echo-swagger v1.4.1
	github.com/swaggo/swag v1.16.2
	golang.org/x/crypto v0.31.0
	gorm.io/gorm v1.25.12
	gorm.io/driver/postgres v1.5.11
	gorm.io/driver/mysql v1.5.2
	gorm.io/driver/sqlite v1.5.4
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.1
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240513163218-0867130af1f8
)`

// ReadmeTemplate is the template for README.md file
const ReadmeTemplate = `# {{.Name}}

{{.Description}}

**Author:** {{.Author}}

## Features

- 🔐 JWT Authentication
- 👤 User Management  
- 🌐 REST API with Echo
- 🔧 gRPC Support
- 🗃️ PostgreSQL with GORM
- 🐳 Docker Support
- 🌍 Internationalization
- ✅ Input Validation
- 📝 Configuration Management
- 🎯 **Flexible Model Generation**

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL
- Docker (optional)

### Installation

1. Clone the repository:
   ` + "```bash" + `
   git clone {{.ModuleName}}
   cd {{.Name}}
   ` + "```" + `

2. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

3. Set up environment variables:
   ` + "```bash" + `
   cp .env.example .env
   # Edit .env with your configuration
   ` + "```" + `

4. Run the application:
   ` + "```bash" + `
   go run main.go
   ` + "```" + `

### Using Docker

` + "```bash" + `
docker-compose up -d
` + "```" + `

## CLI Commands

This project includes a powerful CLI for generating components:

### Add New Model
` + "```bash" + `
go run main.go add model Product -f "Name:string::required" -f "Price:float64::required,gt=0" -f "Description:string"
` + "```" + `

### Add Service Only
` + "```bash" + `
go run main.go add service Payment
` + "```" + `

### Add Handler Only
` + "```bash" + `
go run main.go add handler Analytics
` + "```" + `

### Generate New Project
` + "```bash" + `
go run main.go generate my-new-project --interactive
` + "```" + `

## API Endpoints

### Generated for each model:
- ` + "`POST /api/v1/{model}s`" + ` - Create new {model}
- ` + "`GET /api/v1/{model}s`" + ` - Get all {models}
- ` + "`GET /api/v1/{model}s/{id}`" + ` - Get {model} by ID  
- ` + "`PUT /api/v1/{model}s/{id}`" + ` - Update {model}
- ` + "`DELETE /api/v1/{model}s/{id}`" + ` - Delete {model}

### Authentication (if included):
- ` + "`POST /api/v1/auth/register`" + ` - Register new user
- ` + "`POST /api/v1/auth/login`" + ` - User login
- ` + "`POST /api/v1/auth/refresh`" + ` - Refresh token
- ` + "`POST /api/v1/auth/logout`" + ` - User logout

## Project Structure

` + "```" + `
{{.Name}}/
├── configs/             # Configuration files
├── locales/             # Internationalization files
├── model/               # Data models (generated)
├── repository/          # Data access layer (generated)
├── service/             # Business logic layer (generated)
├── transport/           # Transport layer (HTTP/gRPC)
│   ├── http/
│   │   ├── handler/     # HTTP handlers (generated)
│   │   └── routes/      # Route definitions
│   └── grpc/            # gRPC server
├── utils/               # Utility functions
├── migrations/          # Database migrations
└── docs/               # Documentation
` + "```" + `

## Model Field Types

Supported field types for model generation:

- ` + "`string`" + ` - Text fields
- ` + "`int`, `uint`, `int64`" + ` - Integer fields  
- ` + "`float64`" + ` - Decimal fields
- ` + "`bool`" + ` - Boolean fields
- ` + "`time.Time`" + ` - Timestamp fields
- ` + "`gorm.DeletedAt`" + ` - Soft delete support
- Custom types for relationships

## Validation Tags

Supported validation tags:

- ` + "`required`" + ` - Field is required
- ` + "`email`" + ` - Valid email format
- ` + "`min=N`" + ` - Minimum length/value
- ` + "`max=N`" + ` - Maximum length/value  
- ` + "`gt=N`" + ` - Greater than value
- ` + "`gte=N`" + ` - Greater than or equal

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.
`

// DockerfileTemplate is the template for Dockerfile
const DockerfileTemplate = `FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/locales ./locales

EXPOSE 8080

CMD ["./main"]
`

// DockerComposeTemplate is the template for docker-compose.yml
const DockerComposeTemplate = `version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME={{.Name}}
    depends_on:
      - postgres
    volumes:
      - ./configs:/root/configs
      - ./locales:/root/locales

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB={{.Name}}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
`

// GitignoreTemplate is the template for .gitignore
const GitignoreTemplate = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
main

# Test binary
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# Environment variables
.env

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Database
*.db
*.sqlite

# Temporary files
tmp/
temp/
`

// EnvExampleTemplate is the template for .env.example
const EnvExampleTemplate = `# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME={{.Name}}

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY=24h

# Server
SERVER_PORT=8080
SERVER_MODE=debug

# GRPC
GRPC_PORT=9090
`

// MakefileTemplate is the template for Makefile
const MakefileTemplate = `.PHONY: build run test clean docker-up docker-down generate

# Build the application
build:
	go build -o bin/main main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Generate new model
generate-model:
	@echo "Usage: make generate-model MODEL=Product FIELDS='Name:string::required,Price:float64::required'"
	@if [ -n "$(MODEL)" ]; then \
		go run main.go add model $(MODEL) $(if $(FIELDS),-f $(shell echo $(FIELDS) | sed 's/,/ -f /g'),); \
	fi

# Generate service
generate-service:
	@echo "Usage: make generate-service SERVICE=Payment"
	@if [ -n "$(SERVICE)" ]; then \
		go run main.go add service $(SERVICE); \
	fi

# Generate handler  
generate-handler:
	@echo "Usage: make generate-handler HANDLER=Analytics"
	@if [ -n "$(HANDLER)" ]; then \
		go run main.go add handler $(HANDLER); \
	fi

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker build -t {{.Name}} .

# Database migrations
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/{{.Name}}?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/{{.Name}}?sslmode=disable" down

# Generate Swagger documentation
swag-init:
	swag init -g main.go

# Development
swag-dev:
	swag init -g main.go
	air
`
