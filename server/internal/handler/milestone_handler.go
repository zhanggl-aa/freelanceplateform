package handler

import (
	"net/http"
	"time"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type MilestoneHandler struct {
	milestoneService *service.MilestoneService
}

func NewMilestoneHandler(milestoneService *service.MilestoneService) *MilestoneHandler {
	return &MilestoneHandler{milestoneService: milestoneService}
}

func (h *MilestoneHandler) RegisterRoutes(rg *gin.RouterGroup) {
	milestones := rg.Group("/milestones")
	milestones.POST("/project/:projectId", h.Create)
	milestones.GET("/project/:projectId", h.ListByProject)
	milestones.PUT("/:id", h.Update)
	milestones.DELETE("/:id", h.Delete)
	milestones.POST("/:id/submit", h.Submit)
	milestones.POST("/:id/approve", h.Approve)
	milestones.POST("/:id/reject", h.Reject)
	milestones.POST("/:id/dispute", h.Dispute)
}

func (h *MilestoneHandler) Create(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		BadRequest(c, "project id is required")
		return
	}

	var body struct {
		Title       string   `json:"title" binding:"required"`
		Description *string  `json:"description"`
		Amount      float64  `json:"amount" binding:"required,gt=0"`
		Deadline    *string  `json:"deadline"`
		SortOrder   int      `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	milestone := &model.ProjectMilestone{
		ProjectID:   projectID,
		Title:       body.Title,
		Description: body.Description,
		Amount:      body.Amount,
		SortOrder:   body.SortOrder,
	}
	if body.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *body.Deadline)
		if err != nil {
			BadRequest(c, "invalid deadline format, expected RFC3339")
			return
		}
		milestone.Deadline = &t
	}

	created, err := h.milestoneService.Create(c.Request.Context(), milestone)
	if err != nil {
		Error(c, http.StatusBadRequest, 40080, err.Error())
		return
	}

	Created(c, created)
}

func (h *MilestoneHandler) ListByProject(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		BadRequest(c, "project id is required")
		return
	}

	milestones, err := h.milestoneService.ListByProject(c.Request.Context(), projectID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, milestones)
}

func (h *MilestoneHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "milestone id is required")
		return
	}

	var body struct {
		Title       string   `json:"title"`
		Description *string  `json:"description"`
		Amount      float64  `json:"amount"`
		Deadline    *string  `json:"deadline"`
		SortOrder   int      `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	milestone := &model.ProjectMilestone{
		ID:          id,
		Title:       body.Title,
		Description: body.Description,
		Amount:      body.Amount,
		SortOrder:   body.SortOrder,
	}
	if body.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *body.Deadline)
		if err != nil {
			BadRequest(c, "invalid deadline format, expected RFC3339")
			return
		}
		milestone.Deadline = &t
	}

	updated, err := h.milestoneService.Update(c.Request.Context(), milestone)
	if err != nil {
		Error(c, http.StatusBadRequest, 40081, err.Error())
		return
	}

	Success(c, updated)
}

func (h *MilestoneHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "milestone id is required")
		return
	}

	if err := h.milestoneService.Delete(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40082, err.Error())
		return
	}

	Success(c, gin.H{"message": "milestone deleted successfully"})
}

func (h *MilestoneHandler) Submit(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "milestone id is required")
		return
	}

	var body struct {
		DeliverableURLs []string `json:"deliverable_urls"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	milestone, err := h.milestoneService.Submit(c.Request.Context(), id, body.DeliverableURLs)
	if err != nil {
		Error(c, http.StatusBadRequest, 40083, err.Error())
		return
	}

	Success(c, milestone)
}

func (h *MilestoneHandler) Approve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "milestone id is required")
		return
	}

	var body struct {
		Feedback *string `json:"feedback"`
	}
	_ = c.ShouldBindJSON(&body)

	feedback := ""
	if body.Feedback != nil {
		feedback = *body.Feedback
	}

	milestone, err := h.milestoneService.Approve(c.Request.Context(), id, feedback)
	if err != nil {
		Error(c, http.StatusBadRequest, 40084, err.Error())
		return
	}

	Success(c, milestone)
}

func (h *MilestoneHandler) Reject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "milestone id is required")
		return
	}

	var body struct {
		Feedback string `json:"feedback" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	milestone, err := h.milestoneService.Reject(c.Request.Context(), id, body.Feedback)
	if err != nil {
		Error(c, http.StatusBadRequest, 40085, err.Error())
		return
	}

	Success(c, milestone)
}

func (h *MilestoneHandler) Dispute(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "milestone id is required")
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	milestone, err := h.milestoneService.Dispute(c.Request.Context(), id, body.Reason)
	if err != nil {
		Error(c, http.StatusBadRequest, 40086, err.Error())
		return
	}

	Success(c, milestone)
}
