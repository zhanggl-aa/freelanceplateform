package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/dashboard", h.Dashboard)

	users := rg.Group("/users")
	users.GET("", h.ListUsers)
	users.PUT("/:id/status", h.UpdateUserStatus)

	projects := rg.Group("/projects")
	projects.GET("", h.ListProjects)
	projects.PUT("/:id", h.ModerateProject)

	disputes := rg.Group("/disputes")
	disputes.GET("", h.ListDisputes)
	disputes.PUT("/:id", h.ResolveDispute)

	payments := rg.Group("/payments")
	payments.GET("", h.ListPayments)

	finance := rg.Group("/finance")
	finance.GET("/summary", h.FinancialSummary)
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	userID := middleware.GetUserID(c)
	dashboard, err := h.adminService.Dashboard(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, dashboard)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var query struct {
		Search  *string `form:"search"`
		Status  *string `form:"status"`
		model.PaginationQuery
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	users, total, err := h.adminService.ListUsers(
		c.Request.Context(),
		query.Search,
		query.Status,
		query.Page,
		query.PageSize,
	)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, users, query.Page, query.PageSize, total)
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "user id is required")
		return
	}

	var body struct {
		Status string `json:"status" binding:"required,oneof=active suspended deleted"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	adminID := middleware.GetUserID(c)
	if err := h.adminService.UpdateUserStatus(c.Request.Context(), adminID, id, body.Status); err != nil {
		Error(c, http.StatusBadRequest, 40140, err.Error())
		return
	}

	Success(c, gin.H{"message": "user status updated successfully"})
}

func (h *AdminHandler) ListProjects(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	projects, total, err := h.adminService.ListProjects(c.Request.Context(), query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, projects, query.Page, query.PageSize, total)
}

func (h *AdminHandler) ModerateProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "project id is required")
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	adminID := middleware.GetUserID(c)
	if err := h.adminService.ModerateProject(c.Request.Context(), adminID, id, body); err != nil {
		Error(c, http.StatusBadRequest, 40141, err.Error())
		return
	}

	Success(c, gin.H{"message": "project moderated successfully"})
}

func (h *AdminHandler) ListDisputes(c *gin.Context) {
	var query struct {
		Status *string `form:"status"`
		model.PaginationQuery
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	statusStr := ""
	if query.Status != nil {
		statusStr = *query.Status
	}

	disputes, total, err := h.adminService.ListDisputes(c.Request.Context(), query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	_ = statusStr

	SuccessWithMeta(c, disputes, query.Page, query.PageSize, total)
}

func (h *AdminHandler) ResolveDispute(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "dispute id is required")
		return
	}

	var body struct {
		Resolution     string `json:"resolution" binding:"required"`
		ResolutionType string `json:"resolution_type" binding:"required,oneof=favor_client favor_developer split dismissed"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	adminID := middleware.GetUserID(c)
	if err := h.adminService.ResolveDispute(c.Request.Context(), adminID, id, body.Resolution, body.ResolutionType); err != nil {
		Error(c, http.StatusBadRequest, 40142, err.Error())
		return
	}

	Success(c, gin.H{"message": "dispute resolved successfully"})
}

func (h *AdminHandler) ListPayments(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	payments, total, err := h.adminService.ListPayments(c.Request.Context(), query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, payments, query.Page, query.PageSize, total)
}

func (h *AdminHandler) FinancialSummary(c *gin.Context) {
	summary, err := h.adminService.FinancialSummary(c.Request.Context())
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, summary)
}
