package handler

import (
	"go-booking-system/internal/dto"
	"go-booking-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	accountService service.AccountService
}

// NewAccountHandler creates a new account handler instance
func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// SignUp godoc
// @Summary Register a new user
// @Description Create a new user account with email, password, name, phone, and country
// @Tags Account
// @Accept json
// @Produce json
// @Param input body dto.SignUpRequest true "User registration data"
// @Success 201 {object} dto.SignUp_Success "User registered successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid input data"
// @Failure 409 {object} dto.ErrorResponse "Email already registered"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/account/signup [post]
func (h *AccountHandler) SignUp(c *gin.Context) {
	var input dto.SignUpRequest

	// Validate HTTP input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Call service layer for business logic
	result, err := h.accountService.SignUp(input)
	if err != nil {
		// Handle specific errors
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, result)
}

// SignIn godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Account
// @Accept json
// @Produce json
// @Param input body dto.SignInRequest true "Login credentials"
// @Success 200 {object} dto.SignUp_Success "Login successful"
// @Failure 400 {object} dto.ErrorResponse "Invalid input data"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/account/signin [post]
func (h *AccountHandler) SignIn(c *gin.Context) {
	var input dto.SignInRequest

	// Validate HTTP input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Call service layer for business logic
	result, err := h.accountService.SignIn(input)
	if err != nil {
		// Handle specific errors
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, result)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags Account
// @Security BearerAuth
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c)
// @Success 200 {object} dto.UserResponse "User profile"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/account/profile [get]
func (h *AccountHandler) GetProfile(c *gin.Context) {
	// Get the user UUID that the middleware stored in context
	// The RequireAuth middleware extracts this from the JWT token
	userUUID, exists := c.Get("userUUID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "User not found in context",
		})
		return
	}

	// Convert interface{} to string
	uuid, ok := userUUID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Invalid user UUID format",
		})
		return
	}

	// Call service layer for business logic
	result, err := h.accountService.GetProfile(uuid)
	if err != nil {
		// Handle specific errors
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, result)
}
