package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"task-manager/internal/auth"
	"task-manager/internal/utils"
)

func JWTAuthMiddleware(tokenManager *utils.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		tokenStr := parts[1]

		log.Println("tokenStr:", tokenStr)

		claims, err := tokenManager.Parse(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		log.Println(claims)

		identity := auth.NewIdentity(claims.UserID, claims.Role)

		log.Println(identity)

		ctx := c.Request.Context()
		ctx = auth.WithIdentity(ctx, identity)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
