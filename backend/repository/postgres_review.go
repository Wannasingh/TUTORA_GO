package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresReviewRepository struct {
	db *pgxpool.Pool
}

func NewPostgresReviewRepository(db *pgxpool.Pool) domain.ReviewRepository {
	return &postgresReviewRepository{db: db}
}

func (r *postgresReviewRepository) CreateOrUpdate(ctx context.Context, review *domain.TutorReview) error {
	query := `INSERT INTO tutora_app.tutor_reviews (reviewer_id, tutor_id, rating, body)
	          VALUES ($1, $2, $3, $4)
	          ON CONFLICT (reviewer_id, tutor_id) DO UPDATE SET rating = $3, body = $4
	          RETURNING id, created_at`
	var createdAt time.Time
	err := r.db.QueryRow(ctx, query, review.ReviewerID, review.TutorID, review.Rating, review.Body).
		Scan(&review.ID, &createdAt)
	if err == nil {
		review.CreatedAt = createdAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresReviewRepository) GetByTutorID(ctx context.Context, tutorID int) ([]*domain.TutorReview, error) {
	query := `SELECT tr.id, tr.reviewer_id, tr.tutor_id, tr.rating, tr.body, tr.created_at,
	                 u.name, u.avatar_url
	          FROM tutora_app.tutor_reviews tr
	          JOIN tutora_app.users u ON u.id = tr.reviewer_id
	          WHERE tr.tutor_id = $1
	          ORDER BY tr.created_at DESC`
	rows, err := r.db.Query(ctx, query, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.TutorReview
	for rows.Next() {
		rv := &domain.TutorReview{}
		var createdAt time.Time
		if err := rows.Scan(&rv.ID, &rv.ReviewerID, &rv.TutorID, &rv.Rating, &rv.Body, &createdAt,
			&rv.ReviewerName, &rv.ReviewerAvatar); err != nil {
			return nil, err
		}
		rv.CreatedAt = createdAt.Format(time.RFC3339)
		reviews = append(reviews, rv)
	}
	return reviews, nil
}

func (r *postgresReviewRepository) GetByID(ctx context.Context, id int) (*domain.TutorReview, error) {
	query := `SELECT id, reviewer_id, tutor_id, rating, body, created_at FROM tutora_app.tutor_reviews WHERE id = $1`
	rv := &domain.TutorReview{}
	var createdAt time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&rv.ID, &rv.ReviewerID, &rv.TutorID, &rv.Rating, &rv.Body, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	rv.CreatedAt = createdAt.Format(time.RFC3339)
	return rv, nil
}

func (r *postgresReviewRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tutora_app.tutor_reviews WHERE id = $1`, id)
	return err
}

func (r *postgresReviewRepository) GetAverageRating(ctx context.Context, tutorID int) (float64, error) {
	query := `SELECT COALESCE(AVG(rating), 0) FROM tutora_app.tutor_reviews WHERE tutor_id = $1`
	var avg float64
	err := r.db.QueryRow(ctx, query, tutorID).Scan(&avg)
	return avg, err
}

func (r *postgresReviewRepository) UpdateTutorRating(ctx context.Context, tutorID int, rating float64) error {
	_, err := r.db.Exec(ctx, `UPDATE tutora_app.tutors SET rating = $2, updated_at = NOW() WHERE id = $1`, tutorID, rating)
	return err
}
