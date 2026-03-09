package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InvalidCredentials(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
}

func NotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}
