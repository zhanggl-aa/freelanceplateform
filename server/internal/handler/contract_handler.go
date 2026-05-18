package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	contractService *service.ContractService
}

func NewContractHandler(contractService *service.ContractService) *ContractHandler {
	return &ContractHandler{contractService: contractService}
}

func (h *ContractHandler) RegisterRoutes(rg *gin.RouterGroup) {
	contracts := rg.Group("/contracts")
	contracts.GET("/:id", h.GetByID)
	contracts.GET("/my", h.ListByUser)
	contracts.POST("/:id/start", h.Start)
	contracts.POST("/:id/cancel", h.Cancel)
	contracts.POST("/:id/dispute", h.Dispute)
}

func (h *ContractHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "contract id is required")
		return
	}

	contract, err := h.contractService.GetByID(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if contract == nil {
		NotFound(c, "contract not found")
		return
	}

	Success(c, contract)
}

func (h *ContractHandler) ListByUser(c *gin.Context) {
	var query struct {
		Role     string `form:"role" binding:"required,oneof=client developer"`
		model.PaginationQuery
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	contracts, total, err := h.contractService.ListByUser(c.Request.Context(), userID, query.Role, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, contracts, query.Page, query.PageSize, total)
}

func (h *ContractHandler) Start(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "contract id is required")
		return
	}

	contract, err := h.contractService.Start(c.Request.Context(), id)
	if err != nil {
		Error(c, http.StatusBadRequest, 40071, err.Error())
		return
	}

	Success(c, contract)
}

func (h *ContractHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "contract id is required")
		return
	}

	var body struct {
		Reason *string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&body)

	reason := ""
	if body.Reason != nil {
		reason = *body.Reason
	}

	contract, err := h.contractService.Cancel(c.Request.Context(), id, reason)
	if err != nil {
		Error(c, http.StatusBadRequest, 40072, err.Error())
		return
	}

	Success(c, contract)
}

func (h *ContractHandler) Dispute(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "contract id is required")
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	contract, err := h.contractService.Dispute(c.Request.Context(), id, body.Reason)
	if err != nil {
		Error(c, http.StatusBadRequest, 40073, err.Error())
		return
	}

	Success(c, contract)
}
