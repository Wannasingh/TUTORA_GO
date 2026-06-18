package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/haru/bytestutor/backend/domain"
)

type postgresTutorRepository struct {
	db *pgx.Conn
}

func NewPostgresTutorRepository(db *pgx.Conn) domain.TutorRepository {
	return &postgresTutorRepository{db: db}
}

func (r *postgresTutorRepository) Create(ctx context.Context, tutor *domain.Tutor) error {
	query := `INSERT INTO tutora_app.tutors (user_id, subject, bio, price_per_hour) 
	          VALUES ($1, $2, $3, $4) RETURNING id, rating, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, tutor.UserID, tutor.Subject, tutor.Bio, tutor.PricePerHour).
		Scan(&tutor.ID, &tutor.Rating, &tutor.CreatedAt, &tutor.UpdatedAt)
	return err
}

func (r *postgresTutorRepository) GetByID(ctx context.Context, id int) (*domain.Tutor, error) {
	query := `SELECT t.id, t.user_id, t.subject, t.bio, t.price_per_hour, t.rating, t.created_at, t.updated_at,
	                 u.name, u.email, u.role
	          FROM tutora_app.tutors t
	          JOIN tutora_app.users u ON t.user_id = u.id
	          WHERE t.id = $1`
	tutor := &domain.Tutor{User: &domain.User{}}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&tutor.ID, &tutor.UserID, &tutor.Subject, &tutor.Bio, &tutor.PricePerHour, &tutor.Rating, &tutor.CreatedAt, &tutor.UpdatedAt,
			&tutor.User.Name, &tutor.User.Email, &tutor.User.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	tutor.User.ID = tutor.UserID
	return tutor, nil
}

func (r *postgresTutorRepository) List(ctx context.Context, subject string) ([]*domain.Tutor, error) {
	var rows pgx.Rows
	var err error

	if subject != "" {
		query := `SELECT t.id, t.user_id, t.subject, t.bio, t.price_per_hour, t.rating, t.created_at, t.updated_at,
		                 u.name, u.email, u.role
		          FROM tutora_app.tutors t
		          JOIN tutora_app.users u ON t.user_id = u.id
		          WHERE t.subject ILIKE $1`
		rows, err = r.db.Query(ctx, query, "%"+subject+"%")
	} else {
		query := `SELECT t.id, t.user_id, t.subject, t.bio, t.price_per_hour, t.rating, t.created_at, t.updated_at,
		                 u.name, u.email, u.role
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
		err := rows.Scan(&t.ID, &t.UserID, &t.Subject, &t.Bio, &t.PricePerHour, &t.Rating, &t.CreatedAt, &t.UpdatedAt,
			&t.User.Name, &t.User.Email, &t.User.Role)
		if err != nil {
			return nil, err
		}
		t.User.ID = t.UserID
		tutors = append(tutors, t)
	}

	return tutors, nil
}
