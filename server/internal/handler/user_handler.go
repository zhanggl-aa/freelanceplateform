package handler

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService  *service.UserService
	reviewService *service.ReviewService
}

func NewUserHandler(userService *service.UserService, reviewService *service.ReviewService) *UserHandler {
	return &UserHandler{userService: userService, reviewService: reviewService}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("/me", h.GetCurrentUser)
	users.PUT("/me", h.UpdateProfile)
	users.DELETE("/me", h.DeleteAccount)
	users.GET("/:id", h.GetUserByID)
	users.GET("/:id/reviews", h.GetReviewsByUser)
}

func (h *UserHandler) GetReviewsByUser(c *gin.Context) {
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

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, devProfile, clientProfile, err := h.userService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if user == nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, gin.H{
		"user":           user,
		"dev_profile":    devProfile,
		"client_profile": clientProfile,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var body struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, body.Nickname, body.AvatarURL)
	if err != nil {
		Error(c, http.StatusBadRequest, 40010, err.Error())
		return
	}

	Success(c, user)
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.userService.DeleteAccount(c.Request.Context(), userID); err != nil {
		Error(c, http.StatusBadRequest, 40011, err.Error())
		return
	}

	Success(c, gin.H{"message": "account deleted successfully"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "user id is required")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if user == nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, user)
}
