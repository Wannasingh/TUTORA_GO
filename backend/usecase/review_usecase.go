package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type reviewUsecase struct {
	repo domain.ReviewRepository
}

func NewReviewUsecase(repo domain.ReviewRepository) domain.ReviewUsecase {
	return &reviewUsecase{repo: repo}
}

func (u *reviewUsecase) SubmitReview(ctx context.Context, review *domain.TutorReview) error {
	if review.Rating < 1 || review.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	if err := u.repo.CreateOrUpdate(ctx, review); err != nil {
		return err
	}

	// Recalculate tutor average rating
	avg, err := u.repo.GetAverageRating(ctx, review.TutorID)
	if err != nil {
		return err
	}
	avg = math.Round(avg*100) / 100
	return u.repo.UpdateTutorRating(ctx, review.TutorID, avg)
}

func (u *reviewUsecase) GetTutorReviews(ctx context.Context, tutorID int) ([]*domain.TutorReview, error) {
	return u.repo.GetByTutorID(ctx, tutorID)
}

func (u *reviewUsecase) DeleteReview(ctx context.Context, userID, reviewID int) error {
	review, err := u.repo.GetByID(ctx, reviewID)
	if err != nil {
		return err
	}
	if review == nil {
		return fmt.Errorf("review not found")
	}
	if review.ReviewerID != userID {
		return fmt.Errorf("not authorized to delete this review")
	}
	tutorID := review.TutorID
	if err := u.repo.Delete(ctx, reviewID); err != nil {
		return err
	}
	// Recalculate rating after deletion
	avg, err := u.repo.GetAverageRating(ctx, tutorID)
	if err != nil {
		return err
	}
	avg = math.Round(avg*100) / 100
	return u.repo.UpdateTutorRating(ctx, tutorID, avg)
}
