package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUseCase := usecase.NewAuthUseCase(mockRepo, "test-secret")
	router := NewRouter(authUseCase, "test-secret")
	engine := router.Setup()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "register route",
			method:         http.MethodPost,
			path:           "/api/auth/register",
			expectedStatus: http.StatusBadRequest, // 400 because we're not sending any data
		},
		{
			name:           "login route",
			method:         http.MethodPost,
			path:           "/api/auth/login",
			expectedStatus: http.StatusBadRequest, // 400 because we're not sending any data
		},
		{
			name:           "profile route",
			method:         http.MethodGet,
			path:           "/api/auth/profile",
			expectedStatus: http.StatusUnauthorized, // 401 because we're not sending a token
		},
		{
			name:           "verify token route",
			method:         http.MethodPost,
			path:           "/api/auth/verify",
			expectedStatus: http.StatusBadRequest, // 400 because we're not sending any data
		},
		{
			name:           "refresh token route",
			method:         http.MethodPost,
			path:           "/api/auth/refresh",
			expectedStatus: http.StatusBadRequest, // 400 because we're not sending any data
		},
		{
			name:           "get user by id route",
			method:         http.MethodGet,
			path:           "/api/auth/user/1",
			expectedStatus: http.StatusUnauthorized, // 401 because we're not sending a token
		},
		{
			name:           "non-existent route",
			method:         http.MethodGet,
			path:           "/api/auth/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
