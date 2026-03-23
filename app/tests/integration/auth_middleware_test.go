package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/myerrors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task-manager/internal/auth"
	"task-manager/internal/middleware"
)

// Тест 1: Успешная аутентификация с валидным токеном
func TestJWTAuthMiddleware_Success(t *testing.T) {
	r, tokenManager, identity := setupAuthTest()

	// Создаем тестовый токен
	token, err := tokenManager.GenerateAccessToken(identity.UserID, identity.Role)
	require.NoError(t, err)

	// Защищенный эндпоинт
	r.GET("/protected", middleware.JWTAccessMiddleware(tokenManager), func(c *gin.Context) {
		ctxIdentity := auth.FromContext(c.Request.Context())
		assert.NotNil(t, ctxIdentity)
		assert.Equal(t, identity.UserID, ctxIdentity.UserID)
		assert.Equal(t, identity.Role, ctxIdentity.Role)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["message"])
}

// Тест 2: Отсутствует заголовок Authorization
func TestJWTAuthMiddleware_MissingHeader(t *testing.T) {
	r, tokenManager, _ := setupAuthTest()

	r.GET("/protected", middleware.JWTAccessMiddleware(tokenManager), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	// Не добавляем заголовок

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], myerrors.MissingAuthorizationHeader().Error())
}

// Тест 3: Неправильный формат заголовка
func TestJWTAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	r, tokenManager, _ := setupAuthTest()

	r.GET("/protected", middleware.JWTAccessMiddleware(tokenManager), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], myerrors.InvalidAuthorizationHeader().Error())
}

// Тест 4: Невалидный токен
func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	r, tokenManager, _ := setupAuthTest()

	r.GET("/protected", middleware.JWTAccessMiddleware(tokenManager), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid access token")
}

// Тест 5: Успешная авторизация с правильной ролью
func TestRequireRole_Success(t *testing.T) {
	r, tokenManager, identity := setupAuthTest()
	identity.Role = auth.UserRoleAdmin

	token, err := tokenManager.GenerateAccessToken(identity.UserID, identity.Role)
	require.NoError(t, err)

	// Эндпоинт только для админов
	r.GET("/admin",
		middleware.JWTAccessMiddleware(tokenManager),
		middleware.RequireRole(auth.UserRoleAdmin),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		},
	)

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Тест 6: Неудачная авторизация с НЕ правильной ролью
func TestRequireRole_Fail(t *testing.T) {
	r, tokenManager, identity := setupAuthTest()
	identity.Role = auth.UserRoleViewer

	token, err := tokenManager.GenerateAccessToken(identity.UserID, identity.Role)
	require.NoError(t, err)

	// Эндпоинт только для админов
	r.GET("/admin",
		middleware.JWTAccessMiddleware(tokenManager),
		middleware.RequireRole(auth.UserRoleAdmin),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		},
	)

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Тест 7: Тестирование RequireRolesModerators
func TestRequireRolesModerators(t *testing.T) {
	testCases := []struct {
		name           string
		userRole       auth.UserRole
		expectedStatus int
	}{
		{
			name:           "Admin имеет доступ",
			userRole:       auth.UserRoleAdmin,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Модератор имеет доступ",
			userRole:       auth.UserRoleEditor, // предполагаем, что такая роль есть
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Пользователь не имеет доступа",
			userRole:       auth.UserRoleViewer,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, tokenManager, identity := setupAuthTest()
			identity.Role = tc.userRole

			token, err := tokenManager.GenerateAccessToken(identity.UserID, identity.Role)
			require.NoError(t, err)

			// Используем исправленную версию middleware
			r.GET("/moderator-area",
				middleware.JWTAccessMiddleware(tokenManager),
				middleware.RequireRolesModerators(), // используем исправленную версию
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "access granted"})
				},
			)

			req, _ := http.NewRequest("GET", "/moderator-area", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
