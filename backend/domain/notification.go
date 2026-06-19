package domain

import "context"

type Notification struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Type      string  `json:"type"` // 'message', 'remind', 'promotion', 'system', 'critical'
	Title     string  `json:"title"`
	Body      *string `json:"body,omitempty"`
	DataJSON  *string `json:"data,omitempty"`
	IsRead    bool    `json:"is_read"`
	CreatedAt string  `json:"created_at"`
}

type NotificationRepository interface {
	Create(ctx context.Context, notif *Notification) error
	ListByUser(ctx context.Context, userID int, limit, offset int) ([]*Notification, error)
	CountUnread(ctx context.Context, userID int) (int, error)
	MarkRead(ctx context.Context, id int) error
	MarkAllRead(ctx context.Context, userID int) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*Notification, error)
}

type NotificationUsecase interface {
	Notify(ctx context.Context, userID int, notifType, title string, body *string, dataJSON *string) error
	ListMyNotifications(ctx context.Context, userID int, limit, offset int) ([]*Notification, error)
	GetUnreadCount(ctx context.Context, userID int) (int, error)
	MarkRead(ctx context.Context, userID, notifID int) error
	MarkAllRead(ctx context.Context, userID int) error
	DeleteNotification(ctx context.Context, userID, notifID int) error
}
