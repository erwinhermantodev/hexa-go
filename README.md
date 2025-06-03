# Go Starter CLI - Restructured Project

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

```
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
│   │   ├── utils.go           # Utility templates
│   │   ├── server.go          # Server and config templates
│   │   ├── model.go           # Dynamic model templates
│   │   ├── repository.go      # Repository templates
│   │   ├── service.go         # Service templates
│   │   └── handler.go         # Handler templates
│   ├── prompts/                # Interactive prompts
│   │   └── prompts.go         # User input handling
│   └── utils/                  # Utility functions
│       ├── file.go            # File operations
│       └── parser.go          # Field parsing logic
└── examples/                   # Usage examples
    └── example-project/
```

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/erwinhermantodev/hexa-go
cd hexa-go

# Build the CLI
go build -o go-starter main.go

# Make it globally available (optional)
sudo mv go-starter /usr/local/bin/
```

### Generate a New Project

```bash
# Interactive mode (recommended for first-time users)
go-starter generate my-awesome-api --interactive

# Quick generation with flags
go-starter generate my-api \
  --module "github.com/erwinhermantodev/my-api" \
  --author "Your Name" \
  --description "My awesome API project"

# Minimal project (no authentication)
go-starter generate simple-api --minimal
```

### Add Components to Existing Project

```bash
# Navigate to your project directory
cd my-api

# Add a new model with full CRUD
go-starter add model Product \
  -f "Name:string::required,min=2" \
  -f "Price:float64::required,gt=0" \
  -f "Category:string::required" \
  -f "Description:string"

# Add model without repository
go-starter add model Category --no-repo

# Add standalone service
go-starter add service PaymentProcessor

# Add standalone handler
go-starter add handler HealthCheck
```

## 📖 Usage Examples

### Example 1: E-commerce API

```bash
# Generate the base project
go-starter generate ecommerce-api \
  --module "github.com/mycompany/ecommerce-api" \
  --author "Development Team" \
  --description "E-commerce REST API with authentication"

cd ecommerce-api

# Add Product model
go-starter add model Product \
  -f "Name:string::required,min=2,max=100" \
  -f "SKU:string::required,unique" \
  -f "Price:float64::required,gt=0" \
  -f "Stock:int::required,gte=0" \
  -f "CategoryID:uint::required" \
  -f "Description:string"

# Add Category model
go-starter add model Category \
  -f "Name:string::required,min=2,max=50" \
  -f "Slug:string::required,unique" \
  -f "Description:string"

# Add Order model
go-starter add model Order \
  -f "UserID:uint::required" \
  -f "TotalAmount:float64::required,gt=0" \
  -f "Status:string::required" \
  -f "ShippingAddress:string::required"

# Add payment service
go-starter add service Payment

# Add analytics handler
go-starter add handler Analytics
```

### Example 2: Blog API

```bash
# Generate minimal blog API
go-starter generate blog-api --minimal \
  --module "github.com/myblog/api" \
  --author "Blog Team"

cd blog-api

# Add Post model
go-starter add model Post \
  -f "Title:string::required,min=5,max=200" \
  -f "Slug:string::required,unique" \
  -f "Content:string::required,min=10" \
  -f "AuthorID:uint::required" \
  -f "Published:bool::default=false" \
  -f "PublishedAt:*time.Time"

# Add Comment model
go-starter add model Comment \
  -f "PostID:uint::required" \
  -f "AuthorName:string::required,min=2,max=50" \
  -f "AuthorEmail:string::required,email" \
  -f "Content:string::required,min=5,max=500" \
  -f "Approved:bool::default=false"

# Add search service
go-starter add service Search
```

### Example 3: Microservice

```bash
# Generate microservice with custom components
go-starter generate user-service \
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
├── cmd/server/              # Application entrypoint
│   └── main.go
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
go-starter generate my-project --interactive

cd my-project
go mod tidy
```

### 2. Add Models Incrementally
```bash
# Add core models first
go-starter add model User -f "Name:string::required" -f "Email:string::required,email"

# Add related models
go-starter add model Profile -f "UserID:uint::required" -f "Bio:string"
```

### 3. Add Business Logic
```bash
# Add services for complex business logic
go-starter add service NotificationService
go-starter add service AnalyticsService
```

### 4. Add Custom Handlers
```bash
# Add specialized handlers
go-starter add handler ReportHandler
go-starter add handler WebhookHandler
```

### 5. Run and Test
```bash
# Run the application
go run cmd/server/main.go

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