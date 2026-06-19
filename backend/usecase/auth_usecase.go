package usecase

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/Wannasingh/TUTORA_GO/backend/config"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
	"github.com/Wannasingh/TUTORA_GO/backend/utils"
)

type authUsecase struct {
	userRepo domain.UserRepository
	cfg      *config.Config
}

func NewAuthUsecase(repo domain.UserRepository, cfg *config.Config) domain.AuthUsecase {
	return &authUsecase{
		userRepo: repo,
		cfg:      cfg,
	}
}

func (a *authUsecase) RegisterWithEmail(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if email already exists
	existing, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email is already registered")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	passwordStr := string(hashedPassword)

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: &passwordStr,
	}

	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (a *authUsecase) LoginWithEmail(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.PasswordHash == nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (a *authUsecase) LoginWithGoogle(ctx context.Context, req *domain.OAuthLoginRequest) (*domain.AuthResponse, error) {
	// Verify Google ID Token
	googleID, email, name, err := utils.VerifyGoogleIDToken(req.Token, a.cfg.GoogleClientID)
	if err != nil {
		return nil, fmt.Errorf("google token verification failed: %v", err)
	}

	// Try to find user by Google ID
	user, err := a.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return nil, err
	}

	// If user does not exist, check if email exists to link account, or register new user
	if user == nil {
		user, err = a.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}

		if user != nil {
			// Link existing email to Google ID
			user.GoogleID = &googleID
		} else {
			// Register new OAuth user
			if req.Name != "" {
				name = req.Name
			}
			user = &domain.User{
				Name:     name,
				Email:    email,
				Role:     req.Role,
				GoogleID: &googleID,
			}
			if err := a.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		}
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (a *authUsecase) LoginWithApple(ctx context.Context, req *domain.OAuthLoginRequest) (*domain.AuthResponse, error) {
	// Verify Apple ID Token
	appleID, email, name, err := utils.VerifyAppleIdentityToken(req.Token, a.cfg.AppleBundleID)
	if err != nil {
		return nil, fmt.Errorf("apple token verification failed: %v", err)
	}

	// Try to find user by Apple ID
	user, err := a.userRepo.GetByAppleID(ctx, appleID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = a.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}

		if user != nil {
			user.AppleID = &appleID
		} else {
			if req.Name != "" {
				name = req.Name
			}
			user = &domain.User{
				Name:    name,
				Email:   email,
				Role:    req.Role,
				AppleID: &appleID,
			}
			if err := a.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		}
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}
