package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/jackc/pgx/v5"
)

type NotificationSettingsRepository struct {
	db *DB
}

func NewNotificationSettingsRepository(db *DB) *NotificationSettingsRepository {
	return &NotificationSettingsRepository{db: db}
}

func (r *NotificationSettingsRepository) GetByUserID(ctx context.Context, userID string) (*model.NotificationSettings, error) {
	var ns model.NotificationSettings
	query := `
		SELECT id, user_id, email_enabled, sms_enabled, push_enabled, in_app_enabled,
			bid_notifications, message_notifications, payment_notifications, project_notifications,
			created_at, updated_at
		FROM notification_settings
		WHERE user_id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&ns.ID, &ns.UserID, &ns.EmailEnabled, &ns.SMSEnabled, &ns.PushEnabled, &ns.InAppEnabled,
		&ns.BidNotifications, &ns.MessageNotifications, &ns.PaymentNotifications, &ns.ProjectNotifications,
		&ns.CreatedAt, &ns.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get notification settings by user id: %w", err)
	}
	return &ns, nil
}

func (r *NotificationSettingsRepository) Create(ctx context.Context, userID string) (*model.NotificationSettings, error) {
	var ns model.NotificationSettings
	query := `
		INSERT INTO notification_settings (user_id)
		VALUES ($1)
		RETURNING id, user_id, email_enabled, sms_enabled, push_enabled, in_app_enabled,
			bid_notifications, message_notifications, payment_notifications, project_notifications,
			created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&ns.ID, &ns.UserID, &ns.EmailEnabled, &ns.SMSEnabled, &ns.PushEnabled, &ns.InAppEnabled,
		&ns.BidNotifications, &ns.MessageNotifications, &ns.PaymentNotifications, &ns.ProjectNotifications,
		&ns.CreatedAt, &ns.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create notification settings: %w", err)
	}
	return &ns, nil
}

func (r *NotificationSettingsRepository) GetOrCreateByUserID(ctx context.Context, userID string) (*model.NotificationSettings, error) {
	ns, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if ns != nil {
		return ns, nil
	}
	return r.Create(ctx, userID)
}

func (r *NotificationSettingsRepository) Update(ctx context.Context, userID string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	allowedColumns := map[string]bool{
		"email_enabled":          true,
		"sms_enabled":            true,
		"push_enabled":           true,
		"in_app_enabled":         true,
		"bid_notifications":      true,
		"message_notifications":  true,
		"payment_notifications":  true,
		"project_notifications":  true,
	}

	var setClauses []string
	var args []interface{}
	argIdx := 1

	for col, val := range fields {
		if !allowedColumns[col] {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, val)
		argIdx++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	args = append(args, userID)
	query := fmt.Sprintf(
		"UPDATE notification_settings SET %s WHERE user_id = $%d",
		strings.Join(setClauses, ", "), argIdx,
	)

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update notification settings: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification settings not found")
	}
	return nil
}
