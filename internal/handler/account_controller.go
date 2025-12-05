package handler

import (
	"go-booking-system/config"
	"go-booking-system/internal/domain"
	"go-booking-system/internal/dto"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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
func SignUp(c *gin.Context) {
	var input dto.SignUpRequest

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Check if user already exists
	var existingUser domain.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "Email already registered"})
		return
	}

	var mobile_country_id *uint
	if input.Country != "" {
		var country domain.Country
		if err := config.DB.Where("shortname = ?", input.Country).First(&country).Error; err == nil {
			mobile_country_id = &country.ID
		}
	}

	// Create user
	user := domain.User{
		Email:           input.Email,
		Name:            input.Name,
		Phone:           input.Phone,
		MobileCountryId: mobile_country_id,
	}

	// Hash password
	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to hash password"})
		return
	}

	// Save to database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	// Return success response using DTO
	c.JSON(http.StatusCreated, dto.SignUp_Success{
		Message: "User registered successfully",
		User: dto.UserResponse{
			UUID:      user.UUID,
			Email:     user.Email,
			Name:      user.Name,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
		Token: token,
	})
}

// SignIn godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Account
// @Accept json
// @Produce json
// @Param input body dto.SignInRequest true "Login credentials"
// @Success 200 {object} dto.SignIn_Success "Login successful"
// @Failure 400 {object} dto.ErrorResponse "Invalid input data"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/account/signin [post]
func SignIn(c *gin.Context) {
	var input dto.SignInRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Find user
	var user domain.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Check password
	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Generate token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	// Return success response using DTO
	c.JSON(http.StatusOK, dto.SignUp_Success{
		Message: "Login successful",
		User: dto.UserResponse{
			UUID:      user.UUID,
			Email:     user.Email,
			Name:      user.Name,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
		Token: token,
	})
}

// Generate JWT Token
func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		// "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"exp": time.Now().Add(time.Hour * 1).Unix(), // 1 hours only
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
