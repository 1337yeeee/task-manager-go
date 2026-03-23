package integration

import (
	"github.com/gin-gonic/gin"
	"task-manager/internal/auth"
	"task-manager/internal/utils"
)

func setupAuthTest() (*gin.Engine, *utils.TokenManager, *auth.Identity) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	tokenManager := utils.NewTokenManager("test-secret", utils.DefaultAccessTTL, utils.DefaultRefreshTTL)

	// Создаем тестовую identity
	identity := auth.NewIdentity("test-user-id", auth.UserRoleViewer)

	return r, tokenManager, identity
}
