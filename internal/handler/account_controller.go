package handler

import (
	"go-booking-system/config"
	"go-booking-system/internal/domain"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SignUpInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Country  string `json:"country"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SignUp godoc
// @Summary Sign Up
// @Description Register a new user account
// @Tags Account
// @Accept json
// @Produce json
// @Param input body dto.req_signUp_successfulRegistration true "Sign up credentials"
// @Success 201 {object} dto.res_signUp_successfulRegistration "User created successfully"
// @Success 409  {object} dto.res_signUp_existingEmail "Email Already Exist"
// @Success 500  {object} dto.res_signUp_failed "Failed to hash password / Failed to create user / Failed to generate token"
// @Router /api/account/signup [post]
// Sign Up Handler
func SignUp(c *gin.Context) {
	var input SignUpInput

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser domain.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Save to database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.UUID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": token,
	})
}

// Sign In Handler
func SignIn(c *gin.Context) {
	var input SignInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user domain.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": token,
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
