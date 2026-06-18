package usecase

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"github.com/haru/bytestutor/backend/domain"
)

func TestRegisterWithEmail_Success(t *testing.T) {
	repo := &MockUserRepository{
		users:        make(map[int]*domain.User),
		usersByEmail: make(map[string]*domain.User),
	}
	au := NewAuthUsecase(repo)

	req := &domain.RegisterRequest{
		Name:     "Auth User",
		Email:    "auth@example.com",
		Password: "password123",
		Role:     "tutor",
	}

	resp, err := au.RegisterWithEmail(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.User.Email != "auth@example.com" {
		t.Errorf("expected email to be auth@example.com, got %s", resp.User.Email)
	}

	if resp.Token == "" {
		t.Error("expected JWT token to be generated, got empty string")
	}
}

func TestLoginWithEmail_Success(t *testing.T) {
	// Hash password first
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secretpassword"), bcrypt.DefaultCost)
	passHashStr := string(hashedPassword)

	existingUser := &domain.User{
		ID:           1,
		Name:         "Login User",
		Email:        "login@example.com",
		Role:         "student",
		PasswordHash: &passHashStr,
	}

	repo := &MockUserRepository{
		usersByEmail: map[string]*domain.User{
			"login@example.com": existingUser,
		},
	}
	au := NewAuthUsecase(repo)

	req := &domain.LoginRequest{
		Email:    "login@example.com",
		Password: "secretpassword",
	}

	resp, err := au.LoginWithEmail(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("expected JWT token to be generated, got empty string")
	}
}

func TestLoginWithEmail_InvalidCredentials(t *testing.T) {
	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secretpassword"), bcrypt.DefaultCost)
	passHashStr := string(hashedPassword)

	existingUser := &domain.User{
		ID:           1,
		Name:         "Login User",
		Email:        "login@example.com",
		Role:         "student",
		PasswordHash: &passHashStr,
	}

	repo := &MockUserRepository{
		usersByEmail: map[string]*domain.User{
			"login@example.com": existingUser,
		},
	}
	au := NewAuthUsecase(repo)

	req := &domain.LoginRequest{
		Email:    "login@example.com",
		Password: "wrongpassword",
	}

	_, err := au.LoginWithEmail(context.Background(), req)
	if err == nil {
		t.Fatal("expected login to fail, got nil error")
	}

	if err.Error() != "invalid email or password" {
		t.Errorf("expected invalid credentials error, got '%v'", err.Error())
	}
}

func TestLoginWithGoogle_Mock(t *testing.T) {
	repo := &MockUserRepository{
		users:        make(map[int]*domain.User),
		usersByEmail: make(map[string]*domain.User),
	}
	au := NewAuthUsecase(repo)

	req := &domain.OAuthLoginRequest{
		Token: "mock_google_id123_test@google.com",
		Role:  "student",
	}

	resp, err := au.LoginWithGoogle(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.User.Email != "google-user@bytestutor.com" && resp.User.Email != "test@google.com" {
		t.Errorf("unexpected user email: %s", resp.User.Email)
	}

	if resp.Token == "" {
		t.Error("expected JWT token to be generated, got empty string")
	}
}

func TestLoginWithApple_Mock(t *testing.T) {
	repo := &MockUserRepository{
		users:        make(map[int]*domain.User),
		usersByEmail: make(map[string]*domain.User),
	}
	au := NewAuthUsecase(repo)

	req := &domain.OAuthLoginRequest{
		Token: "mock_apple_id555_test@apple.com",
		Role:  "student",
	}

	resp, err := au.LoginWithApple(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.User.Email != "apple-user@bytestutor.com" && resp.User.Email != "test@apple.com" {
		t.Errorf("unexpected user email: %s", resp.User.Email)
	}

	if resp.Token == "" {
		t.Error("expected JWT token to be generated, got empty string")
	}
}
