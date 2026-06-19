package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type NotificationHandler struct {
	notifUsecase domain.NotificationUsecase
}

func NewNotificationHandler(nu domain.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{notifUsecase: nu}
}

func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, _ := c.Get("userID")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifs, err := h.notifUsecase.ListMyNotifications(c.Request.Context(), userID.(int), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifs)
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("userID")
	count, err := h.notifUsecase.GetUnreadCount(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	notifID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.notifUsecase.MarkRead(c.Request.Context(), userID.(int), notifID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	if err := h.notifUsecase.MarkAllRead(c.Request.Context(), userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, _ := c.Get("userID")
	notifID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.notifUsecase.DeleteNotification(c.Request.Context(), userID.(int), notifID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification deleted"})
}
