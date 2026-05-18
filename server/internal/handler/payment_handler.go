package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	walletService  *service.WalletService
}

func NewPaymentHandler(paymentService *service.PaymentService, walletService *service.WalletService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService, walletService: walletService}
}

func (h *PaymentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	payments := rg.Group("/payments")
	payments.POST("/deposit", h.Deposit)
	payments.POST("/release", h.Release)
	payments.POST("/refund", h.Refund)
	payments.GET("/:id", h.GetByID)
	payments.GET("/my", h.ListByUser)

	wallet := rg.Group("/wallet")
	wallet.GET("/balance", h.GetBalance)
	wallet.GET("/transactions", h.ListTransactions)
	wallet.POST("/withdraw", h.Withdraw)
}

func (h *PaymentHandler) Deposit(c *gin.Context) {
	var body struct {
		ContractID    string  `json:"contract_id" binding:"required"`
		MilestoneID   string  `json:"milestone_id"`
		PayeeID       string  `json:"payee_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		PaymentMethod string  `json:"payment_method"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	payment := &model.Payment{
		ContractID:    body.ContractID,
		PayerID:       userID,
		PayeeID:       body.PayeeID,
		Amount:        body.Amount,
		PaymentMethod: body.PaymentMethod,
	}
	if body.MilestoneID != "" {
		payment.MilestoneID = &body.MilestoneID
	}

	created, err := h.paymentService.Deposit(c.Request.Context(), payment)
	if err != nil {
		Error(c, http.StatusBadRequest, 40090, err.Error())
		return
	}

	Created(c, created)
}

func (h *PaymentHandler) Release(c *gin.Context) {
	var body struct {
		PaymentID string `json:"payment_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	payment, err := h.paymentService.Release(c.Request.Context(), body.PaymentID)
	if err != nil {
		Error(c, http.StatusBadRequest, 40091, err.Error())
		return
	}

	Success(c, payment)
}

func (h *PaymentHandler) Refund(c *gin.Context) {
	var body struct {
		PaymentID string `json:"payment_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	payment, err := h.paymentService.Refund(c.Request.Context(), body.PaymentID)
	if err != nil {
		Error(c, http.StatusBadRequest, 40092, err.Error())
		return
	}

	Success(c, payment)
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "payment id is required")
		return
	}

	payment, err := h.paymentService.GetByID(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if payment == nil {
		NotFound(c, "payment not found")
		return
	}

	Success(c, payment)
}

func (h *PaymentHandler) ListByUser(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	role := c.Query("role")
	if role != "payer" && role != "payee" {
		BadRequest(c, "role must be 'payer' or 'payee'")
		return
	}

	userID := middleware.GetUserID(c)
	payments, total, err := h.paymentService.ListByUser(c.Request.Context(), userID, role, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, payments, query.Page, query.PageSize, total)
}

func (h *PaymentHandler) GetBalance(c *gin.Context) {
	userID := middleware.GetUserID(c)
	wallet, err := h.walletService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, wallet)
}

func (h *PaymentHandler) ListTransactions(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	transactions, total, err := h.walletService.ListTransactions(c.Request.Context(), userID, query.Page, query.PageSize)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithMeta(c, transactions, query.Page, query.PageSize, total)
}

func (h *PaymentHandler) Withdraw(c *gin.Context) {
	var body struct {
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	transaction, err := h.walletService.Withdraw(c.Request.Context(), userID, body.Amount, body.Description)
	if err != nil {
		Error(c, http.StatusBadRequest, 40093, err.Error())
		return
	}

	Success(c, transaction)
}
