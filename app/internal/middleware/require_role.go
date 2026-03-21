package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"task-manager/internal/auth"
)

func RequireRole(role auth.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := auth.FromContext(c.Request.Context())
		if identity == nil || identity.Role != role {
			c.AbortWithStatus(http.StatusForbidden)
		}

		c.Next()
	}
}

func RequireRolesModerators() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := []auth.UserRole{auth.UserRoleAdmin, auth.UserRoleEditor}
		identity := auth.FromContext(c.Request.Context())
		if identity == nil || !slices.Contains(roles, identity.Role) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	}
}
