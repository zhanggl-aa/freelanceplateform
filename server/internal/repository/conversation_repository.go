package repository

import (
	"context"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type ConversationRepository struct {
	db *DB
}

func NewConversationRepository(db *DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(ctx context.Context, convType string, projectID *string) (*model.ChatConversation, error) {
	var c model.ChatConversation
	query := `
		INSERT INTO chat_conversations (type, project_id)
		VALUES ($1, $2)
		RETURNING id, type, project_id, last_message_at, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, convType, projectID).Scan(
		&c.ID, &c.Type, &c.ProjectID, &c.LastMessageAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create conversation: %w", err)
	}
	return &c, nil
}

func (r *ConversationRepository) GetByID(ctx context.Context, id string) (*model.ChatConversation, error) {
	var c model.ChatConversation
	query := `
		SELECT id, type, project_id, last_message_at, created_at, updated_at
		FROM chat_conversations
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Type, &c.ProjectID, &c.LastMessageAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get conversation by id: %w", err)
	}
	return &c, nil
}

func (r *ConversationRepository) FindDirect(ctx context.Context, userID1, userID2 string) (*model.ChatConversation, error) {
	var c model.ChatConversation
	query := `
		SELECT c.id, c.type, c.project_id, c.last_message_at, c.created_at, c.updated_at
		FROM chat_conversations c
		WHERE c.type = 'direct'
		  AND EXISTS (SELECT 1 FROM conversation_participants cp WHERE cp.conversation_id = c.id AND cp.user_id = $1)
		  AND EXISTS (SELECT 1 FROM conversation_participants cp WHERE cp.conversation_id = c.id AND cp.user_id = $2)
		LIMIT 1
	`
	err := r.db.Pool.QueryRow(ctx, query, userID1, userID2).Scan(
		&c.ID, &c.Type, &c.ProjectID, &c.LastMessageAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find direct conversation: %w", err)
	}
	return &c, nil
}

func (r *ConversationRepository) FindByProject(ctx context.Context, projectID string) (*model.ChatConversation, error) {
	var c model.ChatConversation
	query := `
		SELECT id, type, project_id, last_message_at, created_at, updated_at
		FROM chat_conversations
		WHERE project_id = $1
		LIMIT 1
	`
	err := r.db.Pool.QueryRow(ctx, query, projectID).Scan(
		&c.ID, &c.Type, &c.ProjectID, &c.LastMessageAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find conversation by project: %w", err)
	}
	return &c, nil
}

func (r *ConversationRepository) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.ChatConversation, int64, error) {
	var total int64
	countQuery := `
		SELECT COUNT(*)
		FROM chat_conversations c
		INNER JOIN conversation_participants cp ON cp.conversation_id = c.id
		WHERE cp.user_id = $1
	`
	if err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count conversations by user: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT c.id, c.type, c.project_id, c.last_message_at, c.created_at, c.updated_at
		FROM chat_conversations c
		INNER JOIN conversation_participants cp ON cp.conversation_id = c.id
		WHERE cp.user_id = $1
		ORDER BY c.last_message_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list conversations by user: %w", err)
	}
	defer rows.Close()

	var conversations []*model.ChatConversation
	for rows.Next() {
		var c model.ChatConversation
		if err := rows.Scan(
			&c.ID, &c.Type, &c.ProjectID, &c.LastMessageAt, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan conversation: %w", err)
		}
		conversations = append(conversations, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate conversations: %w", err)
	}
	return conversations, total, nil
}

func (r *ConversationRepository) AddParticipant(ctx context.Context, conversationID, userID string) error {
	query := `
		INSERT INTO conversation_participants (conversation_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (conversation_id, user_id) DO NOTHING
	`
	_, err := r.db.Pool.Exec(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("add participant: %w", err)
	}
	return nil
}

func (r *ConversationRepository) UpdateLastReadAt(ctx context.Context, conversationID, userID string) error {
	query := `
		UPDATE conversation_participants
		SET last_read_at = now()
		WHERE conversation_id = $1 AND user_id = $2
	`
	tag, err := r.db.Pool.Exec(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("update last read at: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("participant not found")
	}
	return nil
}

func (r *ConversationRepository) UpdateLastMessageAt(ctx context.Context, conversationID string) error {
	query := `UPDATE chat_conversations SET last_message_at = now() WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, conversationID)
	if err != nil {
		return fmt.Errorf("update last message at: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found")
	}
	return nil
}
