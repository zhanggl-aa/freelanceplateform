package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// RegisterPublicRoutes registers public review routes (no JWT required).
func (h *ReviewHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	reviews := rg.Group("/reviews")
	reviews.GET("/project/:projectId", h.GetByProject)
	reviews.GET("/user/:id", h.GetByReviewee)
}

// RegisterRoutes registers protected review routes (JWT required).
func (h *ReviewHandler) RegisterRoutes(rg *gin.RouterGroup) {
	reviews := rg.Group("/reviews")
	reviews.POST("/project/:projectId", h.Create)
	reviews.PUT("/:id", h.Update)
	reviews.DELETE("/:id", h.Delete)
}

func (h *ReviewHandler) Create(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		BadRequest(c, "project id is required")
		return
	}

	var body struct {
		ContractID          string  `json:"contract_id" binding:"required"`
		RevieweeID          string  `json:"reviewee_id" binding:"required"`
		QualityRating       int     `json:"quality_rating" binding:"required,min=1,max=5"`
		CommunicationRating int     `json:"communication_rating" binding:"required,min=1,max=5"`
		TimelinessRating    int     `json:"timeliness_rating" binding:"required,min=1,max=5"`
		Comment             *string `json:"comment"`
		IsPublic            bool    `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	review := &model.Review{
		ProjectID:           projectID,
		ContractID:          body.ContractID,
		ReviewerID:          userID,
		RevieweeID:          body.RevieweeID,
		QualityRating:       body.QualityRating,
		CommunicationRating: body.CommunicationRating,
		TimelinessRating:    body.TimelinessRating,
		Comment:             body.Comment,
		IsPublic:            body.IsPublic,
	}

	created, err := h.reviewService.Create(c.Request.Context(), review)
	if err != nil {
		Error(c, http.StatusBadRequest, 40110, err.Error())
		return
	}

	Created(c, created)
}

func (h *ReviewHandler) GetByProject(c *gin.Context) {
	projectID := c.Param("projectId")
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

func (h *ReviewHandler) GetByReviewee(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "user id is required")
		return
	}

	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	reviews, total, err := h.reviewService.GetByReviewee(c.Request.Context(), id, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, reviews, query.Page, query.PageSize, total)
}

func (h *ReviewHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "review id is required")
		return
	}

	var body struct {
		QualityRating       int     `json:"quality_rating" binding:"min=1,max=5"`
		CommunicationRating int     `json:"communication_rating" binding:"min=1,max=5"`
		TimelinessRating    int     `json:"timeliness_rating" binding:"min=1,max=5"`
		Comment             *string `json:"comment"`
		IsPublic            bool    `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	review := &model.Review{
		ID:                  id,
		QualityRating:       body.QualityRating,
		CommunicationRating: body.CommunicationRating,
		TimelinessRating:    body.TimelinessRating,
		Comment:             body.Comment,
		IsPublic:            body.IsPublic,
	}

	updated, err := h.reviewService.Update(c.Request.Context(), review)
	if err != nil {
		Error(c, http.StatusBadRequest, 40111, err.Error())
		return
	}

	Success(c, updated)
}

func (h *ReviewHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "review id is required")
		return
	}

	if err := h.reviewService.Delete(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40112, err.Error())
		return
	}

	Success(c, gin.H{"message": "review deleted successfully"})
}
