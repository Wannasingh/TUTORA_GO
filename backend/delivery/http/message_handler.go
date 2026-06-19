package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/Wannasingh/TUTORA_GO/backend/delivery/ws"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
	"github.com/Wannasingh/TUTORA_GO/backend/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type MessageHandler struct {
	messageUsecase domain.MessageUsecase
	hub            *ws.Hub
}

func NewMessageHandler(mu domain.MessageUsecase, hub *ws.Hub) *MessageHandler {
	return &MessageHandler{messageUsecase: mu, hub: hub}
}

func (h *MessageHandler) StartConversation(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req domain.StartConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, err := h.messageUsecase.StartConversation(c.Request.Context(), userID.(int), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, conv)
}

func (h *MessageHandler) ListMyConversations(c *gin.Context) {
	userID, _ := c.Get("userID")
	convos, err := h.messageUsecase.ListMyConversations(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convos)
}

func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID, _ := c.Get("userID")
	convID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	conv, err := h.messageUsecase.GetConversation(c.Request.Context(), userID.(int), convID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conv)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	convID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.messageUsecase.GetConversationHistory(c.Request.Context(), userID.(int), convID, limit, offset)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	convID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req domain.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.messageUsecase.SendMessage(c.Request.Context(), userID.(int), convID, &req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Push via WebSocket to all conversation members
	conv, _ := h.messageUsecase.GetConversation(c.Request.Context(), userID.(int), convID)
	if conv != nil {
		wsMsg := &ws.WSMessage{Type: "message", Payload: msg}
		for _, member := range conv.Members {
			if member.UserID != userID.(int) {
				h.hub.SendToUser(member.UserID, wsMsg)
			}
		}
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	convID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	if err := h.messageUsecase.MarkAsRead(c.Request.Context(), userID.(int), convID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *MessageHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}
	userID := int(sub)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	ws.ServeWS(h.hub, userID, conn)
}

