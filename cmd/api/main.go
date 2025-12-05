package main

import (
	"go-booking-system/config"
	"go-booking-system/internal/domain"
	"go-booking-system/internal/routes"
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

	// Initialize Gin
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

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
