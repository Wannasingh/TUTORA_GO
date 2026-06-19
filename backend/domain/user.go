package domain

import "context"

type User struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"` // e.g. ["student", "tutor"]
	PasswordHash *string  `json:"-"`
	GoogleID     *string  `json:"google_id,omitempty"`
	AppleID      *string  `json:"apple_id,omitempty"`
	AvatarURL    *string  `json:"avatar_url,omitempty"`
	Bio          *string  `json:"bio,omitempty"`
	CoverURL     *string  `json:"cover_url,omitempty"`
	Phone        *string  `json:"phone,omitempty"`
	School       *string  `json:"school,omitempty"`
	Birthdate    *string  `json:"birthdate,omitempty"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	CoverURL  *string `json:"cover_url,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	School    *string `json:"school,omitempty"`
	Birthdate *string `json:"birthdate,omitempty"`
}

type UserProfile struct {
	User           *User        `json:"user"`
	FollowStats    *FollowStats `json:"follow_stats"`
	PostsCount     int          `json:"posts_count"`
	RepostsCount   int          `json:"reposts_count"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)
	GetByAppleID(ctx context.Context, appleID string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
	CountUserPosts(ctx context.Context, userID int) (int, error)
	CountUserReposts(ctx context.Context, userID int) (int, error)
}

type UserUsecase interface {
	Register(ctx context.Context, user *User) error
	GetProfile(ctx context.Context, id int) (*User, error)
	GetFullProfile(ctx context.Context, userID, requesterID int) (*UserProfile, error)
	UpdateProfile(ctx context.Context, userID int, req *UpdateProfileRequest) (*User, error)
	DeleteAccount(ctx context.Context, id int) error
}
