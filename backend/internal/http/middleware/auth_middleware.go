package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/auth"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
)

const (
	UserIDContextKey         = "authenticated_user_id"
	OrganizationIDContextKey = "authenticated_organization_id"
	RoleContextKey           = "authenticated_role"
)

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(auth.AccessTokenCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			return
		}

		claims, err := auth.ParseAccessToken(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired session",
			})
			return
		}

		c.Set(UserIDContextKey, claims.Subject)
		c.Set(OrganizationIDContextKey, claims.OrganizationID)
		c.Set(RoleContextKey, claims.Role)

		c.Next()
	}
}

func OrganizationID(c *gin.Context) (string, bool) {
	value, exists := c.Get(OrganizationIDContextKey)
	if !exists {
		return "", false
	}

	organizationID, ok := value.(string)
	if !ok || organizationID == "" {
		return "", false
	}

	return organizationID, true
}

func RequireRole(allowedRoles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(RoleContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session context",
			})
			return
		}

		role, ok := value.(string)
		if !ok || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session context",
			})
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == string(allowedRole) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
	}
}
