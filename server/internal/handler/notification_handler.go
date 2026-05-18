package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	notifications := rg.Group("/notifications")
	notifications.GET("", h.List)
	notifications.GET("/unread-count", h.GetUnreadCount)
	notifications.PUT("/:id/read", h.MarkAsRead)
	notifications.PUT("/read-all", h.MarkAllAsRead)
	notifications.GET("/settings", h.GetSettings)
	notifications.PUT("/settings", h.UpdateSettings)
}

func (h *NotificationHandler) List(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	notifications, total, err := h.notificationService.ListByUser(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, notifications, query.Page, query.PageSize, total)
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{"unread_count": count})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "notification id is required")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.notificationService.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		Error(c, http.StatusBadRequest, 40120, err.Error())
		return
	}

	Success(c, gin.H{"message": "notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		Error(c, http.StatusBadRequest, 40121, err.Error())
		return
	}

	Success(c, gin.H{"message": "all notifications marked as read"})
}

func (h *NotificationHandler) GetSettings(c *gin.Context) {
	userID := middleware.GetUserID(c)
	settings, err := h.notificationService.GetSettings(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, settings)
}

func (h *NotificationHandler) UpdateSettings(c *gin.Context) {
	var settings model.NotificationSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		BadRequest(c, err.Error())
		return
	}

	settings.UserID = middleware.GetUserID(c)
	updated, err := h.notificationService.UpdateSettings(c.Request.Context(), &settings)
	if err != nil {
		Error(c, http.StatusBadRequest, 40122, err.Error())
		return
	}

	Success(c, updated)
}
