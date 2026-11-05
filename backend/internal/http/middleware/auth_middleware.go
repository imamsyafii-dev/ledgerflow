package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/auth"
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
