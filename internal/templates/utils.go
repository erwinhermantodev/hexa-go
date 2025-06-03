package templates

// UtilsCodesTemplate contains status codes and error codes
const UtilsCodesTemplate = `package utils

// HTTP Status Codes
const (
	StatusOK                   = 200
	StatusCreated              = 201
	StatusBadRequest           = 400
	StatusUnauthorized         = 401
	StatusForbidden            = 403
	StatusNotFound             = 404
	StatusInternalServerError  = 500
)

// Custom Error Codes
const (
	ErrCodeValidation          = "VALIDATION_ERROR"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrCodeUserExists          = "USER_EXISTS"
	ErrCodeInvalidCredentials  = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired        = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid        = "TOKEN_INVALID"
	ErrCodeDuplicateEntry      = "DUPLICATE_ENTRY"
	ErrCodeInvalidInput        = "INVALID_INPUT"
)

// Success Messages
const (
	MsgCreated    = "Resource created successfully"
	MsgUpdated    = "Resource updated successfully"
	MsgDeleted    = "Resource deleted successfully"
	MsgRetrieved  = "Resource retrieved successfully"
)
`

// UtilsConfigTemplate contains configuration utilities
const UtilsConfigTemplate = `package utils

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseHost     string ` + "`mapstructure:\"database_host\"`" + `
	DatabasePort     string ` + "`mapstructure:\"database_port\"`" + `
	DatabaseUser     string ` + "`mapstructure:\"database_user\"`" + `
	DatabasePassword string ` + "`mapstructure:\"database_password\"`" + `
	DatabaseName     string ` + "`mapstructure:\"database_name\"`" + `
	JWTSecret        string ` + "`mapstructure:\"jwt_secret\"`" + `
	JWTExpiry        string ` + "`mapstructure:\"jwt_expiry\"`" + `
	ServerPort       string ` + "`mapstructure:\"server_port\"`" + `
	ServerMode       string ` + "`mapstructure:\"server_mode\"`" + `
	GRPCPort         string ` + "`mapstructure:\"grpc_port\"`" + `
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("database_host", "localhost")
	viper.SetDefault("database_port", "5432")
	viper.SetDefault("database_user", "postgres")
	viper.SetDefault("database_password", "postgres")
	viper.SetDefault("database_name", "{{.Name}}")
	viper.SetDefault("jwt_secret", "your-secret-key")
	viper.SetDefault("jwt_expiry", "24h")
	viper.SetDefault("server_port", "8080")
	viper.SetDefault("server_mode", "debug")
	viper.SetDefault("grpc_port", "9090")

	// Override with environment variables
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		viper.Set("database_host", dbHost)
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		viper.Set("database_port", dbPort)
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		viper.Set("database_user", dbUser)
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		viper.Set("database_password", dbPassword)
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		viper.Set("database_name", dbName)
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		viper.Set("jwt_secret", jwtSecret)
	}
	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		viper.Set("server_port", serverPort)
	}
	if grpcPort := os.Getenv("GRPC_PORT"); grpcPort != "" {
		viper.Set("grpc_port", grpcPort)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
`

// UtilsJwtTemplate contains JWT utilities
const UtilsJwtTemplate = `package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   ` + "`json:\"user_id\"`" + `
	Email  string ` + "`json:\"email\"`" + `
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string ` + "`json:\"access_token\"`" + `
	RefreshToken string ` + "`json:\"refresh_token\"`" + `
	ExpiresAt    int64  ` + "`json:\"expires_at\"`" + `
}

type JWT struct {
	secret string
	expiry time.Duration
}

func NewJWT(secret string, expiry time.Duration) *JWT {
	return &JWT{
		secret: secret,
		expiry: expiry,
	}
}

func (j *JWT) GenerateTokenPair(userID uint, email string) (*TokenPair, error) {
	now := time.Now()
	expiresAt := now.Add(j.expiry)

	// Access token claims
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Create access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	// Create refresh token (longer expiry)
	refreshClaims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenJWT.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}

func (j *JWT) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
`

// UtilsMessagesTemplate contains message utilities
const UtilsMessagesTemplate = `package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Messages struct {
	data map[string]interface{}
}

func LoadMessages(locale string) (*Messages, error) {
	filename := fmt.Sprintf("locales/%s.json", locale)
	file, err := os.ReadFile(filename)
	if err != nil {
		// Fallback to English
		file, err = os.ReadFile("locales/en.json")
		if err != nil {
			return nil, err
		}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &Messages{data: data}, nil
}

func (m *Messages) Get(key string, args ...interface{}) string {
	keys := splitKey(key)
	value := m.getValue(keys, m.data)
	
	if value == "" {
		return key // Return key if not found
	}

	if len(args) > 0 {
		return fmt.Sprintf(value, args...)
	}

	return value
}

func (m *Messages) getValue(keys []string, data map[string]interface{}) string {
	if len(keys) == 0 {
		return ""
	}

	if len(keys) == 1 {
		if val, ok := data[keys[0]].(string); ok {
			return val
		}
		return ""
	}

	if nested, ok := data[keys[0]].(map[string]interface{}); ok {
		return m.getValue(keys[1:], nested)
	}

	return ""
}

func splitKey(key string) []string {
	var keys []string
	current := ""
	
	for _, char := range key {
		if char == '.' {
			if current != "" {
				keys = append(keys, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		keys = append(keys, current)
	}
	
	return keys
}
`

// UtilsPasswordTemplate contains password utilities
const UtilsPasswordTemplate = `package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash for the given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
`

// UtilsValidatorTemplate contains validation utilities
const UtilsValidatorTemplate = `package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()
	
	// Register custom tag name function
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	return &Validator{validate: validate}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, v.formatError(err))
		}
		return fmt.Errorf(strings.Join(errors, ", "))
	}
	return nil
}

func (v *Validator) formatError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
`
