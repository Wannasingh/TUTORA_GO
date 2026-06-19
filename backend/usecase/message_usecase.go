package usecase

import (
	"context"
	"fmt"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type messageUsecase struct {
	repo domain.MessageRepository
}

func NewMessageUsecase(repo domain.MessageRepository) domain.MessageUsecase {
	return &messageUsecase{repo: repo}
}

func (u *messageUsecase) StartConversation(ctx context.Context, userID int, req *domain.StartConversationRequest) (*domain.Conversation, error) {
	// For direct conversations, check if one already exists
	if req.Type == "direct" {
		if len(req.MemberIDs) != 1 {
			return nil, fmt.Errorf("direct conversation requires exactly one other member")
		}
		otherUserID := req.MemberIDs[0]
		if otherUserID == userID {
			return nil, fmt.Errorf("cannot start conversation with yourself")
		}
		existing, err := u.repo.FindDirectConversation(ctx, userID, otherUserID)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	conv := &domain.Conversation{
		Type:  req.Type,
		Title: req.Title,
	}
	if err := u.repo.CreateConversation(ctx, conv); err != nil {
		return nil, err
	}

	// Add all members including the creator
	allMembers := append([]int{userID}, req.MemberIDs...)
	if err := u.repo.AddMembers(ctx, conv.ID, allMembers); err != nil {
		return nil, err
	}

	return u.repo.GetConversationByID(ctx, conv.ID)
}

func (u *messageUsecase) SendMessage(ctx context.Context, userID, conversationID int, req *domain.SendMessageRequest) (*domain.Message, error) {
	isMember, err := u.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("not a member of this conversation")
	}

	msgType := req.MessageType
	if msgType == "" {
		msgType = "text"
	}

	msg := &domain.Message{
		ConversationID: conversationID,
		SenderID:       userID,
		Body:           req.Body,
		ImageURL:       req.ImageURL,
		MessageType:    msgType,
	}

	if err := u.repo.SendMessage(ctx, msg); err != nil {
		return nil, err
	}

	_ = u.repo.UpdateConversationTimestamp(ctx, conversationID)
	return msg, nil
}

func (u *messageUsecase) GetConversationHistory(ctx context.Context, userID, conversationID, limit, offset int) ([]*domain.Message, error) {
	isMember, err := u.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("not a member of this conversation")
	}
	return u.repo.GetMessages(ctx, conversationID, limit, offset)
}

func (u *messageUsecase) ListMyConversations(ctx context.Context, userID int) ([]*domain.ConversationPreview, error) {
	return u.repo.ListConversations(ctx, userID)
}

func (u *messageUsecase) GetConversation(ctx context.Context, userID, conversationID int) (*domain.Conversation, error) {
	isMember, err := u.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("not a member of this conversation")
	}
	return u.repo.GetConversationByID(ctx, conversationID)
}

func (u *messageUsecase) MarkAsRead(ctx context.Context, userID, conversationID int) error {
	isMember, err := u.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return fmt.Errorf("not a member of this conversation")
	}
	return u.repo.MarkRead(ctx, conversationID, userID)
}
