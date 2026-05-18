package middleware

import (
	"net/http"
	"strings"

	"github.com/freelanceplatform/server/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserID   = "user_id"
	ContextUserType = "user_type"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 40101, "message": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 40102, "message": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims := &model.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 40103, "message": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserType, claims.UserType)
		c.Next()
	}
}

func RequireUserType(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get(ContextUserType)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 40104, "message": "unauthorized"})
			c.Abort()
			return
		}

		ut := userType.(string)
		for _, t := range allowedTypes {
			if ut == t || ut == "both" {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"code": 40301, "message": "insufficient permissions"})
		c.Abort()
	}
}

func GetUserID(c *gin.Context) string {
	val, _ := c.Get(ContextUserID)
	return val.(string)
}

func GetUserType(c *gin.Context) string {
	val, _ := c.Get(ContextUserType)
	return val.(string)
}
