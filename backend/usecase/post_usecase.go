package usecase

import (
	"context"
	"errors"

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

	flatComments, err := u.postRepo.GetComments(ctx, id, requesterUserID)
	if err != nil {
		return nil, nil, err
	}

	// Build nested comments tree structure dynamically
	var rootComments []*domain.Comment
	commentMap := make(map[int]*domain.Comment)

	for _, c := range flatComments {
		c.Replies = make([]*domain.Comment, 0)
		commentMap[c.ID] = c
	}

	for _, c := range flatComments {
		if c.ParentID == nil {
			rootComments = append(rootComments, c)
		} else {
			if parent, found := commentMap[*c.ParentID]; found {
				parent.Replies = append(parent.Replies, c)
			} else {
				rootComments = append(rootComments, c)
			}
		}
	}

	return post, rootComments, nil
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
	if comment.Body == "" {
		return errors.New("comment body cannot be empty")
	}
	
	// If parent ID is provided, verify it exists and belongs to the same post
	if comment.ParentID != nil {
		parent, err := u.postRepo.GetCommentByID(ctx, *comment.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("parent comment not found")
		}
		if parent.PostID != comment.PostID {
			return errors.New("reply must belong to the same post as the parent comment")
		}
	}

	return u.postRepo.AddComment(ctx, comment)
}

func (u *postUsecase) ToggleCommentLike(ctx context.Context, commentID int, userID int) (bool, error) {
	return u.postRepo.LikeComment(ctx, commentID, userID)
}

func (u *postUsecase) DeleteComment(ctx context.Context, userID, commentID int) error {
	comment, err := u.postRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("comment not found")
	}

	// Verify authorization: only the author can delete
	if comment.UserID != userID {
		return errors.New("unauthorized: you can only delete your own comments")
	}

	return u.postRepo.DeleteComment(ctx, commentID)
}

func (u *postUsecase) ReportContent(ctx context.Context, report *domain.Report) error {
	if report.Reason == "" {
		return errors.New("report reason cannot be empty")
	}
	if report.TargetType != "post" && report.TargetType != "comment" {
		return errors.New("invalid target type: must be 'post' or 'comment'")
	}
	return u.postRepo.CreateReport(ctx, report)
}

func (u *postUsecase) UpdatePost(ctx context.Context, userID int, post *domain.Post) error {
	existing, err := u.postRepo.GetByID(ctx, post.ID, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("post not found")
	}

	if existing.UserID != userID {
		return errors.New("unauthorized: you can only edit your own posts")
	}

	existing.Title = post.Title
	existing.Body = post.Body
	existing.Subject = post.Subject
	existing.ImageURL = post.ImageURL
	existing.VideoURL = post.VideoURL

	return u.postRepo.Update(ctx, existing)
}

func (u *postUsecase) DeletePost(ctx context.Context, userID, postID int) error {
	existing, err := u.postRepo.GetByID(ctx, postID, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("post not found")
	}

	if existing.UserID != userID {
		return errors.New("unauthorized: you can only delete your own posts")
	}

	return u.postRepo.Delete(ctx, postID)
}

func (u *postUsecase) ToggleRepost(ctx context.Context, postID int, userID int) (bool, error) {
	return u.postRepo.Repost(ctx, postID, userID)
}

func (u *postUsecase) CreateQuotePost(ctx context.Context, post *domain.Post) error {
	if post.OriginalPostID == nil {
		return errors.New("quote post requires an original_post_id")
	}
	return u.postRepo.Create(ctx, post)
}

func (u *postUsecase) GetUserPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	return u.postRepo.GetUserPosts(ctx, userID, requesterUserID)
}

func (u *postUsecase) GetUserLikedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	return u.postRepo.GetUserLikedPosts(ctx, userID, requesterUserID)
}

func (u *postUsecase) GetUserSavedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	return u.postRepo.GetUserSavedPosts(ctx, userID, requesterUserID)
}

func (u *postUsecase) GetUserRepostedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	return u.postRepo.GetUserRepostedPosts(ctx, userID, requesterUserID)
}

