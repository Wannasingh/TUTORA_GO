package domain

import "context"

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"` // "student" or "tutor"
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type OAuthLoginRequest struct {
	Token string `json:"token" binding:"required"` // ID Token from Google or Apple Client
	Name  string `json:"name,omitempty"`           // Optional name (provided on first login/signup)
	Role  string `json:"role" binding:"required"`  // "student" or "tutor"
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type AuthUsecase interface {
	RegisterWithEmail(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	LoginWithEmail(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	LoginWithGoogle(ctx context.Context, req *OAuthLoginRequest) (*AuthResponse, error)
	LoginWithApple(ctx context.Context, req *OAuthLoginRequest) (*AuthResponse, error)
}
