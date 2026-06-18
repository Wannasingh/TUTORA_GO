package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresUserRepository struct {
	db *pgx.Conn
}

func NewPostgresUserRepository(db *pgx.Conn) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO tutora_app.users (name, email, role, password_hash, google_id, apple_id) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.Role, user.PasswordHash, user.GoogleID, user.AppleID).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, name, email, role, password_hash, google_id, apple_id, created_at, updated_at 
	          FROM tutora_app.users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.PasswordHash, &user.GoogleID, &user.AppleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, role, password_hash, google_id, apple_id, created_at, updated_at 
	          FROM tutora_app.users WHERE email = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.PasswordHash, &user.GoogleID, &user.AppleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	query := `SELECT id, name, email, role, password_hash, google_id, apple_id, created_at, updated_at 
	          FROM tutora_app.users WHERE google_id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, googleID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.PasswordHash, &user.GoogleID, &user.AppleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByAppleID(ctx context.Context, appleID string) (*domain.User, error) {
	query := `SELECT id, name, email, role, password_hash, google_id, apple_id, created_at, updated_at 
	          FROM tutora_app.users WHERE apple_id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, appleID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.PasswordHash, &user.GoogleID, &user.AppleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tutora_app.users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
