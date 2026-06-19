package domain

import "context"

type Post struct {
	ID             int     `json:"id"`
	UserID         int     `json:"user_id"`
	User           *User   `json:"user,omitempty"`
	Subject        string  `json:"subject"`
	Title          string  `json:"title"`
	Body           string  `json:"body"`
	ImageURL       *string `json:"image_url,omitempty"`
	VideoURL       *string `json:"video_url,omitempty"`
	OriginalPostID *int    `json:"original_post_id,omitempty"`
	OriginalPost   *Post   `json:"original_post,omitempty"`
	LikesCount     int     `json:"likes_count"`
	CommentsCount  int     `json:"comments_count"`
	SavesCount     int     `json:"saves_count"`
	RepostCount    int     `json:"repost_count"`
	IsLiked        bool    `json:"is_liked"`
	IsSaved        bool    `json:"is_saved"`
	IsReposted     bool    `json:"is_reposted"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type Comment struct {
	ID         int        `json:"id"`
	PostID     int        `json:"post_id"`
	UserID     int        `json:"user_id"`
	User       *User      `json:"user,omitempty"`
	Body       string     `json:"body"`
	ImageURL   *string    `json:"image_url,omitempty"`
	ParentID   *int       `json:"parent_id,omitempty"`
	Status     string     `json:"status"` // 'active', 'deleted'
	LikesCount int        `json:"likes_count"`
	IsLiked    bool       `json:"is_liked"`
	Replies    []*Comment `json:"replies,omitempty"`
	CreatedAt  string     `json:"created_at"`
	UpdatedAt  string     `json:"updated_at"`
}

type Report struct {
	ID         int    `json:"id"`
	ReporterID int    `json:"reporter_id"`
	TargetType string `json:"target_type"` // 'post', 'comment'
	TargetID   int    `json:"target_id"`
	Reason     string `json:"reason"`
	CreatedAt  string `json:"created_at"`
}

type CreatePostRequest struct {
	Subject        string  `json:"subject" binding:"required"`
	Title          string  `json:"title" binding:"required"`
	Body           string  `json:"body" binding:"required"`
	ImageURL       *string `json:"image_url,omitempty"`
	VideoURL       *string `json:"video_url,omitempty"`
	OriginalPostID *int    `json:"original_post_id,omitempty"`
}

type CreateCommentRequest struct {
	Body     string  `json:"body" binding:"required"`
	ImageURL *string `json:"image_url,omitempty"`
	ParentID *int    `json:"parent_id,omitempty"`
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int, requesterUserID int) (*Post, error)
	List(ctx context.Context, subject string, requesterUserID int) ([]*Post, error)
	Like(ctx context.Context, postID int, userID int) (bool, error)
	Save(ctx context.Context, postID int, userID int) (bool, error)
	Repost(ctx context.Context, postID int, userID int) (bool, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int) error

	// User-specific post queries
	GetUserPosts(ctx context.Context, userID, requesterUserID int) ([]*Post, error)
	GetUserLikedPosts(ctx context.Context, userID, requesterUserID int) ([]*Post, error)
	GetUserSavedPosts(ctx context.Context, userID, requesterUserID int) ([]*Post, error)
	GetUserRepostedPosts(ctx context.Context, userID, requesterUserID int) ([]*Post, error)

	// Comment operations
	AddComment(ctx context.Context, comment *Comment) error
	GetComments(ctx context.Context, postID int, requesterUserID int) ([]*Comment, error)
	LikeComment(ctx context.Context, commentID int, userID int) (bool, error)
	DeleteComment(ctx context.Context, id int) error
	GetCommentByID(ctx context.Context, id int) (*Comment, error)

	// Moderation operations
	CreateReport(ctx context.Context, report *Report) error
}

type PostUsecase interface {
	CreatePost(ctx context.Context, post *Post) error
	GetPostDetails(ctx context.Context, id int, requesterUserID int) (*Post, []*Comment, error)
	ListFeed(ctx context.Context, subject string, requesterUserID int) ([]*Post, error)
	ToggleLike(ctx context.Context, postID int, userID int) (bool, error)
	ToggleSave(ctx context.Context, postID int, userID int) (bool, error)
	ToggleRepost(ctx context.Context, postID int, userID int) (bool, error)
	CreateQuotePost(ctx context.Context, post *Post) error
	UpdatePost(ctx context.Context, userID int, post *Post) error
	DeletePost(ctx context.Context, userID, postID int) error

	// User profile tabs
	GetUserPosts(ctx context.Context, userID, requesterUserID int) ([]*Post, error)
	GetUserLikedPosts(ctx context.Context, userID, requesterUserID int) ([]*Post, error)
	GetUserSavedPosts(ctx context.Context, userID, requesterUserID int) ([]*Post, error)
	GetUserRepostedPosts(ctx context.Context, userID, requesterUserID int) ([]*Post, error)

	// Comments
	AddPostComment(ctx context.Context, comment *Comment) error
	ToggleCommentLike(ctx context.Context, commentID int, userID int) (bool, error)
	DeleteComment(ctx context.Context, userID, commentID int) error

	// Moderation
	ReportContent(ctx context.Context, report *Report) error
}
