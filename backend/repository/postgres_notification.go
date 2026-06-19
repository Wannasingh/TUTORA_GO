package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresNotificationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresNotificationRepository(db *pgxpool.Pool) domain.NotificationRepository {
	return &postgresNotificationRepository{db: db}
}

func (r *postgresNotificationRepository) Create(ctx context.Context, notif *domain.Notification) error {
	query := `INSERT INTO tutora_app.notifications (user_id, type, title, body, data_json)
	          VALUES ($1, $2, $3, $4, $5::jsonb) RETURNING id, created_at`
	var createdAt time.Time
	err := r.db.QueryRow(ctx, query, notif.UserID, notif.Type, notif.Title, notif.Body, notif.DataJSON).
		Scan(&notif.ID, &createdAt)
	if err == nil {
		notif.CreatedAt = createdAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresNotificationRepository) ListByUser(ctx context.Context, userID int, limit, offset int) ([]*domain.Notification, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `SELECT id, user_id, type, title, body, data_json::text, is_read, created_at
	          FROM tutora_app.notifications
	          WHERE user_id = $1
	          ORDER BY created_at DESC
	          LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		var createdAt time.Time
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.DataJSON, &n.IsRead, &createdAt); err != nil {
			return nil, err
		}
		n.CreatedAt = createdAt.Format(time.RFC3339)
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *postgresNotificationRepository) CountUnread(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM tutora_app.notifications WHERE user_id = $1 AND is_read = FALSE`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *postgresNotificationRepository) MarkRead(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `UPDATE tutora_app.notifications SET is_read = TRUE WHERE id = $1`, id)
	return err
}

func (r *postgresNotificationRepository) MarkAllRead(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, `UPDATE tutora_app.notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`, userID)
	return err
}

func (r *postgresNotificationRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tutora_app.notifications WHERE id = $1`, id)
	return err
}

func (r *postgresNotificationRepository) GetByID(ctx context.Context, id int) (*domain.Notification, error) {
	query := `SELECT id, user_id, type, title, body, data_json::text, is_read, created_at
	          FROM tutora_app.notifications WHERE id = $1`
	n := &domain.Notification{}
	var createdAt time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.DataJSON, &n.IsRead, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	n.CreatedAt = createdAt.Format(time.RFC3339)
	return n, nil
}
