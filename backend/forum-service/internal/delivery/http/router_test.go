package httpDelivery

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sout1235/forum2/backend/forum-service/internal/delivery/http/middleware"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTopicService struct {
	mock.Mock
}

func (m *MockTopicService) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicService) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Topic), args.Error(1)
}

func (m *MockTopicService) GetAllTopics(ctx context.Context) ([]*entity.Topic, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Topic), args.Error(1)
}

func (m *MockTopicService) UpdateTopic(ctx context.Context, topic *entity.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicService) DeleteTopic(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTopicService) IncrementViews(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTopicService) UpdateCommentCount(ctx context.Context, topicID int64) error {
	args := m.Called(ctx, topicID)
	return args.Error(0)
}

type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) GetRecentMessages(ctx context.Context, limit int) ([]*entity.ChatMessage, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*entity.ChatMessage), args.Error(1)
}

func (m *MockChatRepository) SaveMessage(ctx context.Context, message *entity.ChatMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockChatRepository) DeleteExpiredMessages(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type RouterTestUserRepoMock struct {
	mock.Mock
}

func (m *RouterTestUserRepoMock) GetUsernameByID(ctx context.Context, id int64) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *RouterTestUserRepoMock) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем тестовый сервер авторизации
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"user_id":"1","username":"testuser"}`))
	}))
	defer authServer.Close()

	// Создаем моки
	topicService := new(MockTopicService)
	commentUseCase := new(MockCommentUseCase)
	userRepo := new(RouterTestUserRepoMock)
	chatRepo := new(MockChatRepository)

	// Настраиваем моки
	topicService.On("GetAllTopics", mock.Anything).Return([]*entity.Topic{}, nil)
	topicService.On("GetTopicByID", mock.Anything, int64(1)).Return(&entity.Topic{ID: 1}, nil)
	commentUseCase.On("GetCommentsByTopicID", mock.Anything, int64(1)).Return([]*entity.Comment{}, nil)
	commentUseCase.On("GetCommentByID", mock.Anything, int64(1)).Return(&entity.Comment{ID: 1}, nil)

	// Для негативных сценариев (сначала!)
	topicService.On("CreateTopic", mock.Anything, mock.MatchedBy(func(t *entity.Topic) bool { return t.Title == "" || t.Content == "" })).Return(assert.AnError)
	topicService.On("UpdateTopic", mock.Anything, mock.MatchedBy(func(t *entity.Topic) bool { return t.Title == "" || t.Content == "" })).Return(assert.AnError)
	commentUseCase.On("CreateComment", mock.Anything, mock.MatchedBy(func(c *entity.Comment) bool { return c.Content == "" })).Return(assert.AnError)

	// Для успешных сценариев (после негативных)
	topicService.On("CreateTopic", mock.Anything, mock.Anything).Return(nil)
	topicService.On("UpdateTopic", mock.Anything, mock.AnythingOfType("*entity.Topic")).Return(nil)
	topicService.On("DeleteTopic", mock.Anything, int64(1)).Return(nil)
	commentUseCase.On("CreateComment", mock.Anything, mock.Anything).Return(nil)
	commentUseCase.On("DeleteComment", mock.Anything, int64(1)).Return(nil)
	userRepo.On("GetUsernameByID", mock.Anything, int64(1)).Return("testuser", nil)

	// Создаем роутер с тестовым auth-сервисом
	authCfg := &middleware.AuthConfig{AuthServiceURL: authServer.URL}
	router := NewRouter(topicService, commentUseCase, userRepo, chatRepo, "8080", authCfg)

	// Тестируем маршруты для топиков
	t.Run("Topic Routes", func(t *testing.T) {
		tests := []struct {
			name       string
			method     string
			path       string
			wantStatus int
		}{
			{"Create Topic", "POST", "/api/v1/topics", http.StatusUnauthorized},
			{"Get All Topics", "GET", "/api/v1/topics", http.StatusOK},
			{"Get Topic", "GET", "/api/v1/topics/1", http.StatusOK},
			{"Update Topic", "PUT", "/api/v1/topics/1", http.StatusUnauthorized},
			{"Delete Topic", "DELETE", "/api/v1/topics/1", http.StatusUnauthorized},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(tt.method, tt.path, nil)
				router.engine.ServeHTTP(w, req)
				assert.Equal(t, tt.wantStatus, w.Code)
			})
		}
	})

	// Тестируем маршруты для комментариев
	t.Run("Comment Routes", func(t *testing.T) {
		tests := []struct {
			name       string
			method     string
			path       string
			wantStatus int
		}{
			{"Get All Comments", "GET", "/api/v1/topics/1/comments", http.StatusOK},
			{"Get Comment", "GET", "/api/v1/comments/1", http.StatusOK},
			{"Create Comment", "POST", "/api/v1/topics/1/comments", http.StatusUnauthorized},
			{"Delete Comment", "DELETE", "/api/v1/topics/1/comments/1", http.StatusUnauthorized},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(tt.method, tt.path, nil)
				router.engine.ServeHTTP(w, req)
				assert.Equal(t, tt.wantStatus, w.Code)
			})
		}
	})

	// Тестируем маршруты для чата
	t.Run("Chat Routes", func(t *testing.T) {
		tests := []struct {
			name       string
			method     string
			path       string
			wantStatus int
		}{
			{"Get Messages", "GET", "/api/v1/chat/messages", http.StatusUnauthorized},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(tt.method, tt.path, nil)
				router.engine.ServeHTTP(w, req)
				assert.Equal(t, tt.wantStatus, w.Code)
			})
		}
	})

	// Тестируем несуществующие маршруты
	t.Run("Non-existent Routes", func(t *testing.T) {
		tests := []struct {
			name       string
			method     string
			path       string
			wantStatus int
		}{
			{"Non-existent GET", "GET", "/api/v1/non-existent", http.StatusNotFound},
			{"Non-existent POST", "POST", "/api/v1/non-existent", http.StatusNotFound},
			{"Non-existent PUT", "PUT", "/api/v1/non-existent", http.StatusNotFound},
			{"Non-existent DELETE", "DELETE", "/api/v1/non-existent", http.StatusNotFound},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(tt.method, tt.path, nil)
				router.engine.ServeHTTP(w, req)
				assert.Equal(t, tt.wantStatus, w.Code)
			})
		}
	})

	// Добавляем тесты для успешных авторизованных запросов
	t.Run("Authorized Topic Actions", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"title":"Test","content":"Test content","user_id":1}`
		req, _ := http.NewRequest("POST", "/api/v1/topics", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer testtoken")
		req.Header.Set("Content-Type", "application/json")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		w = httptest.NewRecorder()
		body = `{"title":"Updated","content":"Updated content","user_id":1}`
		req, _ = http.NewRequest("PUT", "/api/v1/topics/1", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer testtoken")
		req.Header.Set("Content-Type", "application/json")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/api/v1/topics/1", nil)
		req.Header.Set("Authorization", "Bearer testtoken")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Authorized Comment Actions", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"content":"Test comment","user_id":1}`
		req, _ := http.NewRequest("POST", "/api/v1/topics/1/comments", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer testtoken")
		req.Header.Set("Content-Type", "application/json")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/api/v1/topics/1/comments/1", nil)
		req.Header.Set("Authorization", "Bearer testtoken")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	// Негативные сценарии (400, 500)
	t.Run("Bad Request Topic - invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"title":123}` // title должен быть строкой
		req, _ := http.NewRequest("POST", "/api/v1/topics", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer testtoken")
		req.Header.Set("Content-Type", "application/json")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Bad Request Topic - empty fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"title":"","content":"","user_id":1}`
		req, _ := http.NewRequest("POST", "/api/v1/topics", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer testtoken")
		req.Header.Set("Content-Type", "application/json")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Bad Request Comment - invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"content":123}` // content должен быть строкой
		req, _ := http.NewRequest("POST", "/api/v1/topics/1/comments", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer testtoken")
		req.Header.Set("Content-Type", "application/json")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Bad Request Comment - empty fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"content":"","user_id":1}`
		req, _ := http.NewRequest("POST", "/api/v1/topics/1/comments", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer testtoken")
		req.Header.Set("Content-Type", "application/json")
		router.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
