package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-manager/internal/auth"
	"task-manager/internal/utils"
)

const RefreshTokenContextKey = "refreshToken"

func JWTRefreshMiddleware(tokenManager *utils.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := tokenManager.Parse(tokenStr)
		if err != nil || claims.Type != utils.TokenTypeRefresh {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid refresh token",
			})
			return
		}

		identity := auth.NewIdentity(claims.UserID, claims.Role)
		ctx := auth.WithIdentity(c.Request.Context(), identity)
		c.Request = c.Request.WithContext(ctx)

		c.Set(RefreshTokenContextKey, tokenStr)

		c.Next()
	}
}
