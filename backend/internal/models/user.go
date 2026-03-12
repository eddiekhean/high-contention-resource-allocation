package models

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User represents a user in the system
type User struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	PasswordHash  string   `json:"-"`
	WalletAddress string   `json:"wallet_address,omitempty"`
	Role          UserRole `json:"role"`
}

// Auth request/response models
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
