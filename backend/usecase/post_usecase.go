package usecase

import (
	"context"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postUsecase struct {
	postRepo domain.PostRepository
}

func NewPostUsecase(repo domain.PostRepository) domain.PostUsecase {
	return &postUsecase{postRepo: repo}
}

func (u *postUsecase) CreatePost(ctx context.Context, post *domain.Post) error {
	return u.postRepo.Create(ctx, post)
}

func (u *postUsecase) GetPostDetails(ctx context.Context, id int, requesterUserID int) (*domain.Post, []*domain.Comment, error) {
	post, err := u.postRepo.GetByID(ctx, id, requesterUserID)
	if err != nil {
		return nil, nil, err
	}
	if post == nil {
		return nil, nil, nil
	}

	comments, err := u.postRepo.GetComments(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return post, comments, nil
}

func (u *postUsecase) ListFeed(ctx context.Context, subject string, requesterUserID int) ([]*domain.Post, error) {
	return u.postRepo.List(ctx, subject, requesterUserID)
}

func (u *postUsecase) ToggleLike(ctx context.Context, postID int, userID int) (bool, error) {
	return u.postRepo.Like(ctx, postID, userID)
}

func (u *postUsecase) ToggleSave(ctx context.Context, postID int, userID int) (bool, error) {
	return u.postRepo.Save(ctx, postID, userID)
}

func (u *postUsecase) AddPostComment(ctx context.Context, comment *domain.Comment) error {
	return u.postRepo.AddComment(ctx, comment)
}
