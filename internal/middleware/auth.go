package middleware

import (
	"go-booking-system/internal/dto"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth validates JWT tokens and protects routes
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Get Authorization header
		// Expected format: "Authorization: Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Authorization header required",
			})
			c.Abort() // Stop processing, don't call next handler
			return
		}

		// Step 2: Extract token string (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// TrimPrefix didn't change anything = "Bearer " wasn't there
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Invalid authorization format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Step 3: Parse and verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// Step 4: Check if token is valid
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Step 5: Extract claims (payload data)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Get the user UUID from the token
			// This matches the "uuid" claim from account_service.go:137
			if userUUID, exists := claims["uuid"]; exists {
				// Store UUID in context so handlers can access it
				// Handlers can get this with: c.Get("userUUID")
				c.Set("userUUID", userUUID)
			}
		}

		// Step 6: Token is valid, proceed to the actual handler
		c.Next()
	}
}
