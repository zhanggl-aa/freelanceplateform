package repository

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
)

type MessageRepository struct {
	db *DB
}

func NewMessageRepository(db *DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, conversationID, senderID, content, messageType string, fileURL, fileName *string, fileSize *int64) (*model.ChatMessage, error) {
	var m model.ChatMessage
	query := `
		INSERT INTO chat_messages (conversation_id, sender_id, content, message_type, file_url, file_name, file_size)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, conversation_id, sender_id, content, message_type, file_url, file_name, file_size, is_read, created_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		conversationID, senderID, nilIfEmpty(content), messageType, fileURL, fileName, fileSize,
	).Scan(
		&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.MessageType,
		&m.FileURL, &m.FileName, &m.FileSize, &m.IsRead, &m.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	return &m, nil
}

func (r *MessageRepository) ListByConversation(ctx context.Context, conversationID string, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM chat_messages WHERE conversation_id = $1`
	if err := r.db.Pool.QueryRow(ctx, countQuery, conversationID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count messages by conversation: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.content, m.message_type,
			m.file_url, m.file_name, m.file_size, m.is_read, m.created_at,
			u.nickname AS sender_name, u.avatar_url AS sender_avatar
		FROM chat_messages m
		LEFT JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, conversationID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list messages by conversation: %w", err)
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		var m model.ChatMessage
		if err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.MessageType,
			&m.FileURL, &m.FileName, &m.FileSize, &m.IsRead, &m.CreatedAt,
			&m.SenderName, &m.SenderAvatar,
		); err != nil {
			return nil, 0, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate messages: %w", err)
	}
	return messages, total, nil
}

func (r *MessageRepository) MarkAsRead(ctx context.Context, conversationID, userID string) error {
	query := `
		UPDATE chat_messages
		SET is_read = true
		WHERE conversation_id = $1
		  AND sender_id != $2
		  AND is_read = false
	`
	_, err := r.db.Pool.Exec(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("mark messages as read: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM chat_messages m
		INNER JOIN conversation_participants cp ON cp.conversation_id = m.conversation_id
		WHERE cp.user_id = $1
		  AND m.sender_id != $1
		  AND m.is_read = false
		  AND m.created_at > cp.last_read_at
	`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get unread count: %w", err)
	}
	return count, nil
}

func (r *MessageRepository) GetRecentMessages(ctx context.Context, conversationID string, limit int) ([]*model.ChatMessage, error) {
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.content, m.message_type,
			m.file_url, m.file_name, m.file_size, m.is_read, m.created_at,
			u.nickname AS sender_name, u.avatar_url AS sender_avatar
		FROM chat_messages m
		LEFT JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Pool.Query(ctx, query, conversationID, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		var m model.ChatMessage
		if err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.MessageType,
			&m.FileURL, &m.FileName, &m.FileSize, &m.IsRead, &m.CreatedAt,
			&m.SenderName, &m.SenderAvatar,
		); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent messages: %w", err)
	}
	return messages, nil
}
