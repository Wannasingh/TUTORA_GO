package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresTutorRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTutorRepository(db *pgxpool.Pool) domain.TutorRepository {
	return &postgresTutorRepository{db: db}
}

func (r *postgresTutorRepository) Create(ctx context.Context, tutor *domain.Tutor) error {
	query := `INSERT INTO tutora_app.tutors (user_id, subject, bio, price_per_hour) 
	          VALUES ($1, $2, $3, $4) RETURNING id, rating, created_at, updated_at`
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, tutor.UserID, tutor.Subject, tutor.Bio, tutor.PricePerHour).
		Scan(&tutor.ID, &tutor.Rating, &createdAt, &updatedAt)
	if err == nil {
		tutor.CreatedAt = createdAt.Format(time.RFC3339)
		tutor.UpdatedAt = updatedAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresTutorRepository) GetByID(ctx context.Context, id int) (*domain.Tutor, error) {
	query := `SELECT t.id, t.user_id, t.subject, t.bio, t.price_per_hour, t.rating, t.created_at, t.updated_at,
	                 u.name, u.email, u.roles
	          FROM tutora_app.tutors t
	          JOIN tutora_app.users u ON t.user_id = u.id
	          WHERE t.id = $1`
	tutor := &domain.Tutor{User: &domain.User{}}
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id).
		Scan(&tutor.ID, &tutor.UserID, &tutor.Subject, &tutor.Bio, &tutor.PricePerHour, &tutor.Rating, &createdAt, &updatedAt,
			&tutor.User.Name, &tutor.User.Email, &tutor.User.Roles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	tutor.User.ID = tutor.UserID
	tutor.CreatedAt = createdAt.Format(time.RFC3339)
	tutor.UpdatedAt = updatedAt.Format(time.RFC3339)
	return tutor, nil
}

func (r *postgresTutorRepository) List(ctx context.Context, subject string) ([]*domain.Tutor, error) {
	var rows pgx.Rows
	var err error

	if subject != "" {
		query := `SELECT t.id, t.user_id, t.subject, t.bio, t.price_per_hour, t.rating, t.created_at, t.updated_at,
		                 u.name, u.email, u.roles
		          FROM tutora_app.tutors t
		          JOIN tutora_app.users u ON t.user_id = u.id
		          WHERE t.subject ILIKE $1`
		rows, err = r.db.Query(ctx, query, "%"+subject+"%")
	} else {
		query := `SELECT t.id, t.user_id, t.subject, t.bio, t.price_per_hour, t.rating, t.created_at, t.updated_at,
		                 u.name, u.email, u.roles
		          FROM tutora_app.tutors t
		          JOIN tutora_app.users u ON t.user_id = u.id`
		rows, err = r.db.Query(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tutors []*domain.Tutor
	for rows.Next() {
		t := &domain.Tutor{User: &domain.User{}}
		var createdAt, updatedAt time.Time
		err := rows.Scan(&t.ID, &t.UserID, &t.Subject, &t.Bio, &t.PricePerHour, &t.Rating, &createdAt, &updatedAt,
			&t.User.Name, &t.User.Email, &t.User.Roles)
		if err != nil {
			return nil, err
		}
		t.User.ID = t.UserID
		t.CreatedAt = createdAt.Format(time.RFC3339)
		t.UpdatedAt = updatedAt.Format(time.RFC3339)
		tutors = append(tutors, t)
	}

	return tutors, nil
}
