package routes

import (
	"go-booking-system/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	health := router.Group("/api/health")
	{
		health.GET("/", handler.HealthStatus)
	}

	account := router.Group("/api/account")
	{
		account.POST("/signup", handler.SignUp)
		account.POST("/signin", handler.SignIn)
	}

	// Protected routes (we'll add these later)
	// api := router.Group("/api")
	// api.Use(middleware.AuthMiddleware())
	// {
	//     api.GET("/profile", controllers.GetProfile)
	// }
}
