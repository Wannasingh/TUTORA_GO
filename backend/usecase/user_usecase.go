package usecase

import (
	"context"
	"errors"

	"github.com/haru/bytestutor/backend/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: repo}
}

func (u *userUsecase) Register(ctx context.Context, user *domain.User) error {
	existing, err := u.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email is already registered")
	}

	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetProfile(ctx context.Context, id int) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}
