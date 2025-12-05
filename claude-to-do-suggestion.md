# Go Booking System - Refactoring Plan & Best Practices

**Date**: 2025-12-04
**Current Scalability Score**: 3.4/10
**Target Score**: 7-8/10
**Estimated Time**: 2 weeks

---

## 📊 Current Assessment

### ✅ What's Working Well
1. Basic MVC separation (Controllers, Models, Routes)
2. Industry-standard libraries (Gin, GORM, JWT)
3. Security awareness (bcrypt password hashing, JWT tokens)
4. Environment configuration (.env files)
5. Database migrations (AutoMigrate)
6. Soft deletes (DeletedAt)
7. Git version control

### ❌ Critical Issues to Fix

#### 1. Fat Controllers - No Service Layer
**Problem**: Controllers handle validation, database queries, business logic, and response formatting.
**Impact**: Unmaintainable as features grow, hard to test, violates single responsibility.

#### 2. Direct Database Access in Controllers
**Problem**: `config.DB.Where()` called directly in controllers.
**Impact**: Can't mock for testing, tight coupling, hard to maintain.

#### 3. Silent Error Handling
**Problem**: Errors silently ignored (country lookup in signup).
**Impact**: Debugging production issues is difficult.

#### 4. No DTO Layer
**Problem**: Request structs mixed with domain models.
**Impact**: Frontend changes force database schema changes.

#### 5. Inconsistent API Responses
**Problem**: Different response formats across endpoints.
**Impact**: Frontend integration is harder.

#### 6. No Testing
**Problem**: No unit or integration tests visible.
**Impact**: Fear of refactoring, bugs in production.

#### 7. Minimal Logging
**Problem**: No structured logging.
**Impact**: Can't debug production issues effectively.

---

## 🎯 Refactoring Roadmap (13 Steps)

### Phase 1: Foundation (~3 hours)

#### **Step 1: Create Project Structure**

**New Structure**:
```
go-booking-system/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/                     # Private application code
│   ├── domain/                   # Business entities (formerly models/)
│   │   ├── user.go
│   │   └── country.go
│   ├── dto/                      # Data Transfer Objects
│   │   ├── requests.go
│   │   └── responses.go
│   ├── repository/               # Database access layer
│   │   ├── user_repository.go
│   │   └── country_repository.go
│   ├── service/                  # Business logic layer
│   │   ├── auth_service.go
│   │   └── user_service.go
│   ├── handler/                  # HTTP handlers (formerly controllers/)
│   │   ├── auth_handler.go
│   │   └── health_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── error_handler.go
│   │   └── logger.go
│   └── routes/
│       └── routes.go
├── pkg/                          # Public/reusable code
│   ├── logger/
│   │   └── logger.go
│   ├── response/
│   │   └── response.go
│   └── errors/
│       └── errors.go
├── config/
│   ├── config.go
│   └── database.go
├── tests/
│   ├── unit/
│   └── integration/
├── docs/                         # Swagger generated
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

**Migration Steps**:
1. Create new folders
2. Move `models/` → `internal/domain/`
3. Move `controllers/` → `internal/handler/` (will refactor later)
4. Move `routes/` → `internal/routes/`
5. Update all imports

---

#### **Step 2: Create DTO Layer**

**File**: `internal/dto/requests.go`
```go
package dto

// SignUpRequest represents user registration payload
type SignUpRequest struct {
    Email    string `json:"email" binding:"required,email" example:"user@example.com"`
    Password string `json:"password" binding:"required,min=6" example:"password123"`
    Name     string `json:"name" binding:"required" example:"John Doe"`
    Phone    string `json:"phone" example:"+1234567890"`
    Country  string `json:"country" example:"US"`
}

// SignInRequest represents user login payload
type SignInRequest struct {
    Email    string `json:"email" binding:"required,email" example:"user@example.com"`
    Password string `json:"password" binding:"required" example:"password123"`
}
```

**File**: `internal/dto/responses.go`
```go
package dto

// UserResponse represents user data in API responses
type UserResponse struct {
    ID        uint   `json:"id" example:"1"`
    UUID      string `json:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
    Email     string `json:"email" example:"user@example.com"`
    Name      string `json:"name" example:"John Doe"`
    Phone     string `json:"phone" example:"+1234567890"`
    CreatedAt string `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// AuthResponse represents authentication response (signup/signin)
type AuthResponse struct {
    Message string       `json:"message" example:"Login successful"`
    User    UserResponse `json:"user"`
    Token   string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// HealthResponse represents health check response
type HealthResponse struct {
    Status  int    `json:"status" example:"0"`
    Message string `json:"message" example:"Server is Healthy"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
    Error string `json:"error" example:"Invalid input"`
}
```

**Benefits**:
- Decouples API contracts from domain models
- Frontend changes don't affect database schema
- Clear API documentation
- Easy to version APIs

---

#### **Step 3: Create Repository Layer**

**File**: `internal/repository/user_repository.go`
```go
package repository

import (
    "go-booking-system/internal/domain"
    "gorm.io/gorm"
)

// UserRepository defines data access methods for User
type UserRepository interface {
    Create(user *domain.User) error
    FindByEmail(email string) (*domain.User, error)
    FindByID(id uint) (*domain.User, error)
    FindByUUID(uuid string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id uint) error
}

// userRepository implements UserRepository
type userRepository struct {
    db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// Create inserts a new user into database
func (r *userRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

// FindByEmail retrieves user by email
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// FindByID retrieves user by ID
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
    var user domain.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// FindByUUID retrieves user by UUID
func (r *userRepository) FindByUUID(uuid string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("uuid = ?", uuid).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Update updates user data
func (r *userRepository) Update(user *domain.User) error {
    return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *userRepository) Delete(id uint) error {
    return r.db.Delete(&domain.User{}, id).Error
}
```

**File**: `internal/repository/country_repository.go`
```go
package repository

import (
    "go-booking-system/internal/domain"
    "gorm.io/gorm"
)

// CountryRepository defines data access methods for Country
type CountryRepository interface {
    FindByShortname(shortname string) (*domain.Country, error)
    FindAll() ([]domain.Country, error)
}

type countryRepository struct {
    db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) CountryRepository {
    return &countryRepository{db: db}
}

func (r *countryRepository) FindByShortname(shortname string) (*domain.Country, error) {
    var country domain.Country
    err := r.db.Where("shortname = ?", shortname).First(&country).Error
    if err != nil {
        return nil, err
    }
    return &country, nil
}

func (r *countryRepository) FindAll() ([]domain.Country, error) {
    var countries []domain.Country
    err := r.db.Find(&countries).Error
    return countries, err
}
```

**Benefits**:
- Database logic isolated in one place
- Easy to mock for testing
- Can swap databases without changing business logic
- Follows Repository pattern (industry standard)

---

### Phase 2: Business Logic (~4 hours)

#### **Step 4: Create Service Layer**

**File**: `internal/service/auth_service.go`
```go
package service

import (
    "errors"
    "go-booking-system/internal/domain"
    "go-booking-system/internal/dto"
    "go-booking-system/internal/repository"
    "go-booking-system/pkg/logger"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "gorm.io/gorm"
)

// AuthService defines authentication business logic
type AuthService interface {
    SignUp(req dto.SignUpRequest) (*dto.AuthResponse, error)
    SignIn(req dto.SignInRequest) (*dto.AuthResponse, error)
}

type authService struct {
    userRepo    repository.UserRepository
    countryRepo repository.CountryRepository
    logger      logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
    userRepo repository.UserRepository,
    countryRepo repository.CountryRepository,
    logger logger.Logger,
) AuthService {
    return &authService{
        userRepo:    userRepo,
        countryRepo: countryRepo,
        logger:      logger,
    }
}

// SignUp registers a new user
func (s *authService) SignUp(req dto.SignUpRequest) (*dto.AuthResponse, error) {
    // Check if user already exists
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err == nil && existingUser != nil {
        s.logger.Info("signup attempt with existing email", "email", req.Email)
        return nil, errors.New("email already registered")
    }
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        s.logger.Error("error checking existing user", err)
        return nil, errors.New("internal server error")
    }

    // Get country ID if provided
    var mobileCountryID *uint
    if req.Country != "" {
        country, err := s.countryRepo.FindByShortname(req.Country)
        if err != nil {
            s.logger.Warn("country not found", "country", req.Country)
            // Continue without country - don't fail signup
        } else {
            mobileCountryID = &country.ID
        }
    }

    // Create user domain object
    user := &domain.User{
        Email:           req.Email,
        Name:            req.Name,
        Phone:           req.Phone,
        MobileCountryId: mobileCountryID,
    }

    // Hash password
    if err := user.HashPassword(req.Password); err != nil {
        s.logger.Error("failed to hash password", err)
        return nil, errors.New("failed to process password")
    }

    // Save to database
    if err := s.userRepo.Create(user); err != nil {
        s.logger.Error("failed to create user", err)
        return nil, errors.New("failed to create user")
    }

    // Generate JWT token
    token, err := s.generateToken(user.ID)
    if err != nil {
        s.logger.Error("failed to generate token", err)
        return nil, errors.New("failed to generate token")
    }

    s.logger.Info("user registered successfully", "email", user.Email, "uuid", user.UUID)

    return &dto.AuthResponse{
        Message: "User registered successfully",
        User: dto.UserResponse{
            ID:        user.ID,
            UUID:      user.UUID,
            Email:     user.Email,
            Name:      user.Name,
            Phone:     user.Phone,
            CreatedAt: user.CreatedAt.Format(time.RFC3339),
        },
        Token: token,
    }, nil
}

// SignIn authenticates a user
func (s *authService) SignIn(req dto.SignInRequest) (*dto.AuthResponse, error) {
    // Find user by email
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            s.logger.Info("signin attempt with non-existent email", "email", req.Email)
            return nil, errors.New("invalid credentials")
        }
        s.logger.Error("error finding user", err)
        return nil, errors.New("internal server error")
    }

    // Check password
    if err := user.CheckPassword(req.Password); err != nil {
        s.logger.Info("signin attempt with wrong password", "email", req.Email)
        return nil, errors.New("invalid credentials")
    }

    // Generate token
    token, err := s.generateToken(user.ID)
    if err != nil {
        s.logger.Error("failed to generate token", err)
        return nil, errors.New("failed to generate token")
    }

    s.logger.Info("user signed in successfully", "email", user.Email)

    return &dto.AuthResponse{
        Message: "Login successful",
        User: dto.UserResponse{
            ID:        user.ID,
            UUID:      user.UUID,
            Email:     user.Email,
            Name:      user.Name,
            Phone:     user.Phone,
            CreatedAt: user.CreatedAt.Format(time.RFC3339),
        },
        Token: token,
    }, nil
}

// generateToken creates a JWT token
func (s *authService) generateToken(userID uint) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
```

**Benefits**:
- Business logic separated from HTTP transport
- Easy to test (mock repositories)
- Can reuse across different interfaces (REST, GraphQL, gRPC)
- Proper error logging and handling
- Single Responsibility Principle

---

#### **Step 5: Add Structured Logging**

**Install dependency**:
```bash
go get github.com/sirupsen/logrus
```

**File**: `pkg/logger/logger.go`
```go
package logger

import (
    "os"

    "github.com/sirupsen/logrus"
)

// Logger defines logging interface
type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, err error, fields ...interface{})
    Warn(msg string, fields ...interface{})
    Debug(msg string, fields ...interface{})
}

type logrusLogger struct {
    logger *logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger() Logger {
    log := logrus.New()
    log.SetOutput(os.Stdout)
    log.SetFormatter(&logrus.JSONFormatter{})

    // Set log level based on environment
    env := os.Getenv("APP_ENV")
    if env == "production" {
        log.SetLevel(logrus.InfoLevel)
    } else {
        log.SetLevel(logrus.DebugLevel)
    }

    return &logrusLogger{logger: log}
}

func (l *logrusLogger) Info(msg string, fields ...interface{}) {
    l.logger.WithFields(parseFields(fields...)).Info(msg)
}

func (l *logrusLogger) Error(msg string, err error, fields ...interface{}) {
    f := parseFields(fields...)
    if err != nil {
        f["error"] = err.Error()
    }
    l.logger.WithFields(f).Error(msg)
}

func (l *logrusLogger) Warn(msg string, fields ...interface{}) {
    l.logger.WithFields(parseFields(fields...)).Warn(msg)
}

func (l *logrusLogger) Debug(msg string, fields ...interface{}) {
    l.logger.WithFields(parseFields(fields...)).Debug(msg)
}

// parseFields converts variadic args to logrus.Fields
func parseFields(fields ...interface{}) logrus.Fields {
    f := logrus.Fields{}
    for i := 0; i < len(fields)-1; i += 2 {
        key, ok := fields[i].(string)
        if ok {
            f[key] = fields[i+1]
        }
    }
    return f
}
```

**Benefits**:
- Structured logs (searchable, parseable)
- Easy to integrate with log aggregators (ELK, Datadog, CloudWatch)
- Different log levels for different environments
- Debug production issues faster

---

### Phase 3: Error Handling & Responses (~2 hours)

#### **Step 6: Centralized Error Handling**

**File**: `pkg/errors/errors.go`
```go
package errors

import "net/http"

// AppError represents application-specific error
type AppError struct {
    Code    string
    Message string
    Status  int
}

func (e *AppError) Error() string {
    return e.Message
}

// Predefined errors
var (
    ErrEmailExists        = &AppError{"EMAIL_EXISTS", "Email already registered", http.StatusConflict}
    ErrInvalidCredentials = &AppError{"INVALID_CREDENTIALS", "Invalid credentials", http.StatusUnauthorized}
    ErrInvalidInput       = &AppError{"INVALID_INPUT", "Invalid input data", http.StatusBadRequest}
    ErrInternal           = &AppError{"INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError}
    ErrUnauthorized       = &AppError{"UNAUTHORIZED", "Unauthorized access", http.StatusUnauthorized}
    ErrNotFound           = &AppError{"NOT_FOUND", "Resource not found", http.StatusNotFound}
    ErrInvalidToken       = &AppError{"INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized}
)

// New creates a new AppError
func New(code, message string, status int) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Status:  status,
    }
}
```

**File**: `internal/middleware/error_handler.go`
```go
package middleware

import (
    "go-booking-system/pkg/errors"
    "go-booking-system/pkg/response"

    "github.com/gin-gonic/gin"
)

// ErrorHandler middleware catches errors and formats responses
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // Check if there are any errors
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err

            // Check if it's an AppError
            if appErr, ok := err.(*errors.AppError); ok {
                c.JSON(appErr.Status, response.Error(appErr.Code, appErr.Message))
                return
            }

            // Default to internal server error
            c.JSON(500, response.Error("INTERNAL_ERROR", "An unexpected error occurred"))
        }
    }
}
```

**Benefits**:
- Consistent error format across all endpoints
- Easy to add new error types
- Centralized error handling logic
- Cleaner handler code

---

#### **Step 7: API Response Wrapper**

**File**: `pkg/response/response.go`
```go
package response

// Response represents standard API response
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorData  `json:"error,omitempty"`
}

// ErrorData represents error details
type ErrorData struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// Success creates a success response
func Success(data interface{}) Response {
    return Response{
        Success: true,
        Data:    data,
    }
}

// Error creates an error response
func Error(code, message string) Response {
    return Response{
        Success: false,
        Error: &ErrorData{
            Code:    code,
            Message: message,
        },
    }
}
```

**Example Responses**:
```json
// Success
{
    "success": true,
    "data": {
        "message": "User registered successfully",
        "user": {...},
        "token": "..."
    }
}

// Error
{
    "success": false,
    "error": {
        "code": "EMAIL_EXISTS",
        "message": "Email already registered"
    }
}
```

**Benefits**:
- Consistent response format
- Easy for frontend to handle
- Clear success/failure indication

---

#### **Step 8: Refactor Handlers (Controllers)**

**File**: `internal/handler/auth_handler.go`
```go
package handler

import (
    "go-booking-system/internal/dto"
    "go-booking-system/internal/service"
    "go-booking-system/pkg/errors"
    "go-booking-system/pkg/logger"
    "go-booking-system/pkg/response"
    "net/http"

    "github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
    authService service.AuthService
    logger      logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, logger logger.Logger) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        logger:      logger,
    }
}

// SignUp godoc
// @Summary Register a new user
// @Description Create a new user account with email, password, name, phone, and country
// @Tags Account
// @Accept json
// @Produce json
// @Param input body dto.SignUpRequest true "User registration data"
// @Success 201 {object} response.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure 400 {object} response.Response{error=response.ErrorData} "Invalid input"
// @Failure 409 {object} response.Response{error=response.ErrorData} "Email already registered"
// @Router /account/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
    var req dto.SignUpRequest

    // Validate input
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Error("INVALID_INPUT", err.Error()))
        return
    }

    // Call service
    result, err := h.authService.SignUp(req)
    if err != nil {
        // Check for specific errors
        if err.Error() == "email already registered" {
            c.JSON(http.StatusConflict, response.Error("EMAIL_EXISTS", err.Error()))
            return
        }
        c.JSON(http.StatusInternalServerError, response.Error("INTERNAL_ERROR", err.Error()))
        return
    }

    c.JSON(http.StatusCreated, response.Success(result))
}

// SignIn godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Account
// @Accept json
// @Produce json
// @Param input body dto.SignInRequest true "Login credentials"
// @Success 200 {object} response.Response{data=dto.AuthResponse} "Login successful"
// @Failure 400 {object} response.Response{error=response.ErrorData} "Invalid input"
// @Failure 401 {object} response.Response{error=response.ErrorData} "Invalid credentials"
// @Router /account/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
    var req dto.SignInRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Error("INVALID_INPUT", err.Error()))
        return
    }

    result, err := h.authService.SignIn(req)
    if err != nil {
        if err.Error() == "invalid credentials" {
            c.JSON(http.StatusUnauthorized, response.Error("INVALID_CREDENTIALS", err.Error()))
            return
        }
        c.JSON(http.StatusInternalServerError, response.Error("INTERNAL_ERROR", err.Error()))
        return
    }

    c.JSON(http.StatusOK, response.Success(result))
}
```

**File**: `internal/handler/health_handler.go`
```go
package handler

import (
    "go-booking-system/internal/dto"
    "go-booking-system/pkg/response"
    "net/http"

    "github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

// HealthStatus godoc
// @Summary Health check endpoint
// @Description Check if the server is running and healthy. Status 0 means healthy.
// @Tags Health
// @Produce json
// @Success 200 {object} response.Response{data=dto.HealthResponse} "Server is healthy"
// @Router /health/ [get]
func (h *HealthHandler) HealthStatus(c *gin.Context) {
    c.JSON(http.StatusOK, response.Success(dto.HealthResponse{
        Status:  0,
        Message: "Server is Healthy",
    }))
}
```

**Benefits**:
- Thin handlers (just HTTP transport)
- Business logic in service layer
- Clean, readable code
- Easy to test

---

### Phase 4: Advanced Features (~3 hours)

#### **Step 9: Dependency Injection**

**File**: `cmd/api/main.go`
```go
package main

import (
    "log"
    "os"

    "go-booking-system/config"
    "go-booking-system/internal/domain"
    "go-booking-system/internal/handler"
    "go-booking-system/internal/middleware"
    "go-booking-system/internal/repository"
    "go-booking-system/internal/routes"
    "go-booking-system/internal/service"
    "go-booking-system/pkg/logger"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"

    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "go-booking-system/docs"
)

// @title Go Booking System API
// @version 1.0
// @description Production-ready booking system API with clean architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Initialize application
    app := initializeApp()

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s...", port)
    log.Printf("Swagger docs available at http://localhost:%s/swagger/index.html", port)

    if err := app.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

// initializeApp sets up all dependencies
func initializeApp() *gin.Engine {
    // Initialize logger
    logger := logger.NewLogger()

    // Connect to database
    config.ConnectDatabase()
    db := config.DB

    // Auto migrate
    if err := db.AutoMigrate(&domain.User{}, &domain.Country{}); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    countryRepo := repository.NewCountryRepository(db)

    // Initialize services
    authService := service.NewAuthService(userRepo, countryRepo, logger)

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService, logger)
    healthHandler := handler.NewHealthHandler()

    // Setup router
    router := gin.Default()

    // Apply middleware
    router.Use(middleware.ErrorHandler())

    // Setup routes
    routes.SetupRoutes(router, authHandler, healthHandler)

    // Swagger
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return router
}
```

**File**: `internal/routes/routes.go`
```go
package routes

import (
    "go-booking-system/internal/handler"
    "go-booking-system/internal/middleware"

    "github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
    router *gin.Engine,
    authHandler *handler.AuthHandler,
    healthHandler *handler.HealthHandler,
) {
    // API group
    api := router.Group("/api")
    {
        // Health check
        health := api.Group("/health")
        {
            health.GET("/", healthHandler.HealthStatus)
        }

        // Public account endpoints
        account := api.Group("/account")
        {
            account.POST("/signup", authHandler.SignUp)
            account.POST("/signin", authHandler.SignIn)
        }

        // Protected routes (example)
        protected := api.Group("/")
        protected.Use(middleware.RequireAuth())
        {
            // protected.GET("/profile", profileHandler.GetProfile)
            // protected.PUT("/profile", profileHandler.UpdateProfile)
        }
    }
}
```

**Benefits**:
- Clear dependency flow
- Easy to test (inject mocks)
- Loose coupling
- Single place to see all dependencies

---

#### **Step 10: Authentication Middleware**

**File**: `internal/middleware/auth.go`
```go
package middleware

import (
    "net/http"
    "os"
    "strings"

    "go-booking-system/pkg/response"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// RequireAuth validates JWT token
func RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get token from header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(
                http.StatusUnauthorized,
                response.Error("UNAUTHORIZED", "Missing authorization header"),
            )
            return
        }

        // Parse Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(
                http.StatusUnauthorized,
                response.Error("INVALID_TOKEN", "Invalid authorization format"),
            )
            return
        }

        tokenString := parts[1]

        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Validate signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(
                http.StatusUnauthorized,
                response.Error("INVALID_TOKEN", "Invalid or expired token"),
            )
            return
        }

        // Extract claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            // Set user ID in context
            if userID, ok := claims["user_id"].(float64); ok {
                c.Set("userID", uint(userID))
            }
        }

        c.Next()
    }
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
    userID, exists := c.Get("userID")
    if !exists {
        return 0, false
    }

    id, ok := userID.(uint)
    return id, ok
}
```

**Usage in handlers**:
```go
func (h *ProfileHandler) GetProfile(c *gin.Context) {
    userID, exists := middleware.GetUserID(c)
    if !exists {
        c.JSON(401, response.Error("UNAUTHORIZED", "User not authenticated"))
        return
    }

    // Use userID to fetch profile
    // ...
}
```

**Benefits**:
- Reusable across protected routes
- Centralized token validation
- Easy to add more auth methods (API keys, OAuth)

---

### Phase 5: Testing & Documentation (~4 hours)

#### **Step 11: Unit Tests**

**Install testing dependencies**:
```bash
go get github.com/stretchr/testify
```

**File**: `internal/service/auth_service_test.go`
```go
package service

import (
    "errors"
    "testing"

    "go-booking-system/internal/domain"
    "go-booking-system/internal/dto"
    "go-booking-system/pkg/logger"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"
)

// MockUserRepository is a mock for UserRepository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByUUID(uuid string) (*domain.User, error) {
    args := m.Called(uuid)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}

// MockCountryRepository is a mock for CountryRepository
type MockCountryRepository struct {
    mock.Mock
}

func (m *MockCountryRepository) FindByShortname(shortname string) (*domain.Country, error) {
    args := m.Called(shortname)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Country), args.Error(1)
}

func (m *MockCountryRepository) FindAll() ([]domain.Country, error) {
    args := m.Called()
    return args.Get(0).([]domain.Country), args.Error(1)
}

// MockLogger is a mock for Logger
type MockLogger struct{}

func (m *MockLogger) Info(msg string, fields ...interface{})           {}
func (m *MockLogger) Error(msg string, err error, fields ...interface{}) {}
func (m *MockLogger) Warn(msg string, fields ...interface{})           {}
func (m *MockLogger) Debug(msg string, fields ...interface{})          {}

func TestAuthService_SignUp_Success(t *testing.T) {
    // Arrange
    mockUserRepo := new(MockUserRepository)
    mockCountryRepo := new(MockCountryRepository)
    mockLogger := &MockLogger{}

    service := NewAuthService(mockUserRepo, mockCountryRepo, mockLogger)

    req := dto.SignUpRequest{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
        Phone:    "+1234567890",
    }

    // Mock: User doesn't exist
    mockUserRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

    // Mock: Create user succeeds
    mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

    // Act
    result, err := service.SignUp(req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "User registered successfully", result.Message)
    assert.Equal(t, req.Email, result.User.Email)
    assert.Equal(t, req.Name, result.User.Name)
    assert.NotEmpty(t, result.Token)

    mockUserRepo.AssertExpectations(t)
}

func TestAuthService_SignUp_EmailExists(t *testing.T) {
    // Arrange
    mockUserRepo := new(MockUserRepository)
    mockCountryRepo := new(MockCountryRepository)
    mockLogger := &MockLogger{}

    service := NewAuthService(mockUserRepo, mockCountryRepo, mockLogger)

    req := dto.SignUpRequest{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    }

    existingUser := &domain.User{
        ID:    1,
        Email: req.Email,
    }

    // Mock: User already exists
    mockUserRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

    // Act
    result, err := service.SignUp(req)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Equal(t, "email already registered", err.Error())

    mockUserRepo.AssertExpectations(t)
}

func TestAuthService_SignIn_Success(t *testing.T) {
    // Arrange
    mockUserRepo := new(MockUserRepository)
    mockCountryRepo := new(MockCountryRepository)
    mockLogger := &MockLogger{}

    service := NewAuthService(mockUserRepo, mockCountryRepo, mockLogger)

    password := "password123"
    user := &domain.User{
        ID:    1,
        Email: "test@example.com",
        Name:  "Test User",
    }
    user.HashPassword(password)

    req := dto.SignInRequest{
        Email:    user.Email,
        Password: password,
    }

    // Mock: Find user
    mockUserRepo.On("FindByEmail", req.Email).Return(user, nil)

    // Act
    result, err := service.SignIn(req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Login successful", result.Message)
    assert.Equal(t, user.Email, result.User.Email)
    assert.NotEmpty(t, result.Token)

    mockUserRepo.AssertExpectations(t)
}

func TestAuthService_SignIn_InvalidCredentials(t *testing.T) {
    // Arrange
    mockUserRepo := new(MockUserRepository)
    mockCountryRepo := new(MockCountryRepository)
    mockLogger := &MockLogger{}

    service := NewAuthService(mockUserRepo, mockCountryRepo, mockLogger)

    req := dto.SignInRequest{
        Email:    "test@example.com",
        Password: "wrongpassword",
    }

    // Mock: User not found
    mockUserRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

    // Act
    result, err := service.SignIn(req)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Equal(t, "invalid credentials", err.Error())

    mockUserRepo.AssertExpectations(t)
}
```

**Run tests**:
```bash
go test ./internal/service/... -v
```

**Benefits**:
- Catch bugs before production
- Safe refactoring
- Living documentation
- Confidence in code changes

---

#### **Step 12: Swagger Documentation**

Already covered in handlers above. Just run:
```bash
swag init -g cmd/api/main.go
```

---

#### **Step 13: Database Config Improvements**

**File**: `config/database.go`
```go
package config

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase establishes database connection with proper configuration
func ConnectDatabase() {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    // Configure GORM
    config := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    }

    database, err := gorm.Open(postgres.Open(dsn), config)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Get underlying SQL database
    sqlDB, err := database.DB()
    if err != nil {
        log.Fatal("Failed to get database instance:", err)
    }

    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)                   // Max idle connections
    sqlDB.SetMaxOpenConns(100)                  // Max open connections
    sqlDB.SetConnMaxLifetime(time.Hour)         // Connection lifetime
    sqlDB.SetConnMaxIdleTime(10 * time.Minute)  // Idle connection timeout

    DB = database
    log.Println("Database connected successfully!")
}
```

**Benefits**:
- Connection pooling prevents connection exhaustion
- Timeouts prevent hanging queries
- Proper logging for debugging

---

## 📅 Implementation Timeline

### Week 1: Foundation & Business Logic
- **Day 1**: Steps 1-2 (Project structure, DTO layer)
- **Day 2**: Step 3 (Repository layer)
- **Day 3**: Step 4 (Service layer)
- **Day 4**: Steps 5-6 (Logging, Error handling)
- **Day 5**: Step 7 (Refactor handlers)

### Week 2: Advanced Features & Testing
- **Day 1**: Steps 8-9 (DI, Response wrapper)
- **Day 2**: Step 10 (Auth middleware)
- **Day 3-4**: Step 11 (Unit tests)
- **Day 5**: Steps 12-13 (Swagger, DB improvements)

---

## 🎯 Success Metrics

After completing all steps, your project will have:

1. ✅ **Scalability Score: 7-8/10** (up from 3.4/10)
2. ✅ **Clean Architecture** with proper separation of concerns
3. ✅ **80%+ Test Coverage** for business logic
4. ✅ **Proper Error Handling** with meaningful messages
5. ✅ **Structured Logging** for production debugging
6. ✅ **API Documentation** with Swagger
7. ✅ **Dependency Injection** for testability
8. ✅ **Production-Ready** codebase

---

## 🚀 Next Steps

1. Start with **Step 1** (Create project structure)
2. Work through steps sequentially
3. Test after each step
4. Commit after each completed step
5. Ask Claude for help when stuck!

---

## 📚 Learning Resources

### Go Clean Architecture
- https://github.com/bxcodec/go-clean-arch
- https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

### Testing in Go
- https://quii.gitbook.io/learn-go-with-tests/

### SOLID Principles
- https://dave.cheney.net/2016/08/20/solid-go-design

### Gin Best Practices
- https://github.com/gin-gonic/gin#documentation

---

**Good luck with your refactoring! You're on the path to becoming a better developer! 🚀**
