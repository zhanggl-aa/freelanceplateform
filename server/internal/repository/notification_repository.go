package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type NotificationRepository struct {
	db *DB
}

func NewNotificationRepository(db *DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, userID, ntype, title, content string, data *string) (*model.Notification, error) {
	var n model.Notification
	query := `
		INSERT INTO notifications (user_id, type, title, content, data)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, type, title, content, data, is_read, read_at, created_at
	`
	err := r.db.Pool.QueryRow(ctx, query,
		userID, ntype, title, nilIfEmpty(content), data,
	).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Content, &n.Data,
		&n.IsRead, &n.ReadAt, &n.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}
	return &n, nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id string) (*model.Notification, error) {
	var n model.Notification
	query := `
		SELECT id, user_id, type, title, content, data, is_read, read_at, created_at
		FROM notifications
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Content, &n.Data,
		&n.IsRead, &n.ReadAt, &n.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get notification by id: %w", err)
	}
	return &n, nil
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.Notification, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	if err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count notifications by user: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, user_id, type, title, content, data, is_read, read_at, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications by user: %w", err)
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Content, &n.Data,
			&n.IsRead, &n.ReadAt, &n.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, &n)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate notifications: %w", err)
	}
	return notifications, total, nil
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get unread notification count: %w", err)
	}
	return count, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE notifications SET is_read = true, read_at = now() WHERE id = $1`
	tag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark notification as read: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `UPDATE notifications SET is_read = true, read_at = now() WHERE user_id = $1 AND is_read = false`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("mark all notifications as read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) DeleteOld(ctx context.Context, before time.Time) error {
	query := `DELETE FROM notifications WHERE created_at < $1 AND is_read = true`
	_, err := r.db.Pool.Exec(ctx, query, before)
	if err != nil {
		return fmt.Errorf("delete old notifications: %w", err)
	}
	return nil
}
