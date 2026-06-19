package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO tutora_app.users (name, email, roles, password_hash, google_id, apple_id, avatar_url) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.Roles, user.PasswordHash, user.GoogleID, user.AppleID, user.AvatarURL).
		Scan(&user.ID, &createdAt, &updatedAt)
	if err == nil {
		user.CreatedAt = createdAt.Format(time.RFC3339)
		user.UpdatedAt = updatedAt.Format(time.RFC3339)
	}
	return err
}

const userSelectCols = `id, name, email, roles, password_hash, google_id, apple_id, avatar_url, bio, cover_url, phone, school, birthdate::text, created_at, updated_at`

func scanUser(row pgx.Row) (*domain.User, error) {
	user := &domain.User{}
	var createdAt, updatedAt time.Time
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Roles, &user.PasswordHash,
		&user.GoogleID, &user.AppleID, &user.AvatarURL, &user.Bio, &user.CoverURL,
		&user.Phone, &user.School, &user.Birthdate, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)
	return user, nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT ` + userSelectCols + ` FROM tutora_app.users WHERE id = $1`
	return scanUser(r.db.QueryRow(ctx, query, id))
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT ` + userSelectCols + ` FROM tutora_app.users WHERE email = $1`
	return scanUser(r.db.QueryRow(ctx, query, email))
}

func (r *postgresUserRepository) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	query := `SELECT ` + userSelectCols + ` FROM tutora_app.users WHERE google_id = $1`
	return scanUser(r.db.QueryRow(ctx, query, googleID))
}

func (r *postgresUserRepository) GetByAppleID(ctx context.Context, appleID string) (*domain.User, error) {
	query := `SELECT ` + userSelectCols + ` FROM tutora_app.users WHERE apple_id = $1`
	return scanUser(r.db.QueryRow(ctx, query, appleID))
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE tutora_app.users SET name=$2, avatar_url=$3, bio=$4, cover_url=$5, phone=$6, school=$7, birthdate=$8::date, updated_at=NOW()
	          WHERE id=$1 RETURNING updated_at`
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, user.ID, user.Name, user.AvatarURL, user.Bio, user.CoverURL, user.Phone, user.School, user.Birthdate).
		Scan(&updatedAt)
	if err == nil {
		user.UpdatedAt = updatedAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tutora_app.users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *postgresUserRepository) CountUserPosts(ctx context.Context, userID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM tutora_app.posts WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *postgresUserRepository) CountUserReposts(ctx context.Context, userID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM tutora_app.reposts WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}
