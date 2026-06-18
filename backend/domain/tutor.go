package domain

import "context"

type Tutor struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	User         *User   `json:"user,omitempty"`
	Subject      string  `json:"subject"`
	Bio          string  `json:"bio"`
	PricePerHour float64 `json:"price_per_hour"`
	Rating       float64 `json:"rating"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type TutorRepository interface {
	Create(ctx context.Context, tutor *Tutor) error
	GetByID(ctx context.Context, id int) (*Tutor, error)
	List(ctx context.Context, subject string) ([]*Tutor, error)
}

type TutorUsecase interface {
	BecomeTutor(ctx context.Context, tutor *Tutor) error
	GetTutorProfile(ctx context.Context, id int) (*Tutor, error)
	SearchTutors(ctx context.Context, subject string) ([]*Tutor, error)
}
