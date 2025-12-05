package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthStatus godoc
// @Summary Health check endpoint
// @Description Check if the server is running and healthy
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthResponse "Server is Healthy"
// @Router /api/health [get]
// HealthStatus
func HealthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "Server is Healthy",
	})
}
