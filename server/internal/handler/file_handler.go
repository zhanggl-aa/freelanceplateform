package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) RegisterRoutes(rg *gin.RouterGroup) {
	files := rg.Group("/files")
	files.POST("/upload", h.Upload)
	files.GET("/:id", h.GetByID)
	files.GET("/:id/download", h.Download)
	files.DELETE("/:id", h.Delete)
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	entityType := c.PostForm("entity_type")
	entityID := c.PostForm("entity_id")

	var entityTypePtr *string
	var entityIDPtr *string
	if entityType != "" {
		entityTypePtr = &entityType
	}
	if entityID != "" {
		entityIDPtr = &entityID
	}

	userID := middleware.GetUserID(c)
	contentType := header.Header.Get("Content-Type")
	attachment, url, err := h.fileService.Upload(c.Request.Context(), userID, header.Filename, file, contentType, header.Size, entityTypePtr, entityIDPtr)
	if err != nil {
		Error(c, http.StatusInternalServerError, 40130, err.Error())
		return
	}

	Created(c, gin.H{"attachment": attachment, "url": url})
}

func (h *FileHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "file id is required")
		return
	}

	attachment, err := h.fileService.GetByID(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if attachment == nil {
		NotFound(c, "file not found")
		return
	}

	Success(c, attachment)
}

func (h *FileHandler) Download(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "file id is required")
		return
	}

	url, err := h.fileService.GetURL(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	c.Redirect(http.StatusFound, url)
}

func (h *FileHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "file id is required")
		return
	}

	if err := h.fileService.Delete(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40131, err.Error())
		return
	}

	Success(c, gin.H{"message": "file deleted successfully"})
}
