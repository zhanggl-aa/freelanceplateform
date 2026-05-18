package middleware

import (
	"net/http"

	"github.com/freelanceplatform/server/internal/repository"
	"github.com/gin-gonic/gin"
)

func RequireAdmin(adminRepo *repository.AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextUserID)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 40100, "message": "unauthorized"})
			c.Abort()
			return
		}

		if !adminRepo.IsAdmin(c.Request.Context(), userID.(string)) {
			c.JSON(http.StatusForbidden, gin.H{"code": 40300, "message": "admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
