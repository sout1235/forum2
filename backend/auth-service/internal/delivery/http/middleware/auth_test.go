package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sout1235/forum2/backend/auth-service/internal/entity"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewAuthUseCase(mocks.NewMockUserRepository(ctrl), "test-key")
	mw := AuthMiddleware(mockUseCase, "internal-token")

	r := gin.New()
	r.GET("/protected", mw, func(c *gin.Context) {
		c.String(200, "ok")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewAuthUseCase(mocks.NewMockUserRepository(ctrl), "test-key")
	mw := AuthMiddleware(mockUseCase, "internal-token")

	r := gin.New()
	r.GET("/protected", mw, func(c *gin.Context) {
		c.String(200, "ok")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestAuthMiddleware_InternalToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewAuthUseCase(mocks.NewMockUserRepository(ctrl), "test-key")
	mw := AuthMiddleware(mockUseCase, "internal-token")

	r := gin.New()
	r.GET("/protected", mw, func(c *gin.Context) {
		c.String(200, "ok")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer internal-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockUseCase := usecase.NewAuthUseCase(mockRepo, "test-key")
	mw := AuthMiddleware(mockUseCase, "internal-token")

	// Создаем тестового пользователя
	testUser := &entity.User{
		ID:       "user-id",
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Генерируем валидный токен
	tokenPair, err := mockUseCase.GenerateTokenPair(testUser)
	assert.NoError(t, err)

	// Мокаем repo.GetByID для usecase.GetUserByID
	mockRepo.EXPECT().GetByID(gomock.Any(), "user-id").Return(testUser, nil)

	r := gin.New()
	r.GET("/protected", mw, func(c *gin.Context) {
		user, exists := c.Get("user")
		if exists {
			c.JSON(200, user)
		} else {
			c.String(401, "no user")
		}
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
}
