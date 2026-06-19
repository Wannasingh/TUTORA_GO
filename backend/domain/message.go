package domain

import "context"

type Conversation struct {
	ID        int                   `json:"id"`
	Type      string                `json:"type"`
	Title     *string               `json:"title,omitempty"`
	Members   []*ConversationMember `json:"members,omitempty"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

type ConversationMember struct {
	UserID     int     `json:"user_id"`
	UserName   string  `json:"user_name"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	LastReadAt string  `json:"last_read_at"`
	JoinedAt   string  `json:"joined_at"`
}

type Message struct {
	ID             int     `json:"id"`
	ConversationID int     `json:"conversation_id"`
	SenderID       int     `json:"sender_id"`
	SenderName     string  `json:"sender_name"`
	SenderAvatar   *string `json:"sender_avatar,omitempty"`
	Body           *string `json:"body,omitempty"`
	ImageURL       *string `json:"image_url,omitempty"`
	MessageType    string  `json:"message_type"`
	CreatedAt      string  `json:"created_at"`
}

type ConversationPreview struct {
	ID          int      `json:"id"`
	Type        string   `json:"type"`
	Title       *string  `json:"title,omitempty"`
	LastMessage *Message `json:"last_message,omitempty"`
	UnreadCount int      `json:"unread_count"`
	OtherUser   *User    `json:"other_user,omitempty"`
	UpdatedAt   string   `json:"updated_at"`
}

type StartConversationRequest struct {
	Type      string  `json:"type" binding:"required"`
	Title     *string `json:"title,omitempty"`
	MemberIDs []int   `json:"member_ids" binding:"required"`
}

type SendMessageRequest struct {
	Body        *string `json:"body,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	MessageType string  `json:"message_type"`
}

type MessageRepository interface {
	CreateConversation(ctx context.Context, conv *Conversation) error
	AddMembers(ctx context.Context, conversationID int, userIDs []int) error
	GetConversationByID(ctx context.Context, id int) (*Conversation, error)
	ListConversations(ctx context.Context, userID int) ([]*ConversationPreview, error)
	IsMember(ctx context.Context, conversationID, userID int) (bool, error)
	FindDirectConversation(ctx context.Context, userID1, userID2 int) (*Conversation, error)
	SendMessage(ctx context.Context, msg *Message) error
	GetMessages(ctx context.Context, conversationID int, limit, offset int) ([]*Message, error)
	MarkRead(ctx context.Context, conversationID, userID int) error
	UpdateConversationTimestamp(ctx context.Context, conversationID int) error
}

type MessageUsecase interface {
	StartConversation(ctx context.Context, userID int, req *StartConversationRequest) (*Conversation, error)
	SendMessage(ctx context.Context, userID, conversationID int, req *SendMessageRequest) (*Message, error)
	GetConversationHistory(ctx context.Context, userID, conversationID, limit, offset int) ([]*Message, error)
	ListMyConversations(ctx context.Context, userID int) ([]*ConversationPreview, error)
	GetConversation(ctx context.Context, userID, conversationID int) (*Conversation, error)
	MarkAsRead(ctx context.Context, userID, conversationID int) error
}
