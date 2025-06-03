// File: internal/templates/server.go

package templates

// MainServerTemplate is the main server application template
const MainServerTemplate = `package main

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"{{.ModuleName}}/transport/http/routes"
	"{{.ModuleName}}/utils"
)

func main() {
	// Load configuration
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	// db, err := connectDB(config)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }

	// log.Println(db)
	// Auto migrate your models here
	// if err := db.AutoMigrate(&model.User{}, &model.Product{}); err != nil {
	//     log.Fatalf("Failed to migrate database: %v", err)
	// }

	// Initialize JWT
	expiry, _ := time.ParseDuration(config.JWTExpiry)
	jwt := utils.NewJWT(config.JWTSecret, expiry)

	// Initialize validator
	validator := utils.NewValidator()
	log.Println(validator)
	// Initialize repositories
	// userRepo := repository.NewUserRepository(db)
	// productRepo := repository.NewProductRepository(db)

	// Initialize services
	// userService := service.NewUserService(userRepo)
	// productService := service.NewProductService(productRepo)

	// Initialize handlers
	// userHandler := handler.NewUserHandler(userService, validator)
	// productHandler := handler.NewProductHandler(productService, validator)

	// Setup Echo
	e := echo.New()

	// Configure Echo
	e.HideBanner = true
	if config.ServerMode == "release" {
		e.Debug = false
	} else {
		e.Debug = true
	}

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Setup routes
	routes.SetupRoutes(e, jwt)

	// Start server
	log.Printf("Server starting on port %s", config.ServerPort)
	log.Printf("Health check: http://localhost:%s/api/v1/health", config.ServerPort)
	
	if err := e.Start(":" + config.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectDB(config *utils.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.DatabaseHost,
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseName,
		config.DatabasePort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
`

// ConfigTemplate is the configuration file template
const ConfigTemplate = `# Database Configuration
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: {{.Name}}

# JWT Configuration
jwt:
  secret: your-super-secret-jwt-key
  expiry: 24h

# Server Configuration
server:
  port: 8080
  mode: debug

# gRPC Configuration
grpc:
  port: 9090

# Logging
log:
  level: info
  format: json
`

// RepositoryInterfacesTemplate contains repository interfaces
const RepositoryInterfacesTemplate = `package repository

import (
	"{{.ModuleName}}/model"
)

// Base repository interface for common CRUD operations
type BaseRepository[T any] interface {
	Create(entity *T) error
	GetByID(id uint) (*T, error)
	GetAll() ([]T, error)
	Update(entity *T) error
	Delete(id uint) error
}

// User repository interface
type UserRepository interface {
	BaseRepository[model.User]
	// GetByEmail(email string) (*model.User, error)
}

// Auth repository interface  
// type AuthRepository interface {
// 	CreateSession(session *model.Session) error
// 	GetSessionByRefreshToken(refreshToken string) (*model.Session, error)
// 	DeleteSession(refreshToken string) error
// 	DeleteUserSessions(userID uint) error
// }

// Add your custom repository interfaces here
// Example:
// type ProductRepository interface {
//     BaseRepository[model.Product]
//     GetByCategory(category string) ([]model.Product, error)
// }
`

const HttpRoutesTemplate = `package routes

import (
	"net/http"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	// "{{.ModuleName}}/transport/http/handler"
	"{{.ModuleName}}/utils"
)

func SetupRoutes(e *echo.Echo, jwtUtil *utils.JWT) {
	api := e.Group("/api/v1")
	
	// Health check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"service": "{{.Name}}",
			"version": "1.0.0",
		})
	})

	// Add your routes here
	// Example for generated models:
	// setupUserRoutes(api, userHandler, jwtUtil)
	// setupProductRoutes(api, productHandler, jwtUtil)
}

// Example route setup for a model
// func setupUserRoutes(api *echo.Group, userHandler *handler.UserHandler, jwtUtil *utils.JWT) {
//     users := api.Group("/users")
//     users.POST("", userHandler.CreateUser)
//     users.GET("", userHandler.GetAllUsers)
//     users.GET("/:id", userHandler.GetUser)
//     users.PUT("/:id", userHandler.UpdateUser)
//     users.DELETE("/:id", userHandler.DeleteUser)
// }

// func JWTAuthMiddleware(jwtUtil *utils.JWT) echo.MiddlewareFunc {
// 	return middleware.JWTWithConfig(middleware.JWTConfig{
// 		SigningKey:  []byte(jwtUtil.GetSecret()),
// 		TokenLookup: "header:Authorization",
// 		AuthScheme:  "Bearer",
// 		Claims:      &utils.JWTClaims{},
// 		ErrorHandler: func(err error) error {
// 			return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
// 				"error": "Invalid or expired token",
// 			})
// 		},
// 		SuccessHandler: func(c echo.Context) {
// 			token := c.Get("user").(*jwt.Token)
// 			if claims, ok := token.Claims.(*utils.JWTClaims); ok {
// 				c.Set("user_id", claims.UserID)
// 				c.Set("email", claims.Email)
// 			}
// 		},
// 	})
// }
`

// GrpcServerTemplate contains gRPC server setup
const GrpcServerTemplate = `package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"{{.ModuleName}}/utils"
)

type Server struct {
	config *utils.Config
}

func NewServer(config *utils.Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+s.config.GRPCPort)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	
	// Register your gRPC services here
	// pb.RegisterYourServiceServer(grpcServer, yourServiceImpl)

	log.Printf("gRPC server listening on port %s", s.config.GRPCPort)
	return grpcServer.Serve(lis)
}
`

// GrpcRunTemplate contains gRPC run logic
const GrpcRunTemplate = `package grpc

import (
	"log"
	"{{.ModuleName}}/utils"
)

func Run(config *utils.Config) {
	server := NewServer(config)
	
	log.Println("Starting gRPC server...")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
`

// DynamicServiceTemplate generates service files for models
const DynamicServiceTemplate = `package service

import (
	"{{.Config.ModuleName}}/model"
	"{{.Config.ModuleName}}/repository"
)

type {{.Model.Name}}Service struct {
	{{ToLower .Model.Name}}Repo repository.{{.Model.Name}}Repository
}

func New{{.Model.Name}}Service({{ToLower .Model.Name}}Repo repository.{{.Model.Name}}Repository) *{{.Model.Name}}Service {
	return &{{.Model.Name}}Service{
		{{ToLower .Model.Name}}Repo: {{ToLower .Model.Name}}Repo,
	}
}

func (s *{{.Model.Name}}Service) Create(req *model.{{.Model.Name}}Request) (*model.{{.Model.Name}}Response, error) {
	{{ToLower .Model.Name}} := &model.{{.Model.Name}}{
{{- range .Model.Fields }}
{{- if and (ne .Name "ID") (ne .Name "CreatedAt") (ne .Name "UpdatedAt") (ne .Name "DeletedAt") }}
		{{.Name}}: req.{{.Name}},
{{- end }}
{{- end }}
	}

	if err := s.{{ToLower .Model.Name}}Repo.Create({{ToLower .Model.Name}}); err != nil {
		return nil, err
	}

	return {{ToLower .Model.Name}}.ToResponse(), nil
}

func (s *{{.Model.Name}}Service) GetByID(id uint) (*model.{{.Model.Name}}Response, error) {
	{{ToLower .Model.Name}}, err := s.{{ToLower .Model.Name}}Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return {{ToLower .Model.Name}}.ToResponse(), nil
}

func (s *{{.Model.Name}}Service) GetAll() ([]model.{{.Model.Name}}Response, error) {
	{{ToLower .Model.Name}}s, err := s.{{ToLower .Model.Name}}Repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []model.{{.Model.Name}}Response
	for _, {{ToLower .Model.Name}} := range {{ToLower .Model.Name}}s {
		responses = append(responses, *{{ToLower .Model.Name}}.ToResponse())
	}

	return responses, nil
}

func (s *{{.Model.Name}}Service) Update(id uint, req *model.{{.Model.Name}}Request) (*model.{{.Model.Name}}Response, error) {
	{{ToLower .Model.Name}}, err := s.{{ToLower .Model.Name}}Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

{{- range .Model.Fields }}
{{- if and (ne .Name "ID") (ne .Name "CreatedAt") (ne .Name "UpdatedAt") (ne .Name "DeletedAt") }}
	{{ToLower $.Model.Name}}.{{.Name}} = req.{{.Name}}
{{- end }}
{{- end }}

	if err := s.{{ToLower .Model.Name}}Repo.Update({{ToLower .Model.Name}}); err != nil {
		return nil, err
	}

	return {{ToLower .Model.Name}}.ToResponse(), nil
}

func (s *{{.Model.Name}}Service) Delete(id uint) error {
	return s.{{ToLower .Model.Name}}Repo.Delete(id)
}
`

// CustomServiceTemplate generates standalone service files
const CustomServiceTemplate = `package service

import (
	"{{.Config.ModuleName}}/repository"
)

type {{.ServiceName}}Service struct {
	// Add your repository dependencies here
	// Example: userRepo repository.UserRepository
}

func New{{.ServiceName}}Service() *{{.ServiceName}}Service {
	return &{{.ServiceName}}Service{
		// Initialize your dependencies here
	}
}

// Add your business logic methods here
// Example:
// func (s *{{.ServiceName}}Service) DoSomething() error {
//     // Business logic implementation
//     return nil
// }
`
