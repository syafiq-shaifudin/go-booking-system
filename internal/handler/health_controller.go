package handler

import (
	"go-booking-system/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check HTTP requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthStatus godoc
// @Summary Health check endpoint
// @Description Check if the server is running and healthy. Status 0 means healthy.
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthResponse "Server is healthy"
// @Router /api/health/ [get]
func (h *HealthHandler) HealthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:  0,
		Message: "Server is Healthy",
	})
}
