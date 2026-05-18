package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type BidHandler struct {
	bidService *service.BidService
}

func NewBidHandler(bidService *service.BidService) *BidHandler {
	return &BidHandler{bidService: bidService}
}

func (h *BidHandler) RegisterRoutes(rg *gin.RouterGroup) {
	bids := rg.Group("/bids")
	bids.POST("/project/:projectId", middleware.RequireUserType("developer", "both"), h.Create)
	bids.GET("/project/:projectId", h.ListByProject)
	bids.GET("/:id", h.GetByID)
	bids.PUT("/:id", h.Update)
	bids.DELETE("/:id", h.Withdraw)
	bids.POST("/:id/accept", middleware.RequireUserType("client", "both"), h.Accept)
	bids.POST("/:id/reject", middleware.RequireUserType("client", "both"), h.Reject)
	bids.POST("/:id/shortlist", middleware.RequireUserType("client", "both"), h.Shortlist)
	bids.POST("/:id/counter-offer", middleware.RequireUserType("client", "both"), h.CounterOffer)

	developers := rg.Group("/developers")
	developers.GET("/me/bids", h.ListMyBids)
}

func (h *BidHandler) Create(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		BadRequest(c, "project id is required")
		return
	}

	var body struct {
		CoverLetter    string  `json:"cover_letter" binding:"required"`
		EstimatedDays  int     `json:"estimated_days" binding:"required,min=1"`
		ProposedBudget float64 `json:"proposed_budget" binding:"required,gt=0"`
		BudgetType     string  `json:"budget_type"`
		MilestonePlan  *string `json:"milestone_plan"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	bid := &model.Bid{
		ProjectID:      projectID,
		DeveloperID:    userID,
		CoverLetter:    body.CoverLetter,
		EstimatedDays:  body.EstimatedDays,
		ProposedBudget: body.ProposedBudget,
		BudgetType:     body.BudgetType,
		MilestonePlan:  body.MilestonePlan,
	}

	created, err := h.bidService.Create(c.Request.Context(), bid)
	if err != nil {
		Error(c, http.StatusBadRequest, 40060, err.Error())
		return
	}

	Created(c, created)
}

func (h *BidHandler) ListByProject(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		BadRequest(c, "project id is required")
		return
	}

	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	bids, total, err := h.bidService.ListByProject(c.Request.Context(), projectID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, bids, query.Page, query.PageSize, total)
}

func (h *BidHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "bid id is required")
		return
	}

	bid, err := h.bidService.GetByID(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if bid == nil {
		NotFound(c, "bid not found")
		return
	}

	Success(c, bid)
}

func (h *BidHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "bid id is required")
		return
	}

	var body struct {
		CoverLetter    *string  `json:"cover_letter"`
		EstimatedDays  *int     `json:"estimated_days"`
		ProposedBudget *float64 `json:"proposed_budget"`
		BudgetType     *string  `json:"budget_type"`
		MilestonePlan  *string  `json:"milestone_plan"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	bid := &model.Bid{
		ID: id,
	}
	if body.CoverLetter != nil {
		bid.CoverLetter = *body.CoverLetter
	}
	if body.EstimatedDays != nil {
		bid.EstimatedDays = *body.EstimatedDays
	}
	if body.ProposedBudget != nil {
		bid.ProposedBudget = *body.ProposedBudget
	}
	if body.BudgetType != nil {
		bid.BudgetType = *body.BudgetType
	}
	if body.MilestonePlan != nil {
		bid.MilestonePlan = body.MilestonePlan
	}

	updated, err := h.bidService.Update(c.Request.Context(), bid)
	if err != nil {
		Error(c, http.StatusBadRequest, 40061, err.Error())
		return
	}

	Success(c, updated)
}

func (h *BidHandler) Withdraw(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "bid id is required")
		return
	}

	if err := h.bidService.Withdraw(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40062, err.Error())
		return
	}

	Success(c, gin.H{"message": "bid withdrawn successfully"})
}

func (h *BidHandler) Accept(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "bid id is required")
		return
	}

	contract, err := h.bidService.Accept(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusBadRequest, 40063, err.Error())
		return
	}

	Success(c, contract)
}

func (h *BidHandler) Reject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "bid id is required")
		return
	}

	var body struct {
		Message *string `json:"message"`
	}
	// Message is optional, so ignore bind errors for empty body
	_ = c.ShouldBindJSON(&body)

	message := ""
	if body.Message != nil {
		message = *body.Message
	}

	if err := h.bidService.Reject(c.Request.Context(), id, message); err != nil {
		Error(c, http.StatusBadRequest, 40064, err.Error())
		return
	}

	Success(c, gin.H{"message": "bid rejected successfully"})
}

func (h *BidHandler) Shortlist(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "bid id is required")
		return
	}

	if err := h.bidService.Shortlist(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40065, err.Error())
		return
	}

	Success(c, gin.H{"message": "bid shortlisted successfully"})
}

func (h *BidHandler) CounterOffer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "bid id is required")
		return
	}

	var body struct {
		ProposedBudget float64 `json:"proposed_budget" binding:"required,gt=0"`
		EstimatedDays  int     `json:"estimated_days" binding:"required,min=1"`
		Message        *string `json:"message"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	message := ""
	if body.Message != nil {
		message = *body.Message
	}

	updated, err := h.bidService.CounterOffer(c.Request.Context(), id, body.ProposedBudget, body.EstimatedDays, message)
	if err != nil {
		Error(c, http.StatusBadRequest, 40066, err.Error())
		return
	}

	Success(c, updated)
}

func (h *BidHandler) ListMyBids(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	bids, total, err := h.bidService.ListByDeveloper(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, bids, query.Page, query.PageSize, total)
}
