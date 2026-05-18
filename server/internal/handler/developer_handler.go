package handler

import (
	"net/http"
	"strings"

	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/model"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/gin-gonic/gin"
)

// Suppress unused import warnings
var _ = strings.TrimSpace
var _ = middleware.GetUserID

type DeveloperHandler struct {
	developerService *service.DeveloperService
	userService      *service.UserService
	fileService      *service.FileService
}

func NewDeveloperHandler(developerService *service.DeveloperService, userService *service.UserService, fileService *service.FileService) *DeveloperHandler {
	return &DeveloperHandler{developerService: developerService, userService: userService, fileService: fileService}
}

// RegisterPublicRoutes registers public developer routes (no JWT required).
func (h *DeveloperHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	developers := rg.Group("/developers")
	developers.GET("", h.Search)
	developers.GET("/:id", h.GetByUserID)
}

// RegisterRoutes registers protected developer routes (JWT required).
func (h *DeveloperHandler) RegisterRoutes(rg *gin.RouterGroup) {
	developers := rg.Group("/developers")
	developers.POST("/profile", h.CreateProfile)
	developers.GET("/profile", h.GetProfile)
	developers.PUT("/profile", h.UpdateProfile)
	developers.POST("/skills", h.AddSkill)
	developers.PUT("/skills/:id", h.UpdateSkill)
	developers.DELETE("/skills/:id", h.DeleteSkill)
	developers.POST("/portfolio", h.AddPortfolio)
	developers.PUT("/portfolio/:id", h.UpdatePortfolio)
	developers.DELETE("/portfolio/:id", h.DeletePortfolio)
}

func (h *DeveloperHandler) Search(c *gin.Context) {
	var query struct {
		Keyword      string   `form:"keyword"`
		Skills       string   `form:"skills"`
		Availability string   `form:"availability"`
		MinRate      *float64 `form:"min_rate"`
		MaxRate      *float64 `form:"max_rate"`
		Sort         string   `form:"sort"`
		model.PaginationQuery
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, err.Error())
		return
	}

	var skills []string
	if query.Skills != "" {
		skills = strings.Split(query.Skills, ",")
	}

	profiles, total, err := h.developerService.Search(
		c.Request.Context(),
		query.Keyword,
		skills,
		query.Availability,
		query.MinRate,
		query.MaxRate,
		query.Page,
		query.PageSize,
	)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	// Enrich each profile with user info, skills, and portfolio
	type enrichedProfile struct {
		*model.DeveloperProfile
		Nickname   string               `json:"nickname"`
		AvatarURL  string               `json:"avatar_url"`
		Skills     []*model.DeveloperSkill     `json:"skills"`
		Portfolio  []*model.DeveloperPortfolio `json:"portfolio"`
	}

	items := make([]enrichedProfile, 0, len(profiles))
	for _, p := range profiles {
		ep := enrichedProfile{DeveloperProfile: p}
		user, err := h.userService.GetUserByID(c.Request.Context(), p.UserID)
		if err == nil && user != nil {
			ep.Nickname = user.Nickname
			if user.AvatarURL != nil {
				ep.AvatarURL = *user.AvatarURL
			}
		}
		sk, err := h.developerService.ListSkills(c.Request.Context(), p.ID)
		if err == nil {
			ep.Skills = sk
		}
		pf, err := h.developerService.ListPortfolio(c.Request.Context(), p.ID)
		if err == nil {
			ep.Portfolio = pf
		}
		items = append(items, ep)
	}

	SuccessWithMeta(c, items, query.Page, query.PageSize, total)
}

func (h *DeveloperHandler) GetByUserID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "developer id is required")
		return
	}

	// Try to get by user ID first
	profile, err := h.developerService.GetByUserID(c.Request.Context(), id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	// If not found, try to get by profile ID directly
	if profile == nil {
		profile, err = h.developerService.GetProfile(c.Request.Context(), id)
		if err != nil {
			InternalError(c, err.Error())
			return
		}
	}
	if profile == nil {
		NotFound(c, "developer profile not found")
		return
	}

	type enrichedProfile struct {
		*model.DeveloperProfile
		Nickname  string               `json:"nickname"`
		AvatarURL string               `json:"avatar_url"`
		Skills    []*model.DeveloperSkill     `json:"skills"`
		Portfolio []*model.DeveloperPortfolio `json:"portfolio"`
	}

	ep := enrichedProfile{DeveloperProfile: profile}
	user, err := h.userService.GetUserByID(c.Request.Context(), profile.UserID)
	if err == nil && user != nil {
		ep.Nickname = user.Nickname
		if user.AvatarURL != nil {
			ep.AvatarURL = *user.AvatarURL
		}
	}
	sk, err := h.developerService.ListSkills(c.Request.Context(), profile.ID)
	if err == nil {
		ep.Skills = sk
	}
	pf, err := h.developerService.ListPortfolio(c.Request.Context(), profile.ID)
	if err == nil {
		ep.Portfolio = pf
	}

	Success(c, ep)
}

func (h *DeveloperHandler) CreateProfile(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	profile, err := h.developerService.CreateProfile(c.Request.Context(), userID, body)
	if err != nil {
		Error(c, http.StatusBadRequest, 40020, err.Error())
		return
	}

	Created(c, profile)
}

func (h *DeveloperHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	profile, err := h.developerService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if profile == nil {
		NotFound(c, "developer profile not found")
		return
	}

	Success(c, profile)
}

func (h *DeveloperHandler) UpdateProfile(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	// Service expects the profile ID, not the user ID. Look up the profile first.
	existing, err := h.developerService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if existing == nil {
		NotFound(c, "developer profile not found")
		return
	}

	profile, err := h.developerService.UpdateProfile(c.Request.Context(), existing.ID, body)
	if err != nil {
		Error(c, http.StatusBadRequest, 40021, err.Error())
		return
	}

	Success(c, profile)
}

func (h *DeveloperHandler) AddSkill(c *gin.Context) {
	var body struct {
		SkillName       string `json:"skill_name" binding:"required"`
		Proficiency     string `json:"proficiency"`
		YearsExperience *int   `json:"years_experience"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	skill, err := h.developerService.AddSkill(c.Request.Context(), userID, body.SkillName, body.Proficiency, body.YearsExperience)
	if err != nil {
		Error(c, http.StatusBadRequest, 40022, err.Error())
		return
	}

	Created(c, skill)
}

func (h *DeveloperHandler) UpdateSkill(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "skill id is required")
		return
	}

	var body struct {
		Proficiency     string `json:"proficiency" binding:"required"`
		YearsExperience *int   `json:"years_experience"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	if err := h.developerService.UpdateSkill(c.Request.Context(), id, body.Proficiency, body.YearsExperience); err != nil {
		Error(c, http.StatusBadRequest, 40023, err.Error())
		return
	}

	Success(c, gin.H{"message": "skill updated successfully"})
}

func (h *DeveloperHandler) DeleteSkill(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "skill id is required")
		return
	}

	if err := h.developerService.DeleteSkill(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40024, err.Error())
		return
	}

	Success(c, gin.H{"message": "skill deleted successfully"})
}

func (h *DeveloperHandler) AddPortfolio(c *gin.Context) {
	var body model.DeveloperPortfolio
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	portfolio, err := h.developerService.AddPortfolio(c.Request.Context(), userID, &body)
	if err != nil {
		Error(c, http.StatusBadRequest, 40025, err.Error())
		return
	}

	Created(c, portfolio)
}

func (h *DeveloperHandler) UpdatePortfolio(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "portfolio id is required")
		return
	}

	var body model.DeveloperPortfolio
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}

	portfolio, err := h.developerService.UpdatePortfolio(c.Request.Context(), id, &body)
	if err != nil {
		Error(c, http.StatusBadRequest, 40026, err.Error())
		return
	}

	Success(c, portfolio)
}

func (h *DeveloperHandler) DeletePortfolio(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "portfolio id is required")
		return
	}

	if err := h.developerService.DeletePortfolio(c.Request.Context(), id); err != nil {
		Error(c, http.StatusBadRequest, 40027, err.Error())
		return
	}

	Success(c, gin.H{"message": "portfolio deleted successfully"})
}
