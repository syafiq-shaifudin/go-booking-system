package service

import (
	"errors"
	"go-booking-system/internal/domain"
	"go-booking-system/internal/dto"
	"go-booking-system/internal/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AccountService defines account management business logic
type AccountService interface {
	SignUp(req dto.SignUpRequest) (*dto.SignUp_Success, error)
	SignIn(req dto.SignInRequest) (*dto.SignUp_Success, error)
}

// accountService implements AccountService
type accountService struct {
	userRepo    repository.UserRepository
	countryRepo repository.CountryRepository
}

// NewAccountService creates a new account service instance
func NewAccountService(
	userRepo repository.UserRepository,
	countryRepo repository.CountryRepository,
) AccountService {
	return &accountService{
		userRepo:    userRepo,
		countryRepo: countryRepo,
	}
}

// SignUp registers a new user
func (s *accountService) SignUp(req dto.SignUpRequest) (*dto.SignUp_Success, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check existing user")
	}

	// Get country ID if provided
	var mobileCountryID *uint
	if req.Country != "" {
		country, err := s.countryRepo.FindByShortname(req.Country)
		if err != nil {
			// Country not found - continue without it (non-critical)
			// In production, you might want to log this
		} else {
			mobileCountryID = &country.ID
		}
	}

	// Create user domain object
	user := &domain.User{
		Email:           req.Email,
		Name:            req.Name,
		Phone:           req.Phone,
		MobileCountryId: mobileCountryID,
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, errors.New("failed to process password")
	}

	// Save to database via repository
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Build response DTO
	return &dto.SignUp_Success{
		Message: "User registered successfully",
		User: dto.UserResponse{
			UUID:      user.UUID,
			Email:     user.Email,
			Name:      user.Name,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
		Token: token,
	}, nil
}

// SignIn authenticates a user
func (s *accountService) SignIn(req dto.SignInRequest) (*dto.SignUp_Success, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("failed to find user")
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Build response DTO
	return &dto.SignUp_Success{
		Message: "Login successful",
		User: dto.UserResponse{
			UUID:      user.UUID,
			Email:     user.Email,
			Name:      user.Name,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
		Token: token,
	}, nil
}

// generateToken creates a JWT token for the user
func (s *accountService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
