package routes

import (
	"go-booking-system/internal/handler"
	"go-booking-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes with handler dependencies
func SetupRoutes(
	router *gin.Engine,
	accountHandler *handler.AccountHandler,
	healthHandler *handler.HealthHandler,
) {
	// Health check routes
	health := router.Group("/api/health")
	{
		health.GET("/", healthHandler.HealthStatus)
	}

	// Account routes (public - no authentication required)
	account := router.Group("/api/account")
	{
		account.POST("/signup", accountHandler.SignUp)
		account.POST("/signin", accountHandler.SignIn)
	}

	// Protected routes (require JWT authentication)
	protected := router.Group("/api/account")
	protected.Use(middleware.RequireAuth()) // Apply JWT verification middleware
	{
		protected.GET("/profile", accountHandler.GetProfile)
	}
}
