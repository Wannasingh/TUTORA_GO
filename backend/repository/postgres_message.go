package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresMessageRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMessageRepository(db *pgxpool.Pool) domain.MessageRepository {
	return &postgresMessageRepository{db: db}
}

func (r *postgresMessageRepository) CreateConversation(ctx context.Context, conv *domain.Conversation) error {
	query := `INSERT INTO tutora_app.conversations (type, title) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, conv.Type, conv.Title).Scan(&conv.ID, &createdAt, &updatedAt)
	if err == nil {
		conv.CreatedAt = createdAt.Format(time.RFC3339)
		conv.UpdatedAt = updatedAt.Format(time.RFC3339)
	}
	return err
}

func (r *postgresMessageRepository) AddMembers(ctx context.Context, conversationID int, userIDs []int) error {
	for _, uid := range userIDs {
		_, err := r.db.Exec(ctx,
			`INSERT INTO tutora_app.conversation_members (conversation_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			conversationID, uid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *postgresMessageRepository) GetConversationByID(ctx context.Context, id int) (*domain.Conversation, error) {
	query := `SELECT id, type, title, created_at, updated_at FROM tutora_app.conversations WHERE id = $1`
	conv := &domain.Conversation{}
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&conv.ID, &conv.Type, &conv.Title, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	conv.CreatedAt = createdAt.Format(time.RFC3339)
	conv.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Load members
	membersQuery := `SELECT cm.user_id, u.name, u.avatar_url, cm.last_read_at, cm.joined_at
	                 FROM tutora_app.conversation_members cm
	                 JOIN tutora_app.users u ON u.id = cm.user_id
	                 WHERE cm.conversation_id = $1`
	rows, err := r.db.Query(ctx, membersQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := &domain.ConversationMember{}
		var lastReadAt, joinedAt time.Time
		if err := rows.Scan(&m.UserID, &m.UserName, &m.AvatarURL, &lastReadAt, &joinedAt); err != nil {
			return nil, err
		}
		m.LastReadAt = lastReadAt.Format(time.RFC3339)
		m.JoinedAt = joinedAt.Format(time.RFC3339)
		conv.Members = append(conv.Members, m)
	}
	return conv, nil
}

func (r *postgresMessageRepository) ListConversations(ctx context.Context, userID int) ([]*domain.ConversationPreview, error) {
	query := `SELECT c.id, c.type, c.title, c.updated_at,
	                 (SELECT COUNT(*) FROM tutora_app.messages m
	                  WHERE m.conversation_id = c.id
	                    AND m.created_at > cm.last_read_at) as unread_count
	          FROM tutora_app.conversations c
	          JOIN tutora_app.conversation_members cm ON cm.conversation_id = c.id AND cm.user_id = $1
	          ORDER BY c.updated_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var previews []*domain.ConversationPreview
	for rows.Next() {
		p := &domain.ConversationPreview{}
		var updatedAt time.Time
		if err := rows.Scan(&p.ID, &p.Type, &p.Title, &updatedAt, &p.UnreadCount); err != nil {
			return nil, err
		}
		p.UpdatedAt = updatedAt.Format(time.RFC3339)

		// Get last message
		lastMsgQuery := `SELECT m.id, m.conversation_id, m.sender_id, u.name, u.avatar_url, m.body, m.image_url, m.message_type, m.created_at
		                 FROM tutora_app.messages m
		                 JOIN tutora_app.users u ON u.id = m.sender_id
		                 WHERE m.conversation_id = $1
		                 ORDER BY m.created_at DESC LIMIT 1`
		msg := &domain.Message{}
		var msgCreatedAt time.Time
		err := r.db.QueryRow(ctx, lastMsgQuery, p.ID).
			Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.SenderName, &msg.SenderAvatar,
				&msg.Body, &msg.ImageURL, &msg.MessageType, &msgCreatedAt)
		if err == nil {
			msg.CreatedAt = msgCreatedAt.Format(time.RFC3339)
			p.LastMessage = msg
		}

		// For direct convos, get the other user
		if p.Type == "direct" {
			otherQuery := `SELECT u.id, u.name, u.avatar_url
			               FROM tutora_app.conversation_members cm
			               JOIN tutora_app.users u ON u.id = cm.user_id
			               WHERE cm.conversation_id = $1 AND cm.user_id != $2 LIMIT 1`
			otherUser := &domain.User{}
			_ = r.db.QueryRow(ctx, otherQuery, p.ID, userID).Scan(&otherUser.ID, &otherUser.Name, &otherUser.AvatarURL)
			if otherUser.ID > 0 {
				p.OtherUser = otherUser
			}
		}
		previews = append(previews, p)
	}
	return previews, nil
}

func (r *postgresMessageRepository) IsMember(ctx context.Context, conversationID, userID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tutora_app.conversation_members WHERE conversation_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, conversationID, userID).Scan(&exists)
	return exists, err
}

func (r *postgresMessageRepository) FindDirectConversation(ctx context.Context, userID1, userID2 int) (*domain.Conversation, error) {
	query := `SELECT c.id FROM tutora_app.conversations c
	          WHERE c.type = 'direct'
	            AND EXISTS (SELECT 1 FROM tutora_app.conversation_members WHERE conversation_id = c.id AND user_id = $1)
	            AND EXISTS (SELECT 1 FROM tutora_app.conversation_members WHERE conversation_id = c.id AND user_id = $2)
	          LIMIT 1`
	var convID int
	err := r.db.QueryRow(ctx, query, userID1, userID2).Scan(&convID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.GetConversationByID(ctx, convID)
}

func (r *postgresMessageRepository) SendMessage(ctx context.Context, msg *domain.Message) error {
	query := `INSERT INTO tutora_app.messages (conversation_id, sender_id, body, image_url, message_type)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	var createdAt time.Time
	err := r.db.QueryRow(ctx, query, msg.ConversationID, msg.SenderID, msg.Body, msg.ImageURL, msg.MessageType).
		Scan(&msg.ID, &createdAt)
	if err != nil {
		return err
	}
	msg.CreatedAt = createdAt.Format(time.RFC3339)

	// Fetch sender info
	nameQuery := `SELECT name, avatar_url FROM tutora_app.users WHERE id = $1`
	_ = r.db.QueryRow(ctx, nameQuery, msg.SenderID).Scan(&msg.SenderName, &msg.SenderAvatar)
	return nil
}

func (r *postgresMessageRepository) GetMessages(ctx context.Context, conversationID int, limit, offset int) ([]*domain.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	query := fmt.Sprintf(`SELECT m.id, m.conversation_id, m.sender_id, u.name, u.avatar_url, m.body, m.image_url, m.message_type, m.created_at
	          FROM tutora_app.messages m
	          JOIN tutora_app.users u ON u.id = m.sender_id
	          WHERE m.conversation_id = $1
	          ORDER BY m.created_at DESC
	          LIMIT %d OFFSET %d`, limit, offset)
	rows, err := r.db.Query(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		msg := &domain.Message{}
		var createdAt time.Time
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.SenderName, &msg.SenderAvatar,
			&msg.Body, &msg.ImageURL, &msg.MessageType, &createdAt); err != nil {
			return nil, err
		}
		msg.CreatedAt = createdAt.Format(time.RFC3339)
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *postgresMessageRepository) MarkRead(ctx context.Context, conversationID, userID int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE tutora_app.conversation_members SET last_read_at = NOW() WHERE conversation_id = $1 AND user_id = $2`,
		conversationID, userID)
	return err
}

func (r *postgresMessageRepository) UpdateConversationTimestamp(ctx context.Context, conversationID int) error {
	_, err := r.db.Exec(ctx, `UPDATE tutora_app.conversations SET updated_at = NOW() WHERE id = $1`, conversationID)
	return err
}
