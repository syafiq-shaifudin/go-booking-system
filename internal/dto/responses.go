package dto

// UserResponse represents user data in API responses
type UserResponse struct {
	UUID      string `json:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string `json:"email" example:"user@example.com"`
	Name      string `json:"name" example:"John Doe"`
	Phone     string `json:"phone,omitempty" example:"+1234567890"`
	CreatedAt string `json:"created_at" example:"2024-12-05T08:00:00Z"`
}

type SignUp_Success struct {
	Message string       `json:"message" example:"Login successful"`
	User    UserResponse `json:"user"`
	Token   string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type SignIn_Success struct {
	Message string       `json:"message" example:"Login successful"`
	User    UserResponse `json:"user"`
	Token   string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  int    `json:"status" example:"0"`
	Message string `json:"message" example:"Server is Healthy"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error" example:"Error Message"`
}
