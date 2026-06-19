package domain

import "context"

type User struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Roles        []string `json:"roles"` // e.g. ["student", "tutor"]
	PasswordHash *string `json:"-"`
	GoogleID     *string `json:"google_id,omitempty"`
	AppleID      *string `json:"apple_id,omitempty"`
	AvatarURL    *string `json:"avatar_url,omitempty"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)
	GetByAppleID(ctx context.Context, appleID string) (*User, error)
	Delete(ctx context.Context, id int) error
}

type UserUsecase interface {
	Register(ctx context.Context, user *User) error
	GetProfile(ctx context.Context, id int) (*User, error)
	DeleteAccount(ctx context.Context, id int) error
}
