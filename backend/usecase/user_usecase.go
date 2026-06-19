package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type userUsecase struct {
	userRepo   domain.UserRepository
	followRepo domain.FollowRepository
}

func NewUserUsecase(repo domain.UserRepository, followRepo ...domain.FollowRepository) domain.UserUsecase {
	u := &userUsecase{userRepo: repo}
	if len(followRepo) > 0 {
		u.followRepo = followRepo[0]
	}
	return u
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

func (u *userUsecase) GetFullProfile(ctx context.Context, userID, requesterID int) (*domain.UserProfile, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	profile := &domain.UserProfile{User: user}

	// Follow stats
	if u.followRepo != nil {
		followersCount, _ := u.followRepo.CountFollowers(ctx, userID)
		followingCount, _ := u.followRepo.CountFollowing(ctx, userID)
		isFollowing := false
		if requesterID > 0 && requesterID != userID {
			isFollowing, _ = u.followRepo.IsFollowing(ctx, requesterID, userID)
		}
		profile.FollowStats = &domain.FollowStats{
			FollowersCount: followersCount,
			FollowingCount: followingCount,
			IsFollowing:    isFollowing,
		}
	}

	// Post counts
	profile.PostsCount, _ = u.userRepo.CountUserPosts(ctx, userID)
	profile.RepostsCount, _ = u.userRepo.CountUserReposts(ctx, userID)

	return profile, nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, userID int, req *domain.UpdateProfileRequest) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Apply updates only for provided fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.CoverURL != nil {
		user.CoverURL = req.CoverURL
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.School != nil {
		user.School = req.School
	}
	if req.Birthdate != nil {
		user.Birthdate = req.Birthdate
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) DeleteAccount(ctx context.Context, id int) error {
	return u.userRepo.Delete(ctx, id)
}
