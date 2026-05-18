package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientService *service.ClientService
	fileService   *service.FileService
}

func NewClientHandler(clientService *service.ClientService, fileService *service.FileService) *ClientHandler {
	return &ClientHandler{clientService: clientService, fileService: fileService}
}

func (h *ClientHandler) RegisterRoutes(rg *gin.RouterGroup) {
	clients := rg.Group("/clients")
	clients.POST("/profile", h.CreateProfile)
	clients.GET("/profile", h.GetProfile)
	clients.PUT("/profile", h.UpdateProfile)
}

func (h *ClientHandler) CreateProfile(c *gin.Context) {
	var profile model.ClientProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		BadRequest(c, err.Error())
		return
	}

	profile.UserID = middleware.GetUserID(c)
	created, err := h.clientService.CreateProfile(c.Request.Context(), &profile)
	if err != nil {
		Error(c, http.StatusBadRequest, 40030, err.Error())
		return
	}

	Created(c, created)
}

func (h *ClientHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	profile, err := h.clientService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if profile == nil {
		NotFound(c, "client profile not found")
		return
	}

	Success(c, profile)
}

func (h *ClientHandler) UpdateProfile(c *gin.Context) {
	var profile model.ClientProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	// Service expects a *model.ClientProfile with the ID field set. Look up the existing profile first.
	existing, err := h.clientService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if existing == nil {
		NotFound(c, "client profile not found")
		return
	}

	profile.ID = existing.ID
	profile.UserID = userID
	updated, err := h.clientService.UpdateProfile(c.Request.Context(), &profile)
	if err != nil {
		Error(c, http.StatusBadRequest, 40031, err.Error())
		return
	}

	Success(c, updated)
}
