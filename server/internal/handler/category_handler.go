package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// RegisterPublicRoutes registers public category routes (no JWT required).
func (h *CategoryHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")
	categories.GET("", h.GetTree)
	categories.GET("/:id", h.GetByID)
}

// RegisterRoutes registers protected category routes for admin CRUD.
func (h *CategoryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")
	categories.POST("", h.Create)
	categories.PUT("/:id", h.Update)
	categories.DELETE("/:id", h.Delete)
}

func (h *CategoryHandler) GetTree(c *gin.Context) {
	tree, err := h.categoryService.GetTree(c.Request.Context())
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, tree)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "category id is required")
		return
	}

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if category == nil {
		NotFound(c, "category not found")
		return
	}

	Success(c, category)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var body struct {
		Name        string  `json:"name" binding:"required"`
		Slug        string  `json:"slug" binding:"required"`
		Description *string `json:"description"`
		IconURL     *string `json:"icon_url"`
		ParentID    *string `json:"parent_id"`
		SortOrder   int     `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	category, err := h.categoryService.Create(
		c.Request.Context(),
		body.Name,
		body.Slug,
		body.Description,
		body.IconURL,
		body.ParentID,
		body.SortOrder,
	)
	if err != nil {
		Error(c, http.StatusBadRequest, 40040, err.Error())
		return
	}

	Created(c, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "category id is required")
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	if err := h.categoryService.Update(c.Request.Context(), id, body); err != nil {
		Error(c, http.StatusBadRequest, 40041, err.Error())
		return
	}

	Success(c, gin.H{"message": "category updated successfully"})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "category id is required")
		return
	}

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40042, err.Error())
		return
	}

	Success(c, gin.H{"message": "category deleted successfully"})
}
