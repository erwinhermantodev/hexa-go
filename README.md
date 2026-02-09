# HexaGo Starter CLI - Restructured Project

A powerful CLI tool for generating flexible Go projects with hexagonal architecture, featuring customizable models, repositories, services, and handlers.

## 🎯 Features

- **Hexagonal Architecture**: Follows Domain-Driven Design principles
- **Modular Structure**: Well-organized packages and separation of concerns
- **Template System**: Flexible template engine for code generation
- **Interactive Mode**: Guided project setup with prompts
- **Minimal Mode**: Generate lightweight projects without authentication
- **Incremental Development**: Add components to existing projects
- **Production Ready**: Includes Docker, configuration, logging, and more

## 📁 Project Structure

````
hexa-go/
├── main.go                     # Entry point
├── cmd/                        # CLI commands
│   ├── root.go                 # Root command configuration
│   ├── generate.go             # Project generation command
│   ├── add.go                  # Add component command
│   └── add_commands.go         # Model, service, handler commands
├── internal/                   # Internal packages
│   ├── config/                 # Configuration types
│   │   └── types.go           # ProjectConfig, ModelConfig, FieldConfig
│   ├── generator/              # Code generation logic
│   │   ├── generator.go       # Main generator struct
│   │   ├── project.go         # Project creation logic
│   │   ├── model.go           # Model generation
│   │   ├── service.go         # Service generation
│   │   └── handler.go         # Handler generation
│   ├── templates/              # Template definitions
│   │   ├── base.go            # Base project templates
│   │   ├── locale.go          # Localization templates
# hexa-go 🚀

`hexa-go` is a powerful CLI tool designed to scaffold Go projects using **Hexagonal Architecture**. It helps developers build scalable, maintainable, and production-ready applications by automating the creation of models, services, and handlers.

## ✨ Features

- **🏹 Hexagonal Architecture**: Separation of concerns between domain logic and infrastructure.
- **🔌 Multi-Transport**: Built-in support for HTTP (Echo) and gRPC.
- **🗃️ Database Ready**: Integrated with GORM and PostgreSQL.
- **🔐 Security**: Pre-configured JWT authentication and password hashing.
- **🌍 I18n**: Support for internationalization out-of-the-box.
- **🐳 Containerized**: Ready-to-use Dockerfile and docker-compose.
- **🛠️ Extensible**: Easily add models, services, and handlers via CLI.

## 📚 Documentation

- **[Full Documentation](DOCUMENTATION.md)**: Comprehensive guide on how to use `hexa-go`.
- **[Enhancement Recommendations](RECOMMENDATIONS.md)**: Ideas and roadmap for future improvements.

## 🚀 Quick Start

### Installation

```bash
go install github.com/erwinhermantodev/hexa-go
````

### Create a New Project

```bash
hexa-go generate my-project --interactive
```

---

## 🛠️ Commands Reference

| Command       | Description                                  |
| ------------- | -------------------------------------------- |
| `generate`    | Create a new project structure               |
| `add model`   | Add a domain model with Repo/Service/Handler |
| `add service` | Add a standalone business logic service      |
| `add handler` | Add a standalone HTTP handler                |

---

## 🏗️ Project Structure

```text
my-project/
├── cmd/                # Entry points
├── internal/           # Private application code
│   ├── model/          # Domain entities
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic
│   └── transport/      # HTTP/gRPC handlers
├── utils/              # Helper functions
└── ...
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
D:uint::required" \
 -f "TotalAmount:float64::required,gt=0" \
 -f "Status:string::required" \
 -f "ShippingAddress:string::required"

# Add payment service

hexa-go add service Payment

# Add analytics handler

hexa-go add handler Analytics

````

### Example 2: Blog API

```bash
# Generate minimal blog API
hexa-go generate blog-api --minimal \
  --module "github.com/myblog/api" \
  --author "Blog Team"

cd blog-api

# Add Post model
hexa-go add model Post \
  -f "Title:string::required,min=5,max=200" \
  -f "Slug:string::required,unique" \
  -f "Content:string::required,min=10" \
  -f "AuthorID:uint::required" \
  -f "Published:bool::default=false" \
  -f "PublishedAt:*time.Time"

# Add Comment model
hexa-go add model Comment \
  -f "PostID:uint::required" \
  -f "AuthorName:string::required,min=2,max=50" \
  -f "AuthorEmail:string::required,email" \
  -f "Content:string::required,min=5,max=500" \
  -f "Approved:bool::default=false"

# Add search service
hexa-go add service Search
````

### Example 3: Microservice

```bash
# Generate microservice with custom components
hexa-go generate user-service \
  --module "github.com/company/user-service" \
  --interactive

# During interactive mode:
# - Add Profile model with custom fields
# - Add Notification service
# - Add Metrics handler
```

## 🔧 Configuration

### Field Types Supported

- `string` - Text fields
- `int`, `uint`, `int64` - Integer fields
- `float64` - Decimal fields
- `bool` - Boolean fields
- `time.Time` - Timestamp fields
- `*time.Time` - Optional timestamp fields
- `gorm.DeletedAt` - Soft delete support
- Custom types for relationships

### Validation Tags

- `required` - Field is required
- `email` - Valid email format
- `min=N` - Minimum length/value
- `max=N` - Maximum length/value
- `gt=N` - Greater than value
- `gte=N` - Greater than or equal
- `lt=N` - Less than value
- `lte=N` - Less than or equal
- `unique` - Unique constraint (GORM)

### GORM Tags

- `primaryKey` - Primary key field
- `unique` - Unique constraint
- `not null` - Not null constraint
- `default:value` - Default value
- `index` - Create index
- `autoIncrement` - Auto increment

## 🏗️ Generated Project Structure

```
your-project/
|--  main.go                 # Application entrypoint
├── configs/                 # Configuration files
│   └── config.yaml
├── locales/                 # Internationalization
│   ├── en.json
│   └── id.json
├── model/                   # Data models (generated)
│   ├── user.go
│   └── product.go
├── repository/              # Data access layer (generated)
│   ├── interfaces.go
│   ├── user.go
│   └── product.go
├── service/                 # Business logic layer (generated)
│   ├── user.go
│   └── product.go
├── transport/               # Transport layer
│   ├── http/
│   │   ├── handler/         # HTTP handlers (generated)
│   │   │   ├── user_handler.go
│   │   │   └── product_handler.go
│   │   └── routes/          # Route definitions
│   │       └── routes.go
│   └── grpc/                # gRPC server
│       ├── server.go
│       └── run.go
├── utils/                   # Utility functions
│   ├── codes.go
│   ├── config.go
│   ├── jwt.go
│   ├── messages.go
│   ├── password.go
│   └── validator.go
├── migrations/              # Database migrations
├── docs/                    # Documentation
├── test/                    # Tests
├── go.mod
├── go.sum
├── README.md
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .gitignore
└── .env.example
```

## 🔄 Development Workflow

### 1. Initial Setup

```bash
# Generate project
hexa-go generate my-project --interactive

cd my-project
go mod tidy
```

### 2. Add Models Incrementally

```bash
# Add core models first
hexa-go add model User -f "Name:string::required" -f "Email:string::required,email"

# Add related models
hexa-go add model Profile -f "UserID:uint::required" -f "Bio:string"
```

### 3. Add Business Logic

```bash
# Add services for complex business logic
hexa-go add service NotificationService
hexa-go add service AnalyticsService
```

### 4. Add Custom Handlers

```bash
# Add specialized handlers
hexa-go add handler ReportHandler
hexa-go add handler WebhookHandler
```

### 5. Run and Test

```bash
# Run the application
go run main.go

# Or use Docker
docker-compose up -d
```

## 🌟 Advanced Features

### Custom Templates

You can modify templates in `internal/templates/` to customize the generated code:

- `base.go` - Project structure templates
- `model.go` - Model generation templates
- `repository.go` - Repository templates
- `service.go` - Service templates
- `handler.go` - Handler templates

### Environment Configuration

The generated projects support multiple configuration methods:

1. **YAML Config**: `configs/config.yaml`
2. **Environment Variables**: `.env` file
3. **Command Line**: Environment variable overrides

### Docker Support

Every generated project includes:

- Multi-stage Dockerfile for optimized images
- docker-compose.yml with PostgreSQL
- Production-ready container configuration

### Database Migrations

Use the included migration commands:

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 📚 [Documentation](https://github.com/erwinhermantodev/hexa-go/wiki)
- 🐛 [Issue Tracker](https://github.com/erwinhermantodev/hexa-go/issues)
- 💬 [Discussions](https://github.com/erwinhermantodev/hexa-go/discussions)

## 🏆 Examples in the Wild

- [E-commerce API](examples/ecommerce-api/)
- [Blog Platform](examples/blog-platform/)
- [User Management Service](examples/user-service/)
- [File Upload Service](examples/file-service/)

---

**Happy coding! 🚀**
