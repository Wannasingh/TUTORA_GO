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
	query := `INSERT INTO tutora_app.posts (user_id, subject, title, body, image_url) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, post.UserID, post.Subject, post.Title, post.Body, post.ImageURL).
		Scan(&post.ID, &createdAt, &updatedAt)
	if err == nil {
		post.CreatedAt = createdAt.Format(time.RFC3339)
		post.UpdatedAt = updatedAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresPostRepository) GetByID(ctx context.Context, id int, requesterUserID int) (*domain.Post, error) {
	query := `SELECT p.id, p.user_id, p.subject, p.title, p.body, p.image_url, p.created_at, p.updated_at,
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
		Scan(&post.ID, &post.UserID, &post.Subject, &post.Title, &post.Body, &post.ImageURL, &createdAt, &updatedAt,
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

	baseQuery := `SELECT p.id, p.user_id, p.subject, p.title, p.body, p.image_url, p.created_at, p.updated_at,
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
		err := rows.Scan(&p.ID, &p.UserID, &p.Subject, &p.Title, &p.Body, &p.ImageURL, &createdAt, &updatedAt,
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
	// Check if already liked
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tutora_app.post_likes WHERE post_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, checkQuery, postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		// Unlike
		deleteQuery := `DELETE FROM tutora_app.post_likes WHERE post_id = $1 AND user_id = $2`
		_, err = r.db.Exec(ctx, deleteQuery, postID, userID)
		return false, err
	} else {
		// Like
		insertQuery := `INSERT INTO tutora_app.post_likes (post_id, user_id) VALUES ($1, $2)`
		_, err = r.db.Exec(ctx, insertQuery, postID, userID)
		return true, err
	}
}

func (r *postgresPostRepository) Save(ctx context.Context, postID int, userID int) (bool, error) {
	// Check if already saved
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tutora_app.post_saves WHERE post_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, checkQuery, postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		// Unsave
		deleteQuery := `DELETE FROM tutora_app.post_saves WHERE post_id = $1 AND user_id = $2`
		_, err = r.db.Exec(ctx, deleteQuery, postID, userID)
		return false, err
	} else {
		// Save
		insertQuery := `INSERT INTO tutora_app.post_saves (post_id, user_id) VALUES ($1, $2)`
		_, err = r.db.Exec(ctx, insertQuery, postID, userID)
		return true, err
	}
}

func (r *postgresPostRepository) AddComment(ctx context.Context, comment *domain.Comment) error {
	query := `INSERT INTO tutora_app.comments (post_id, user_id, body, image_url) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, comment.PostID, comment.UserID, comment.Body, comment.ImageURL).
		Scan(&comment.ID, &createdAt, &updatedAt)
	if err == nil {
		comment.CreatedAt = createdAt.Format(time.RFC3339)
		comment.UpdatedAt = updatedAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresPostRepository) GetComments(ctx context.Context, postID int) ([]*domain.Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.body, c.image_url, c.created_at, c.updated_at,
	                 u.name, u.email, u.roles, u.avatar_url
	          FROM tutora_app.comments c
	          JOIN tutora_app.users u ON c.user_id = u.id
	          WHERE c.post_id = $1
	          ORDER BY c.created_at ASC`
	rows, err := r.db.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		c := &domain.Comment{User: &domain.User{}}
		var createdAt, updatedAt time.Time
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Body, &c.ImageURL, &createdAt, &updatedAt,
			&c.User.Name, &c.User.Email, &c.User.Roles, &c.User.AvatarURL)
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
