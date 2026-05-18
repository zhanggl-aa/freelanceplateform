package handler

import (
	"net/http"
	"time"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService  *service.ProjectService
	fileService     *service.FileService
	bookmarkService *service.BookmarkService
	reviewService   *service.ReviewService
}

func NewProjectHandler(projectService *service.ProjectService, fileService *service.FileService, bookmarkService *service.BookmarkService, reviewService *service.ReviewService) *ProjectHandler {
	return &ProjectHandler{
		projectService:  projectService,
		fileService:     fileService,
		bookmarkService: bookmarkService,
		reviewService:   reviewService,
	}
}

// RegisterPublicRoutes registers public project routes (no JWT required).
func (h *ProjectHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.GET("", h.Search)
	projects.GET("/:id", h.GetByID)
	projects.GET("/:id/reviews", h.GetReviewsByProject)
}

func (h *ProjectHandler) GetReviewsByProject(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		BadRequest(c, "project id is required")
		return
	}

	reviews, err := h.reviewService.GetByProject(c.Request.Context(), projectID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, reviews)
}

// RegisterRoutes registers protected project routes (JWT required).
func (h *ProjectHandler) RegisterRoutes(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.POST("", middleware.RequireUserType("client", "both"), h.Create)
	projects.PUT("/:id", h.Update)
	projects.DELETE("/:id", h.Delete)
	projects.POST("/:id/publish", h.Publish)
	projects.POST("/:id/close", h.Close)
	projects.POST("/:id/attachments", h.UploadAttachment)
	projects.DELETE("/:id/attachments/:attachmentId", h.DeleteAttachment)
	projects.POST("/:id/bookmark", h.Bookmark)
	projects.DELETE("/:id/bookmark", h.RemoveBookmark)
	projects.GET("/my/posted", h.ListMyPosted)
	projects.GET("/my/bidding", h.ListMyBidding)
	projects.GET("/my/working", h.ListMyWorking)
	projects.GET("/my/completed", h.ListMyCompleted)
}

func (h *ProjectHandler) Search(c *gin.Context) {
	var query struct {
		CategoryID *string  `form:"category_id"`
		Status     *string  `form:"status"`
		Keyword    *string  `form:"keyword"`
		TechStack  []string `form:"tech_stack"`
		model.PaginationQuery
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	keyword := ""
	if query.Keyword != nil {
		keyword = *query.Keyword
	}
	categoryID := ""
	if query.CategoryID != nil {
		categoryID = *query.CategoryID
	}
	status := ""
	if query.Status != nil {
		status = *query.Status
	}

	projects, total, err := h.projectService.Search(
		c.Request.Context(),
		keyword,
		categoryID,
		query.TechStack,
		status,
		query.Page,
		query.PageSize,
	)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, projects, query.Page, query.PageSize, total)
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	project, err := h.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if project == nil {
		NotFound(c, "project not found")
		return
	}

	// Increment view count asynchronously
	go h.projectService.IncrementView(c.Request.Context(), id)

	Success(c, project)
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var body struct {
		CategoryID  string   `json:"category_id" binding:"required"`
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description" binding:"required"`
		BudgetMin   *float64 `json:"budget_min"`
		BudgetMax   *float64 `json:"budget_max"`
		BudgetType  string   `json:"budget_type"`
		Deadline    *string  `json:"deadline"`
		TechStack   []string `json:"tech_stack"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	project := &model.Project{
		ClientID:   userID,
		CategoryID: body.CategoryID,
		Title:      body.Title,
		Description: body.Description,
		BudgetMin:  body.BudgetMin,
		BudgetMax:  body.BudgetMax,
		BudgetType: body.BudgetType,
		TechStack:  body.TechStack,
	}

	if body.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *body.Deadline)
		if err != nil {
			BadRequest(c, "invalid deadline format, expected RFC3339")
			return
		}
		project.Deadline = &t
	}

	created, err := h.projectService.Create(c.Request.Context(), project)
	if err != nil {
		Error(c, http.StatusBadRequest, 40050, err.Error())
		return
	}

	Created(c, created)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	var body struct {
		CategoryID  *string  `json:"category_id"`
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		BudgetMin   *float64 `json:"budget_min"`
		BudgetMax   *float64 `json:"budget_max"`
		BudgetType  *string  `json:"budget_type"`
		Deadline    *string  `json:"deadline"`
		TechStack   []string `json:"tech_stack"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	project := &model.Project{
		ID: id,
	}
	if body.CategoryID != nil {
		project.CategoryID = *body.CategoryID
	}
	if body.Title != nil {
		project.Title = *body.Title
	}
	if body.Description != nil {
		project.Description = *body.Description
	}
	if body.BudgetMin != nil {
		project.BudgetMin = body.BudgetMin
	}
	if body.BudgetMax != nil {
		project.BudgetMax = body.BudgetMax
	}
	if body.BudgetType != nil {
		project.BudgetType = *body.BudgetType
	}
	if body.TechStack != nil {
		project.TechStack = body.TechStack
	}
	if body.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *body.Deadline)
		if err != nil {
			BadRequest(c, "invalid deadline format, expected RFC3339")
			return
		}
		project.Deadline = &t
	}

	updated, err := h.projectService.Update(c.Request.Context(), project)
	if err != nil {
		Error(c, http.StatusBadRequest, 40051, err.Error())
		return
	}

	Success(c, updated)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	if err := h.projectService.Delete(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40052, err.Error())
		return
	}

	Success(c, gin.H{"message": "project deleted successfully"})
}

func (h *ProjectHandler) Publish(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	project, err := h.projectService.Publish(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusBadRequest, 40053, err.Error())
		return
	}

	Success(c, project)
}

func (h *ProjectHandler) Close(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	project, err := h.projectService.Close(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusBadRequest, 40054, err.Error())
		return
	}

	Success(c, project)
}

func (h *ProjectHandler) UploadAttachment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	userID := middleware.GetUserID(c)
	entityType := "project"
	attachment, url, err := h.fileService.Upload(c.Request.Context(), userID, header.Filename, file, header.Header.Get("Content-Type"), header.Size, &entityType, &id)
	if err != nil {
		Error(c, http.StatusInternalServerError, 40055, err.Error())
		return
	}

	Created(c, gin.H{"attachment": attachment, "url": url})
}

func (h *ProjectHandler) DeleteAttachment(c *gin.Context) {
	attachmentID := c.Param("attachmentId")
	if attachmentID == "" {
		BadRequest(c, "attachment id is required")
		return
	}

	if err := h.fileService.Delete(c.Request.Context(), attachmentID); err != nil {
		Error(c, http.StatusBadRequest, 40056, err.Error())
		return
	}

	Success(c, gin.H{"message": "attachment deleted successfully"})
}

func (h *ProjectHandler) Bookmark(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	userID := middleware.GetUserID(c)
	bookmark, err := h.bookmarkService.Create(c.Request.Context(), userID, id)
	if err != nil {
		Error(c, http.StatusBadRequest, 40057, err.Error())
		return
	}

	Created(c, bookmark)
}

func (h *ProjectHandler) RemoveBookmark(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.bookmarkService.Delete(c.Request.Context(), userID, id); err != nil {
		Error(c, http.StatusBadRequest, 40058, err.Error())
		return
	}

	Success(c, gin.H{"message": "bookmark removed successfully"})
}

func (h *ProjectHandler) ListMyPosted(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	projects, total, err := h.projectService.ListByClient(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, projects, query.Page, query.PageSize, total)
}

func (h *ProjectHandler) ListMyBidding(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	projects, total, err := h.projectService.ListByDeveloper(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, projects, query.Page, query.PageSize, total)
}

func (h *ProjectHandler) ListMyWorking(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	projects, total, err := h.projectService.ListByDeveloper(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, projects, query.Page, query.PageSize, total)
}

func (h *ProjectHandler) ListMyCompleted(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	projects, total, err := h.projectService.ListByDeveloper(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, projects, query.Page, query.PageSize, total)
}
