package usecase

import (
	"context"
	"fmt"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type followUsecase struct {
	repo domain.FollowRepository
}

func NewFollowUsecase(repo domain.FollowRepository) domain.FollowUsecase {
	return &followUsecase{repo: repo}
}

func (u *followUsecase) ToggleFollow(ctx context.Context, followerID, followingID int) (bool, error) {
	if followerID == followingID {
		return false, fmt.Errorf("cannot follow yourself")
	}
	return u.repo.ToggleFollow(ctx, followerID, followingID)
}

func (u *followUsecase) GetFollowStats(ctx context.Context, userID, requesterID int) (*domain.FollowStats, error) {
	followersCount, err := u.repo.CountFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}
	followingCount, err := u.repo.CountFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}
	isFollowing := false
	if requesterID > 0 && requesterID != userID {
		isFollowing, err = u.repo.IsFollowing(ctx, requesterID, userID)
		if err != nil {
			return nil, err
		}
	}
	return &domain.FollowStats{
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		IsFollowing:    isFollowing,
	}, nil
}

func (u *followUsecase) ListFollowers(ctx context.Context, userID int) ([]*domain.User, error) {
	return u.repo.GetFollowers(ctx, userID)
}

func (u *followUsecase) ListFollowing(ctx context.Context, userID int) ([]*domain.User, error) {
	return u.repo.GetFollowing(ctx, userID)
}
