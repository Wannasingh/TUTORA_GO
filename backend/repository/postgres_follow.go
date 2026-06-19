package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresFollowRepository struct {
	db *pgxpool.Pool
}

func NewPostgresFollowRepository(db *pgxpool.Pool) domain.FollowRepository {
	return &postgresFollowRepository{db: db}
}

func (r *postgresFollowRepository) ToggleFollow(ctx context.Context, followerID, followingID int) (bool, error) {
	// Try to delete first; if nothing deleted, insert
	deleteQuery := `DELETE FROM tutora_app.follows WHERE follower_id = $1 AND following_id = $2`
	res, err := r.db.Exec(ctx, deleteQuery, followerID, followingID)
	if err != nil {
		return false, err
	}
	if res.RowsAffected() > 0 {
		return false, nil // unfollowed
	}

	insertQuery := `INSERT INTO tutora_app.follows (follower_id, following_id) VALUES ($1, $2)`
	_, err = r.db.Exec(ctx, insertQuery, followerID, followingID)
	if err != nil {
		return false, err
	}
	return true, nil // followed
}

func (r *postgresFollowRepository) IsFollowing(ctx context.Context, followerID, followingID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tutora_app.follows WHERE follower_id = $1 AND following_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *postgresFollowRepository) GetFollowers(ctx context.Context, userID int) ([]*domain.User, error) {
	query := `SELECT u.id, u.name, u.email, u.roles, u.avatar_url, u.bio, u.created_at, u.updated_at
	          FROM tutora_app.users u
	          JOIN tutora_app.follows f ON f.follower_id = u.id
	          WHERE f.following_id = $1
	          ORDER BY f.created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Roles, &u.AvatarURL, &u.Bio, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		u.CreatedAt = createdAt.Format(time.RFC3339)
		u.UpdatedAt = updatedAt.Format(time.RFC3339)
		users = append(users, u)
	}
	return users, nil
}

func (r *postgresFollowRepository) GetFollowing(ctx context.Context, userID int) ([]*domain.User, error) {
	query := `SELECT u.id, u.name, u.email, u.roles, u.avatar_url, u.bio, u.created_at, u.updated_at
	          FROM tutora_app.users u
	          JOIN tutora_app.follows f ON f.following_id = u.id
	          WHERE f.follower_id = $1
	          ORDER BY f.created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Roles, &u.AvatarURL, &u.Bio, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		u.CreatedAt = createdAt.Format(time.RFC3339)
		u.UpdatedAt = updatedAt.Format(time.RFC3339)
		users = append(users, u)
	}
	return users, nil
}

func (r *postgresFollowRepository) CountFollowers(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM tutora_app.follows WHERE following_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *postgresFollowRepository) CountFollowing(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM tutora_app.follows WHERE follower_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}
