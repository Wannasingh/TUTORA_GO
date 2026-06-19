package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresPostRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPostRepository(db *pgxpool.Pool) domain.PostRepository {
	return &postgresPostRepository{db: db}
}

func (r *postgresPostRepository) Create(ctx context.Context, post *domain.Post) error {
	query := `INSERT INTO tutora_app.posts (user_id, subject, title, body, image_url, video_url) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, post.UserID, post.Subject, post.Title, post.Body, post.ImageURL, post.VideoURL).
		Scan(&post.ID, &createdAt, &updatedAt)
	if err == nil {
		post.CreatedAt = createdAt.Format(time.RFC3339)
		post.UpdatedAt = updatedAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresPostRepository) GetByID(ctx context.Context, id int, requesterUserID int) (*domain.Post, error) {
	query := `SELECT p.id, p.user_id, p.subject, p.title, p.body, p.image_url, p.video_url, p.created_at, p.updated_at,
	                 u.name, u.email, u.roles, u.avatar_url,
	                 (SELECT COUNT(*) FROM tutora_app.post_likes WHERE post_id = p.id) as likes_count,
	                 (SELECT COUNT(*) FROM tutora_app.comments WHERE post_id = p.id) as comments_count,
	                 (SELECT COUNT(*) FROM tutora_app.post_saves WHERE post_id = p.id) as saves_count,
	                 EXISTS (SELECT 1 FROM tutora_app.post_likes WHERE post_id = p.id AND user_id = $2) as is_liked,
	                 EXISTS (SELECT 1 FROM tutora_app.post_saves WHERE post_id = p.id AND user_id = $2) as is_saved
	          FROM tutora_app.posts p
	          JOIN tutora_app.users u ON p.user_id = u.id
	          WHERE p.id = $1`
	post := &domain.Post{User: &domain.User{}}
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id, requesterUserID).
		Scan(&post.ID, &post.UserID, &post.Subject, &post.Title, &post.Body, &post.ImageURL, &post.VideoURL, &createdAt, &updatedAt,
			&post.User.Name, &post.User.Email, &post.User.Roles, &post.User.AvatarURL,
			&post.LikesCount, &post.CommentsCount, &post.SavesCount,
			&post.IsLiked, &post.IsSaved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	post.User.ID = post.UserID
	post.CreatedAt = createdAt.Format(time.RFC3339)
	post.UpdatedAt = updatedAt.Format(time.RFC3339)
	return post, nil
}

func (r *postgresPostRepository) List(ctx context.Context, subject string, requesterUserID int) ([]*domain.Post, error) {
	var rows pgx.Rows
	var err error

	baseQuery := `SELECT p.id, p.user_id, p.subject, p.title, p.body, p.image_url, p.video_url, p.created_at, p.updated_at,
	                     u.name, u.email, u.roles, u.avatar_url,
	                     (SELECT COUNT(*) FROM tutora_app.post_likes WHERE post_id = p.id) as likes_count,
	                     (SELECT COUNT(*) FROM tutora_app.comments WHERE post_id = p.id) as comments_count,
	                     (SELECT COUNT(*) FROM tutora_app.post_saves WHERE post_id = p.id) as saves_count,
	                     EXISTS (SELECT 1 FROM tutora_app.post_likes WHERE post_id = p.id AND user_id = $1) as is_liked,
	                     EXISTS (SELECT 1 FROM tutora_app.post_saves WHERE post_id = p.id AND user_id = $1) as is_saved
	              FROM tutora_app.posts p
	              JOIN tutora_app.users u ON p.user_id = u.id`

	var orderClause string
	if requesterUserID > 0 {
		orderClause = ` ORDER BY (
		                     (
		                       (SELECT COUNT(*) FROM tutora_app.post_likes WHERE post_id = p.id) * 1.0 +
		                       (SELECT COUNT(*) FROM tutora_app.comments WHERE post_id = p.id) * 3.0 +
		                       (SELECT COUNT(*) FROM tutora_app.post_saves WHERE post_id = p.id) * 2.0 + 1.0
		                     ) / 
		                     POWER(EXTRACT(EPOCH FROM (NOW() - p.created_at)) / 3600 + 2, 1.5)
		                   ) * (
		                     CASE WHEN p.subject IN (
		                       SELECT subject FROM tutora_app.tutors WHERE user_id = $1
		                       UNION
		                       SELECT p2.subject FROM tutora_app.post_likes pl JOIN tutora_app.posts p2 ON pl.post_id = p2.id WHERE pl.user_id = $1
		                       UNION
		                       SELECT p3.subject FROM tutora_app.post_saves ps JOIN tutora_app.posts p3 ON ps.post_id = p3.id WHERE ps.user_id = $1
		                     ) THEN 1.5 ELSE 1.0 END
		                   ) DESC, p.created_at DESC`
	} else {
		orderClause = ` ORDER BY (
		                    (
		                      (SELECT COUNT(*) FROM tutora_app.post_likes WHERE post_id = p.id) * 1.0 +
		                      (SELECT COUNT(*) FROM tutora_app.comments WHERE post_id = p.id) * 3.0 +
		                      (SELECT COUNT(*) FROM tutora_app.post_saves WHERE post_id = p.id) * 2.0 + 1.0
		                    ) / 
		                    POWER(EXTRACT(EPOCH FROM (NOW() - p.created_at)) / 3600 + 2, 1.5)
		                  ) DESC, p.created_at DESC`
	}

	if subject != "" {
		query := baseQuery + ` WHERE p.subject ILIKE $2` + orderClause
		rows, err = r.db.Query(ctx, query, requesterUserID, "%"+subject+"%")
	} else {
		query := baseQuery + orderClause
		rows, err = r.db.Query(ctx, query, requesterUserID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		p := &domain.Post{User: &domain.User{}}
		var createdAt, updatedAt time.Time
		err := rows.Scan(&p.ID, &p.UserID, &p.Subject, &p.Title, &p.Body, &p.ImageURL, &p.VideoURL, &createdAt, &updatedAt,
			&p.User.Name, &p.User.Email, &p.User.Roles, &p.User.AvatarURL,
			&p.LikesCount, &p.CommentsCount, &p.SavesCount,
			&p.IsLiked, &p.IsSaved)
		if err != nil {
			return nil, err
		}
		p.User.ID = p.UserID
		p.CreatedAt = createdAt.Format(time.RFC3339)
		p.UpdatedAt = updatedAt.Format(time.RFC3339)
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *postgresPostRepository) Like(ctx context.Context, postID int, userID int) (bool, error) {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tutora_app.post_likes WHERE post_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, checkQuery, postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		deleteQuery := `DELETE FROM tutora_app.post_likes WHERE post_id = $1 AND user_id = $2`
		_, err = r.db.Exec(ctx, deleteQuery, postID, userID)
		return false, err
	} else {
		insertQuery := `INSERT INTO tutora_app.post_likes (post_id, user_id) VALUES ($1, $2)`
		_, err = r.db.Exec(ctx, insertQuery, postID, userID)
		return true, err
	}
}

func (r *postgresPostRepository) Save(ctx context.Context, postID int, userID int) (bool, error) {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tutora_app.post_saves WHERE post_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, checkQuery, postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		deleteQuery := `DELETE FROM tutora_app.post_saves WHERE post_id = $1 AND user_id = $2`
		_, err = r.db.Exec(ctx, deleteQuery, postID, userID)
		return false, err
	} else {
		insertQuery := `INSERT INTO tutora_app.post_saves (post_id, user_id) VALUES ($1, $2)`
		_, err = r.db.Exec(ctx, insertQuery, postID, userID)
		return true, err
	}
}

func (r *postgresPostRepository) AddComment(ctx context.Context, comment *domain.Comment) error {
	query := `INSERT INTO tutora_app.comments (post_id, user_id, body, image_url, parent_id, status) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	var createdAt, updatedAt time.Time
	comment.Status = "active"
	err := r.db.QueryRow(ctx, query, comment.PostID, comment.UserID, comment.Body, comment.ImageURL, comment.ParentID, comment.Status).
		Scan(&comment.ID, &createdAt, &updatedAt)
	if err == nil {
		comment.CreatedAt = createdAt.Format(time.RFC3339)
		comment.UpdatedAt = updatedAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresPostRepository) GetComments(ctx context.Context, postID int, requesterUserID int) ([]*domain.Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.body, c.image_url, c.parent_id, c.status, c.created_at, c.updated_at,
	                 u.name, u.email, u.roles, u.avatar_url,
	                 (SELECT COUNT(*) FROM tutora_app.comment_likes WHERE comment_id = c.id) as likes_count,
	                 EXISTS (SELECT 1 FROM tutora_app.comment_likes WHERE comment_id = c.id AND user_id = $2) as is_liked
	          FROM tutora_app.comments c
	          JOIN tutora_app.users u ON c.user_id = u.id
	          WHERE c.post_id = $1
	          ORDER BY 
	            (EXISTS (SELECT 1 FROM tutora_app.users u2 WHERE u2.id = c.user_id AND 'tutor' = ANY(u2.roles))) DESC, 
	            (SELECT COUNT(*) FROM tutora_app.comment_likes WHERE comment_id = c.id) DESC, 
	            c.created_at ASC`
	rows, err := r.db.Query(ctx, query, postID, requesterUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		c := &domain.Comment{User: &domain.User{}}
		var createdAt, updatedAt time.Time
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Body, &c.ImageURL, &c.ParentID, &c.Status, &createdAt, &updatedAt,
			&c.User.Name, &c.User.Email, &c.User.Roles, &c.User.AvatarURL, &c.LikesCount, &c.IsLiked)
		if err != nil {
			return nil, err
		}
		c.User.ID = c.UserID
		c.CreatedAt = createdAt.Format(time.RFC3339)
		c.UpdatedAt = updatedAt.Format(time.RFC3339)
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *postgresPostRepository) LikeComment(ctx context.Context, commentID int, userID int) (bool, error) {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tutora_app.comment_likes WHERE comment_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, checkQuery, commentID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		deleteQuery := `DELETE FROM tutora_app.comment_likes WHERE comment_id = $1 AND user_id = $2`
		_, err = r.db.Exec(ctx, deleteQuery, commentID, userID)
		return false, err
	} else {
		insertQuery := `INSERT INTO tutora_app.comment_likes (comment_id, user_id) VALUES ($1, $2)`
		_, err = r.db.Exec(ctx, insertQuery, commentID, userID)
		return true, err
	}
}

func (r *postgresPostRepository) DeleteComment(ctx context.Context, id int) error {
	query := `UPDATE tutora_app.comments SET status = 'deleted', body = '[This comment has been deleted.]', image_url = NULL WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *postgresPostRepository) GetCommentByID(ctx context.Context, id int) (*domain.Comment, error) {
	query := `SELECT id, post_id, user_id, body, status FROM tutora_app.comments WHERE id = $1`
	c := &domain.Comment{}
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.PostID, &c.UserID, &c.Body, &c.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

func (r *postgresPostRepository) CreateReport(ctx context.Context, report *domain.Report) error {
	query := `INSERT INTO tutora_app.reports (reporter_id, target_type, target_id, reason) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	var createdAt time.Time
	err := r.db.QueryRow(ctx, query, report.ReporterID, report.TargetType, report.TargetID, report.Reason).
		Scan(&report.ID, &createdAt)
	if err == nil {
		report.CreatedAt = createdAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresPostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `UPDATE tutora_app.posts SET subject = $1, title = $2, body = $3, image_url = $4, video_url = $5, updated_at = NOW() 
	          WHERE id = $6`
	_, err := r.db.Exec(ctx, query, post.Subject, post.Title, post.Body, post.ImageURL, post.VideoURL, post.ID)
	return err
}

func (r *postgresPostRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tutora_app.posts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *postgresPostRepository) Repost(ctx context.Context, postID int, userID int) (bool, error) {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tutora_app.reposts WHERE post_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, checkQuery, postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		_, err = r.db.Exec(ctx, `DELETE FROM tutora_app.reposts WHERE post_id = $1 AND user_id = $2`, postID, userID)
		return false, err
	} else {
		_, err = r.db.Exec(ctx, `INSERT INTO tutora_app.reposts (user_id, post_id) VALUES ($1, $2)`, userID, postID)
		return true, err
	}
}

// userPostBaseQuery is the shared SELECT for user profile tab queries
const userPostBaseQuery = `SELECT p.id, p.user_id, p.subject, p.title, p.body, p.image_url, p.video_url, p.created_at, p.updated_at,
                     u.name, u.email, u.roles, u.avatar_url,
                     (SELECT COUNT(*) FROM tutora_app.post_likes WHERE post_id = p.id) as likes_count,
                     (SELECT COUNT(*) FROM tutora_app.comments WHERE post_id = p.id) as comments_count,
                     (SELECT COUNT(*) FROM tutora_app.post_saves WHERE post_id = p.id) as saves_count,
                     EXISTS (SELECT 1 FROM tutora_app.post_likes WHERE post_id = p.id AND user_id = $1) as is_liked,
                     EXISTS (SELECT 1 FROM tutora_app.post_saves WHERE post_id = p.id AND user_id = $1) as is_saved
              FROM tutora_app.posts p
              JOIN tutora_app.users u ON p.user_id = u.id`

func (r *postgresPostRepository) scanPostRows(rows pgx.Rows) ([]*domain.Post, error) {
	var posts []*domain.Post
	for rows.Next() {
		p := &domain.Post{User: &domain.User{}}
		var createdAt, updatedAt time.Time
		err := rows.Scan(&p.ID, &p.UserID, &p.Subject, &p.Title, &p.Body, &p.ImageURL, &p.VideoURL, &createdAt, &updatedAt,
			&p.User.Name, &p.User.Email, &p.User.Roles, &p.User.AvatarURL,
			&p.LikesCount, &p.CommentsCount, &p.SavesCount,
			&p.IsLiked, &p.IsSaved)
		if err != nil {
			return nil, err
		}
		p.User.ID = p.UserID
		p.CreatedAt = createdAt.Format(time.RFC3339)
		p.UpdatedAt = updatedAt.Format(time.RFC3339)
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *postgresPostRepository) GetUserPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	query := userPostBaseQuery + ` WHERE p.user_id = $2 ORDER BY p.created_at DESC`
	rows, err := r.db.Query(ctx, query, requesterUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanPostRows(rows)
}

func (r *postgresPostRepository) GetUserLikedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	query := userPostBaseQuery + ` JOIN tutora_app.post_likes pl ON pl.post_id = p.id AND pl.user_id = $2 ORDER BY pl.created_at DESC`
	rows, err := r.db.Query(ctx, query, requesterUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanPostRows(rows)
}

func (r *postgresPostRepository) GetUserSavedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	query := userPostBaseQuery + ` JOIN tutora_app.post_saves ps ON ps.post_id = p.id AND ps.user_id = $2 ORDER BY ps.created_at DESC`
	rows, err := r.db.Query(ctx, query, requesterUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanPostRows(rows)
}

func (r *postgresPostRepository) GetUserRepostedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	query := userPostBaseQuery + ` JOIN tutora_app.reposts rp ON rp.post_id = p.id AND rp.user_id = $2 ORDER BY rp.created_at DESC`
	rows, err := r.db.Query(ctx, query, requesterUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanPostRows(rows)
}

