package usecase

import (
	"context"
	"fmt"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type notificationUsecase struct {
	repo domain.NotificationRepository
}

func NewNotificationUsecase(repo domain.NotificationRepository) domain.NotificationUsecase {
	return &notificationUsecase{repo: repo}
}

func (u *notificationUsecase) Notify(ctx context.Context, userID int, notifType, title string, body *string, dataJSON *string) error {
	notif := &domain.Notification{
		UserID:   userID,
		Type:     notifType,
		Title:    title,
		Body:     body,
		DataJSON: dataJSON,
	}
	return u.repo.Create(ctx, notif)
}

func (u *notificationUsecase) ListMyNotifications(ctx context.Context, userID int, limit, offset int) ([]*domain.Notification, error) {
	return u.repo.ListByUser(ctx, userID, limit, offset)
}

func (u *notificationUsecase) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	return u.repo.CountUnread(ctx, userID)
}

func (u *notificationUsecase) MarkRead(ctx context.Context, userID, notifID int) error {
	notif, err := u.repo.GetByID(ctx, notifID)
	if err != nil {
		return err
	}
	if notif == nil {
		return fmt.Errorf("notification not found")
	}
	if notif.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.MarkRead(ctx, notifID)
}

func (u *notificationUsecase) MarkAllRead(ctx context.Context, userID int) error {
	return u.repo.MarkAllRead(ctx, userID)
}

func (u *notificationUsecase) DeleteNotification(ctx context.Context, userID, notifID int) error {
	notif, err := u.repo.GetByID(ctx, notifID)
	if err != nil {
		return err
	}
	if notif == nil {
		return fmt.Errorf("notification not found")
	}
	if notif.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.Delete(ctx, notifID)
}
