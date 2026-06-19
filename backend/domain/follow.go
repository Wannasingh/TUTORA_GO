package domain

import "context"

type FollowStats struct {
	FollowersCount int  `json:"followers_count"`
	FollowingCount int  `json:"following_count"`
	IsFollowing    bool `json:"is_following"`
}

type FollowRepository interface {
	ToggleFollow(ctx context.Context, followerID, followingID int) (bool, error)
	IsFollowing(ctx context.Context, followerID, followingID int) (bool, error)
	GetFollowers(ctx context.Context, userID int) ([]*User, error)
	GetFollowing(ctx context.Context, userID int) ([]*User, error)
	CountFollowers(ctx context.Context, userID int) (int, error)
	CountFollowing(ctx context.Context, userID int) (int, error)
}

type FollowUsecase interface {
	ToggleFollow(ctx context.Context, followerID, followingID int) (bool, error)
	GetFollowStats(ctx context.Context, userID, requesterID int) (*FollowStats, error)
	ListFollowers(ctx context.Context, userID int) ([]*User, error)
	ListFollowing(ctx context.Context, userID int) ([]*User, error)
}
