package main

import (
	"go-booking-system/config"
	"go-booking-system/internal/domain"
	"go-booking-system/internal/handler"
	"go-booking-system/internal/repository"
	"go-booking-system/internal/routes"
	"go-booking-system/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "go-booking-system/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	config.ConnectDatabase()

	// Auto migrate database
	config.DB.AutoMigrate(&domain.User{}, &domain.Country{})

	// Initialize repositories
	userRepo := repository.NewUserRepository(config.DB)
	countryRepo := repository.NewCountryRepository(config.DB)

	// Initialize services
	accountService := service.NewAccountService(userRepo, countryRepo)

	// Initialize handlers
	accountHandler := handler.NewAccountHandler(accountService)
	healthHandler := handler.NewHealthHandler()

	// Initialize Gin router
	router := gin.Default()

	// Setup routes with handler dependencies
	routes.SetupRoutes(router, accountHandler, healthHandler)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	log.Printf("Swagger docs available at http://localhost:%s/swagger/index.html", port)
	router.Run(":" + port)
}
