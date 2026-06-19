package domain

import "context"

type TutorReview struct {
	ID             int     `json:"id"`
	ReviewerID     int     `json:"reviewer_id"`
	TutorID        int     `json:"tutor_id"`
	Rating         float64 `json:"rating"`
	Body           *string `json:"body,omitempty"`
	CreatedAt      string  `json:"created_at"`
	ReviewerName   string  `json:"reviewer_name"`
	ReviewerAvatar *string `json:"reviewer_avatar,omitempty"`
}

type ReviewRepository interface {
	CreateOrUpdate(ctx context.Context, review *TutorReview) error
	GetByTutorID(ctx context.Context, tutorID int) ([]*TutorReview, error)
	GetByID(ctx context.Context, id int) (*TutorReview, error)
	Delete(ctx context.Context, id int) error
	GetAverageRating(ctx context.Context, tutorID int) (float64, error)
	UpdateTutorRating(ctx context.Context, tutorID int, rating float64) error
}

type ReviewUsecase interface {
	SubmitReview(ctx context.Context, review *TutorReview) error
	GetTutorReviews(ctx context.Context, tutorID int) ([]*TutorReview, error)
	DeleteReview(ctx context.Context, userID, reviewID int) error
}
