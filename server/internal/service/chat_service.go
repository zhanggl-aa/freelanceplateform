package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type ChatService struct {
	conversationRepo *repository.ConversationRepository
	messageRepo      *repository.MessageRepository
}

func NewChatService(conversationRepo *repository.ConversationRepository, messageRepo *repository.MessageRepository) *ChatService {
	return &ChatService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}

func (s *ChatService) CreateConversation(ctx context.Context, convType string, projectID *string, participantIDs []string) (*model.ChatConversation, error) {
	if len(participantIDs) < 2 {
		return nil, errors.New("at least 2 participants are required")
	}
	if convType == "" {
		convType = "direct"
	}
	validTypes := map[string]bool{"direct": true, "group": true, "project": true}
	if !validTypes[convType] {
		return nil, errors.New("conversation type must be 'direct', 'group', or 'project'")
	}

	created, err := s.conversationRepo.Create(ctx, convType, projectID)
	if err != nil {
		return nil, fmt.Errorf("create conversation: %w", err)
	}

	for _, pid := range participantIDs {
		if err := s.conversationRepo.AddParticipant(ctx, created.ID, pid); err != nil {
			return nil, fmt.Errorf("add participant: %w", err)
		}
	}

	return created, nil
}

func (s *ChatService) GetOrCreateDirectConversation(ctx context.Context, userID1, userID2 string) (*model.ChatConversation, error) {
	if userID1 == "" || userID2 == "" {
		return nil, errors.New("both user ids are required")
	}
	if userID1 == userID2 {
		return nil, errors.New("cannot create a conversation with yourself")
	}

	conv, err := s.conversationRepo.FindDirect(ctx, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("find direct conversation: %w", err)
	}
	if conv != nil {
		return conv, nil
	}

	created, err := s.conversationRepo.Create(ctx, "direct", nil)
	if err != nil {
		return nil, fmt.Errorf("create direct conversation: %w", err)
	}

	if err := s.conversationRepo.AddParticipant(ctx, created.ID, userID1); err != nil {
		return nil, fmt.Errorf("add participant 1: %w", err)
	}
	if err := s.conversationRepo.AddParticipant(ctx, created.ID, userID2); err != nil {
		return nil, fmt.Errorf("add participant 2: %w", err)
	}

	return created, nil
}

func (s *ChatService) GetByID(ctx context.Context, id string) (*model.ChatConversation, error) {
	if id == "" {
		return nil, errors.New("conversation id is required")
	}
	conv, err := s.conversationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get conversation: %w", err)
	}
	if conv == nil {
		return nil, errors.New("conversation not found")
	}
	return conv, nil
}

func (s *ChatService) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.ChatConversation, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.conversationRepo.ListByUser(ctx, userID, page, pageSize)
}

func (s *ChatService) SendMessage(ctx context.Context, conversationID, senderID, content, messageType string, fileURL, fileName *string, fileSize *int64) (*model.ChatMessage, error) {
	if conversationID == "" {
		return nil, errors.New("conversation_id is required")
	}
	if senderID == "" {
		return nil, errors.New("sender_id is required")
	}

	hasContent := content != ""
	hasFile := fileURL != nil && *fileURL != ""
	if !hasContent && !hasFile {
		return nil, errors.New("message must have content or a file attachment")
	}

	if messageType == "" {
		if hasFile {
			messageType = "file"
		} else {
			messageType = "text"
		}
	}

	validMessageTypes := map[string]bool{"text": true, "image": true, "file": true, "system": true}
	if !validMessageTypes[messageType] {
		return nil, errors.New("message_type must be 'text', 'image', 'file', or 'system'")
	}

	created, err := s.messageRepo.Create(ctx, conversationID, senderID, content, messageType, fileURL, fileName, fileSize)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	_ = s.conversationRepo.UpdateLastMessageAt(ctx, conversationID)

	return created, nil
}

func (s *ChatService) ListMessages(ctx context.Context, conversationID string, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	if conversationID == "" {
		return nil, 0, errors.New("conversation_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.messageRepo.ListByConversation(ctx, conversationID, page, pageSize)
}

func (s *ChatService) MarkAsRead(ctx context.Context, conversationID, userID string) error {
	if conversationID == "" {
		return errors.New("conversation_id is required")
	}
	if userID == "" {
		return errors.New("user_id is required")
	}
	if err := s.messageRepo.MarkAsRead(ctx, conversationID, userID); err != nil {
		return fmt.Errorf("mark as read: %w", err)
	}
	return s.conversationRepo.UpdateLastReadAt(ctx, conversationID, userID)
}

func (s *ChatService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, errors.New("user_id is required")
	}
	return s.messageRepo.GetUnreadCount(ctx, userID)
}
