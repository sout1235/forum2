package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sout1235/forum2/backend/auth-service/internal/entity"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUseCase := usecase.NewAuthUseCase(mockRepo, "test-secret")
	handler := NewAuthHandler(authUseCase, "test-secret")

	tests := []struct {
		name           string
		payload        map[string]interface{}
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful registration",
			payload: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					ExistsByUsername(gomock.Any(), "testuser").
					Return(false, nil)
				mockRepo.EXPECT().
					ExistsByEmail(gomock.Any(), "test@example.com").
					Return(false, nil)
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"access_token":  gomock.Any(),
				"refresh_token": gomock.Any(),
				"user": map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"role":     "user",
				},
			},
		},
		{
			name: "duplicate username",
			payload: map[string]interface{}{
				"username": "existinguser",
				"email":    "test@example.com",
				"password": "password123",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					ExistsByUsername(gomock.Any(), "existinguser").
					Return(true, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"error": "username already exists",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			router := gin.New()
			router.POST("/register", handler.Register)

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedStatus == http.StatusCreated {
				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
				assert.Contains(t, response, "user")
			} else {
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUseCase := usecase.NewAuthUseCase(mockRepo, "test-secret")
	handler := NewAuthHandler(authUseCase, "test-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entity.User{
		ID:           "1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name           string
		payload        map[string]interface{}
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful login",
			payload: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByUsername(gomock.Any(), "testuser").
					Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"access_token":  gomock.Any(),
				"refresh_token": gomock.Any(),
				"user": map[string]interface{}{
					"id":       "1",
					"username": "testuser",
					"email":    "test@example.com",
					"role":     "user",
				},
			},
		},
		{
			name: "user not found",
			payload: map[string]interface{}{
				"username": "nonexistent",
				"password": "password123",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByUsername(gomock.Any(), "nonexistent").
					Return(nil, entity.ErrUserNotFound)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			router := gin.New()
			router.POST("/login", handler.Login)

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
				assert.Contains(t, response, "user")
			} else {
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
		})
	}
}

func TestAuthHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUseCase := usecase.NewAuthUseCase(mockRepo, "test-secret")
	handler := NewAuthHandler(authUseCase, "test-secret")

	testUser := &entity.User{
		ID:        "1",
		Username:  "testuser",
		Email:     "test@example.com",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	router := gin.New()
	router.GET("/profile", func(c *gin.Context) {
		c.Set("user", testUser)
		handler.GetProfile(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response entity.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, response.ID)
	assert.Equal(t, testUser.Username, response.Username)
	assert.Equal(t, testUser.Email, response.Email)
}

func TestAuthHandler_VerifyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUseCase := usecase.NewAuthUseCase(mockRepo, "test-secret")
	handler := NewAuthHandler(authUseCase, "test-secret")

	// Создаем тестового пользователя и генерируем токен
	testUser := &entity.User{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	tokenPair, err := authUseCase.GenerateTokenPair(testUser)
	assert.NoError(t, err)

	// Мокаем GetByID для валидации токена
	mockRepo.EXPECT().
		GetByID(gomock.Any(), "1").
		Return(testUser, nil).
		AnyTimes()

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedValid  bool
	}{
		{
			name:           "valid token",
			token:          tokenPair.AccessToken,
			expectedStatus: http.StatusOK,
			expectedValid:  true,
		},
		{
			name:           "invalid token",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/verify", handler.VerifyToken)

			body, _ := json.Marshal(map[string]string{"token": tt.token})
			req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedValid {
				assert.Equal(t, "1", response["user_id"])
				assert.Equal(t, "testuser", response["username"])
			} else {
				assert.Contains(t, response, "error", "response: %v", response)
			}
		})
	}
}

func TestAuthHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUseCase := usecase.NewAuthUseCase(mockRepo, "test-secret")
	handler := NewAuthHandler(authUseCase, "test-secret")

	testUser := &entity.User{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}

	tests := []struct {
		name           string
		userID         string
		mockBehavior   func()
		expectedStatus int
		expectedUser   *entity.User
	}{
		{
			name:   "successful get user",
			userID: "1",
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "1").
					Return(testUser, nil)
			},
			expectedStatus: http.StatusOK,
			expectedUser:   testUser,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "nonexistent").
					Return(nil, entity.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedUser:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			router := gin.New()
			router.GET("/user/:id", func(c *gin.Context) {
				c.Set("user_id", "1") // Устанавливаем user_id в контексте
				handler.GetUserByID(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/user/"+tt.userID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedUser != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, response["id"])
				assert.Equal(t, tt.expectedUser.Username, response["username"])
			}
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUseCase := usecase.NewAuthUseCase(mockRepo, "test-secret")
	handler := NewAuthHandler(authUseCase, "test-secret")

	// Создаем тестового пользователя и генерируем токены
	testUser := &entity.User{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	tokenPair, err := authUseCase.GenerateTokenPair(testUser)
	assert.NoError(t, err)

	// Мокаем GetByID для валидации токена
	mockRepo.EXPECT().
		GetByID(gomock.Any(), "1").
		Return(testUser, nil).
		AnyTimes()

	tests := []struct {
		name           string
		refreshToken   string
		expectedStatus int
		expectTokens   bool
	}{
		{
			name:           "valid refresh token",
			refreshToken:   tokenPair.RefreshToken,
			expectedStatus: http.StatusOK,
			expectTokens:   true,
		},
		{
			name:           "invalid refresh token",
			refreshToken:   "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectTokens:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/refresh", handler.RefreshToken)

			body, _ := json.Marshal(map[string]string{"refresh_token": tt.refreshToken})
			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectTokens {
				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
			} else {
				assert.Contains(t, response, "error")
			}
		})
	}
}
