package usecase

import (
	"context"
	"errors"

	"github.com/haru/bytestutor/backend/domain"
)

type tutorUsecase struct {
	tutorRepo domain.TutorRepository
	userRepo  domain.UserRepository
}

func NewTutorUsecase(tRepo domain.TutorRepository, uRepo domain.UserRepository) domain.TutorUsecase {
	return &tutorUsecase{
		tutorRepo: tRepo,
		userRepo:  uRepo,
	}
}

func (u *tutorUsecase) BecomeTutor(ctx context.Context, tutor *domain.Tutor) error {
	user, err := u.userRepo.GetByID(ctx, tutor.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user does not exist")
	}

	return u.tutorRepo.Create(ctx, tutor)
}

func (u *tutorUsecase) GetTutorProfile(ctx context.Context, id int) (*domain.Tutor, error) {
	return u.tutorRepo.GetByID(ctx, id)
}

func (u *tutorUsecase) SearchTutors(ctx context.Context, subject string) ([]*domain.Tutor, error) {
	return u.tutorRepo.List(ctx, subject)
}
