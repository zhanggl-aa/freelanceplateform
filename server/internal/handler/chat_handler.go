package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *service.ChatService
	fileService *service.FileService
}

func NewChatHandler(chatService *service.ChatService, fileService *service.FileService) *ChatHandler {
	return &ChatHandler{chatService: chatService, fileService: fileService}
}

func (h *ChatHandler) RegisterRoutes(rg *gin.RouterGroup) {
	conversations := rg.Group("/conversations")
	conversations.GET("", h.List)
	conversations.GET("/unread-count", h.GetUnreadCount)
	conversations.GET("/:id", h.GetWithMessages)
	conversations.POST("", h.Create)
	conversations.POST("/:id/messages", h.SendMessage)
	conversations.POST("/:id/files", h.SendFile)
	conversations.PUT("/:id/read", h.MarkAsRead)
}

func (h *ChatHandler) List(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	conversations, total, err := h.chatService.ListByUser(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, conversations, query.Page, query.PageSize, total)
}

func (h *ChatHandler) GetWithMessages(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "conversation id is required")
		return
	}

	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	conversation, err := h.chatService.GetByID(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusBadRequest, 40500, err.Error())
		return
	}
	if conversation == nil {
		NotFound(c, "conversation not found")
		return
	}

	messages, total, err := h.chatService.ListMessages(c.Request.Context(), id, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	_ = userID
	Success(c, gin.H{
		"conversation": conversation,
		"messages":     messages,
		"meta": gin.H{
			"page":      query.Page,
			"page_size": query.PageSize,
			"total":     total,
		},
	})
}

func (h *ChatHandler) Create(c *gin.Context) {
	var body struct {
		Type      string  `json:"type" binding:"required,oneof=direct group project"`
		UserID    string  `json:"user_id" binding:"required"`
		ProjectID *string `json:"project_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	var conv *model.ChatConversation
	var err error
	if body.Type == "direct" {
		conv, err = h.chatService.GetOrCreateDirectConversation(c.Request.Context(), userID, body.UserID)
	} else {
		conv, err = h.chatService.CreateConversation(c.Request.Context(), body.Type, body.ProjectID, []string{userID, body.UserID})
	}
	if err != nil {
		Error(c, http.StatusBadRequest, 40501, err.Error())
		return
	}

	Created(c, conv)
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "conversation id is required")
		return
	}

	var body struct {
		Content     string `json:"content" binding:"required"`
		MessageType string `json:"message_type" binding:"required,oneof=text image file system"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	message, err := h.chatService.SendMessage(c.Request.Context(), id, userID, body.Content, body.MessageType, nil, nil, nil)
	if err != nil {
		Error(c, http.StatusBadRequest, 40502, err.Error())
		return
	}

	Created(c, message)
}

func (h *ChatHandler) SendFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "conversation id is required")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		BadRequest(c, "file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		InternalError(c, "failed to open file")
		return
	}
	defer file.Close()

	userID := middleware.GetUserID(c)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, fileURL, err := h.fileService.Upload(c.Request.Context(), userID, fileHeader.Filename, file, contentType, fileHeader.Size, nil, nil)
	if err != nil {
		Error(c, http.StatusInternalServerError, 40503, "failed to upload file")
		return
	}

	fileName := fileHeader.Filename
	fileSize := fileHeader.Size
	messageType := "file"

	message, err := h.chatService.SendMessage(c.Request.Context(), id, userID, "", messageType, &fileURL, &fileName, &fileSize)
	if err != nil {
		Error(c, http.StatusInternalServerError, 40503, err.Error())
		return
	}

	Created(c, message)
}

func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "conversation id is required")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.chatService.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		Error(c, http.StatusBadRequest, 40504, err.Error())
		return
	}

	Success(c, gin.H{"message": "conversation marked as read"})
}

func (h *ChatHandler) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	count, err := h.chatService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{"unread_count": count})
}
