package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/haru/bytestutor/backend/domain"
)

// MockUserRepository implements domain.UserRepository interface
type MockUserRepository struct {
	users         map[int]*domain.User
	usersByEmail  map[string]*domain.User
	createFunc    func(user *domain.User) error
	getByIDFunc   func(id int) (*domain.User, error)
	getByEmailFunc func(email string) (*domain.User, error)
	deleteFunc    func(id int) error
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(user)
	}
	user.ID = len(m.users) + 1
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return m.users[id], nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(email)
	}
	return m.usersByEmail[email], nil
}

func (m *MockUserRepository) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetByAppleID(ctx context.Context, appleID string) (*domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	delete(m.users, id)
	return nil
}

func TestRegisterUser_Success(t *testing.T) {
	repo := &MockUserRepository{
		users:        make(map[int]*domain.User),
		usersByEmail: make(map[string]*domain.User),
	}
	u := NewUserUsecase(repo)

	user := &domain.User{
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "student",
	}

	err := u.Register(context.Background(), user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID to be 1, got %d", user.ID)
	}
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	existingUser := &domain.User{
		ID:    1,
		Name:  "Existing User",
		Email: "test@example.com",
	}
	repo := &MockUserRepository{
		users: map[int]*domain.User{
			1: existingUser,
		},
		usersByEmail: map[string]*domain.User{
			"test@example.com": existingUser,
		},
	}
	u := NewUserUsecase(repo)

	user := &domain.User{
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "student",
	}

	err := u.Register(context.Background(), user)
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}

	if err.Error() != "email is already registered" {
		t.Errorf("expected 'email is already registered' error, got '%v'", err.Error())
	}
}

func TestGetProfile_Success(t *testing.T) {
	existingUser := &domain.User{
		ID:    42,
		Name:  "Profile User",
		Email: "profile@example.com",
	}
	repo := &MockUserRepository{
		users: map[int]*domain.User{
			42: existingUser,
		},
	}
	u := NewUserUsecase(repo)

	profile, err := u.GetProfile(context.Background(), 42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if profile == nil || profile.Name != "Profile User" {
		t.Errorf("expected user profile name to be 'Profile User', got %v", profile)
	}
}

func TestDeleteAccount_Success(t *testing.T) {
	repo := &MockUserRepository{
		users: map[int]*domain.User{
			10: {ID: 10, Name: "Delete Me", Email: "delete@example.com"},
		},
	}
	u := NewUserUsecase(repo)

	err := u.DeleteAccount(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	profile, _ := u.GetProfile(context.Background(), 10)
	if profile != nil {
		t.Error("expected profile to be nil after deletion")
	}
}

func TestRegisterUser_RepoError(t *testing.T) {
	repo := &MockUserRepository{
		getByEmailFunc: func(email string) (*domain.User, error) {
			return nil, errors.New("db error")
		},
	}
	u := NewUserUsecase(repo)

	user := &domain.User{
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "student",
	}

	err := u.Register(context.Background(), user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
