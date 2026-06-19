package domain

import "context"

type Post struct {
	ID            int      `json:"id"`
	UserID        int      `json:"user_id"`
	User          *User    `json:"user,omitempty"`
	Subject       string   `json:"subject"`
	Title         string   `json:"title"`
	Body          string   `json:"body"`
	ImageURL      *string  `json:"image_url,omitempty"`
	LikesCount    int      `json:"likes_count"`
	CommentsCount int      `json:"comments_count"`
	SavesCount    int      `json:"saves_count"`
	IsLiked       bool     `json:"is_liked"`
	IsSaved       bool     `json:"is_saved"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type Comment struct {
	ID        int     `json:"id"`
	PostID    int     `json:"post_id"`
	UserID    int     `json:"user_id"`
	User      *User   `json:"user,omitempty"`
	Body      string  `json:"body"`
	ImageURL  *string `json:"image_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type CreatePostRequest struct {
	Subject  string  `json:"subject" binding:"required"`
	Title    string  `json:"title" binding:"required"`
	Body     string  `json:"body" binding:"required"`
	ImageURL *string `json:"image_url,omitempty"`
}

type CreateCommentRequest struct {
	Body     string  `json:"body" binding:"required"`
	ImageURL *string `json:"image_url,omitempty"`
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int, requesterUserID int) (*Post, error)
	List(ctx context.Context, subject string, requesterUserID int) ([]*Post, error)
	Like(ctx context.Context, postID int, userID int) (bool, error)   // Returns true if liked, false if unliked
	Save(ctx context.Context, postID int, userID int) (bool, error)   // Returns true if saved, false if unsaved
	AddComment(ctx context.Context, comment *Comment) error
	GetComments(ctx context.Context, postID int) ([]*Comment, error)
}

type PostUsecase interface {
	CreatePost(ctx context.Context, post *Post) error
	GetPostDetails(ctx context.Context, id int, requesterUserID int) (*Post, []*Comment, error)
	ListFeed(ctx context.Context, subject string, requesterUserID int) ([]*Post, error)
	ToggleLike(ctx context.Context, postID int, userID int) (bool, error)
	ToggleSave(ctx context.Context, postID int, userID int) (bool, error)
	AddPostComment(ctx context.Context, comment *Comment) error
}
