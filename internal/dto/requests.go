package dto

// SignUpRequest represents user registration payload
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Phone    string `json:"phone" example:"234567890"`
	Country  string `json:"country" example:"US"`
}

// SignInRequest represents user login payload
type SignInRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}
