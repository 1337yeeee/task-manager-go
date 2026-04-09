package middleware

import (
	"net/http"
	"strings"
	"task-manager/internal/myerrors"

	"github.com/gin-gonic/gin"

	"task-manager/internal/auth"
	"task-manager/internal/utils"
)

func JWTAccessMiddleware(tokenManager utils.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := extractAccessToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := tokenManager.Parse(tokenStr)
		if err != nil || claims.Type != utils.TokenTypeAccess {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid access token",
			})
			return
		}

		identity := auth.NewIdentity(claims.UserID, claims.Role)

		ctx := auth.WithIdentity(c.Request.Context(), identity)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func extractAccessToken(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", myerrors.MissingAuthorizationHeader()
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", myerrors.InvalidAuthorizationHeader()
	}

	return parts[1], nil
}
