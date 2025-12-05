package routes

import (
	"go-booking-system/internal/handler"

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

	// Account routes (public)
	account := router.Group("/api/account")
	{
		account.POST("/signup", accountHandler.SignUp)
		account.POST("/signin", accountHandler.SignIn)
	}

	// Protected routes (we'll add these later)
	// protected := router.Group("/api")
	// protected.Use(middleware.RequireAuth())
	// {
	//     protected.GET("/profile", profileHandler.GetProfile)
	// }
}
