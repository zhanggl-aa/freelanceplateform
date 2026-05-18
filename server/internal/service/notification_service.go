package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/repository"
)

type NotificationService struct {
	notificationRepo        *repository.NotificationRepository
	notificationSettingsRepo *repository.NotificationSettingsRepository
}

func NewNotificationService(notificationRepo *repository.NotificationRepository, notificationSettingsRepo *repository.NotificationSettingsRepository) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo, notificationSettingsRepo: notificationSettingsRepo}
}

// Create sends a notification to a user.
func (s *NotificationService) Create(ctx context.Context, notification *model.Notification) (*model.Notification, error) {
	if notification.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if notification.Type == "" {
		return nil, errors.New("type is required")
	}
	if notification.Title == "" {
		return nil, errors.New("title is required")
	}

	validTypes := map[string]bool{
		"bid": true, "message": true, "payment": true, "project": true,
		"milestone": true, "contract": true, "dispute": true, "system": true, "review": true,
	}
	if !validTypes[notification.Type] {
		return nil, errors.New("invalid notification type")
	}

	// Check user notification settings before sending
	settings, err := s.notificationSettingsRepo.GetOrCreateByUserID(ctx, notification.UserID)
	if err != nil {
		return nil, fmt.Errorf("get notification settings: %w", err)
	}
	if settings != nil {
		// Respect user preferences for specific notification types
		switch notification.Type {
		case "bid":
			if !settings.BidNotifications {
				return nil, nil // silently skip
			}
		case "message":
			if !settings.MessageNotifications {
				return nil, nil
			}
		case "payment":
			if !settings.PaymentNotifications {
				return nil, nil
			}
		case "project", "milestone", "contract":
			if !settings.ProjectNotifications {
				return nil, nil
			}
		}

		// Check delivery channel preferences
		if !settings.InAppEnabled {
			return nil, nil // In-app disabled, skip creation
		}
	}

	notification.IsRead = false

	content := ""
	if notification.Content != nil {
		content = *notification.Content
	}

	created, err := s.notificationRepo.Create(ctx, notification.UserID, notification.Type, notification.Title, content, notification.Data)
	if err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}

	return created, nil
}

// ListByUser returns paginated notifications for a user, newest first.
func (s *NotificationService) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.Notification, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user_id is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	notifications, total, err := s.notificationRepo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}

	return notifications, total, nil
}

// GetUnreadCount returns the number of unread notifications for a user.
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, errors.New("user_id is required")
	}

	count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("get unread count: %w", err)
	}

	return count, nil
}

// MarkAsRead marks a single notification as read.
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	if notificationID == "" {
		return errors.New("notification_id is required")
	}
	if userID == "" {
		return errors.New("user_id is required")
	}

	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("get notification: %w", err)
	}
	if notification == nil {
		return errors.New("notification not found")
	}
	if notification.UserID != userID {
		return errors.New("notification does not belong to this user")
	}
	if notification.IsRead {
		return nil // Already read, no-op
	}

	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}

// MarkAllAsRead marks all notifications for a user as read.
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

// GetSettings returns the notification settings for a user.
func (s *NotificationService) GetSettings(ctx context.Context, userID string) (*model.NotificationSettings, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	settings, err := s.notificationSettingsRepo.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get notification settings: %w", err)
	}
	if settings == nil {
		// Return default settings
		settings = &model.NotificationSettings{
			UserID:               userID,
			EmailEnabled:         true,
			SMSEnabled:           false,
			PushEnabled:          true,
			InAppEnabled:         true,
			BidNotifications:     true,
			MessageNotifications: true,
			PaymentNotifications: true,
			ProjectNotifications: true,
		}
	}

	return settings, nil
}

// UpdateSettings updates the notification preferences for a user.
func (s *NotificationService) UpdateSettings(ctx context.Context, settings *model.NotificationSettings) (*model.NotificationSettings, error) {
	if settings.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	updated, err := s.notificationSettingsRepo.GetOrCreateByUserID(ctx, settings.UserID)
	if err != nil {
		return nil, fmt.Errorf("get notification settings: %w", err)
	}

	fields := map[string]interface{}{
		"email_enabled":          settings.EmailEnabled,
		"sms_enabled":            settings.SMSEnabled,
		"push_enabled":           settings.PushEnabled,
		"in_app_enabled":         settings.InAppEnabled,
		"bid_notifications":      settings.BidNotifications,
		"message_notifications":  settings.MessageNotifications,
		"payment_notifications":  settings.PaymentNotifications,
		"project_notifications":  settings.ProjectNotifications,
	}
	err = s.notificationSettingsRepo.Update(ctx, settings.UserID, fields)
	if err != nil {
		return nil, fmt.Errorf("update notification settings: %w", err)
	}

	updated, err = s.notificationSettingsRepo.GetByUserID(ctx, settings.UserID)
	if err != nil {
		return nil, fmt.Errorf("update notification settings: %w", err)
	}

	return updated, nil
}
