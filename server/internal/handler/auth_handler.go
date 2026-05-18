package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes registers public auth routes (no JWT required).
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
	auth.POST("/forgot-password", h.ForgotPassword)
	auth.POST("/reset-password", h.ResetPassword)
}

// RegisterRoutesProtected registers auth routes that require JWT.
func (h *AuthHandler) RegisterRoutesProtected(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/logout", h.Logout)
	auth.PUT("/change-password", h.ChangePassword)
	auth.POST("/verify-email", h.VerifyEmail)
	auth.POST("/verify-phone", h.VerifyPhone)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	user, accessToken, refreshToken, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	Created(c, gin.H{
		"user":         user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		Unauthorized(c, err.Error())
		return
	}

	Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		Unauthorized(c, err.Error())
		return
	}

	Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	if err := h.authService.Logout(c.Request.Context(), body.RefreshToken); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var body struct {
		Account string `json:"account" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	err := h.authService.ForgotPassword(c.Request.Context(), body.Account)
	if err != nil {
		Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}

	// In MVP, retrieve the verification code via GetResetCode for testing.
	code := h.authService.GetResetCode(body.Account)
	Success(c, gin.H{"code": code, "message": "check your email/phone"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var body struct {
		Account     string `json:"account" binding:"required"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), body.Account, body.Code, body.NewPassword); err != nil {
		Error(c, http.StatusBadRequest, 40003, err.Error())
		return
	}

	Success(c, gin.H{"message": "password reset successfully"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.authService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		Error(c, http.StatusBadRequest, 40004, err.Error())
		return
	}

	Success(c, gin.H{"message": "password changed successfully"})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.authService.VerifyEmail(c.Request.Context(), userID); err != nil {
		Error(c, http.StatusBadRequest, 40005, err.Error())
		return
	}

	Success(c, gin.H{"message": "email verified successfully"})
}

func (h *AuthHandler) VerifyPhone(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.authService.VerifyPhone(c.Request.Context(), userID); err != nil {
		Error(c, http.StatusBadRequest, 40006, err.Error())
		return
	}

	Success(c, gin.H{"message": "phone verified successfully"})
}
