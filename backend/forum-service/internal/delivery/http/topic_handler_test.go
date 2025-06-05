package httpDelivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
	"github.com/sout1235/forum2/backend/forum-service/internal/service"
	"github.com/sout1235/forum2/backend/forum-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *TopicHandler) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	topicRepo := &MockTopicRepository{}
	commentRepo := &MockCommentRepository{}
	userRepo := &MockUserRepository{}

	topicRepoIface := repository.TopicRepository(topicRepo)
	commentRepoIface := repository.CommentRepository(commentRepo)
	userRepoIface := repository.UserRepository(userRepo)

	uc := usecase.NewTopicUseCase(topicRepoIface, commentRepoIface, userRepoIface)
	var topicService service.TopicService = uc
	h := NewTopicHandler(topicService, userRepoIface)

	return r, h
}

func TestTopicHandler_CreateTopic(t *testing.T) {
	r, h := setupTestRouter()
	r.POST("/topics", h.CreateTopic)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "success",
			body: map[string]interface{}{
				"title":   "Test Topic",
				"content": "Test Content",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing title",
			body: map[string]interface{}{
				"content": "Test Content",
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "missing content",
			body: map[string]interface{}{
				"title": "Test Topic",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/topics", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestTopicHandler_GetTopic(t *testing.T) {
	r, h := setupTestRouter()
	r.GET("/topics/:id", h.GetTopic)

	tests := []struct {
		name       string
		topicID    string
		wantStatus int
	}{
		{
			name:       "success",
			topicID:    "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			topicID:    "999",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/topics/"+tt.topicID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestTopicHandler_GetAllTopics(t *testing.T) {
	r, h := setupTestRouter()
	r.GET("/topics", h.GetAllTopics)

	req, _ := http.NewRequest(http.MethodGet, "/topics", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
}

func TestTopicHandler_UpdateTopic(t *testing.T) {
	r, h := setupTestRouter()
	r.PUT("/topics/:id", h.UpdateTopic)

	tests := []struct {
		name       string
		topicID    string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name:    "success",
			topicID: "1",
			body: map[string]interface{}{
				"title":   "Updated Topic",
				"content": "Updated Content",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "missing title",
			topicID: "1",
			body: map[string]interface{}{
				"content": "Updated Content",
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:    "missing content",
			topicID: "1",
			body: map[string]interface{}{
				"title": "Updated Topic",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPut, "/topics/"+tt.topicID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestTopicHandler_DeleteTopic(t *testing.T) {
	r, h := setupTestRouter()
	r.DELETE("/topics/:id", h.DeleteTopic)

	tests := []struct {
		name       string
		topicID    string
		wantStatus int
	}{
		{
			name:       "success",
			topicID:    "1",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "not found",
			topicID:    "999",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/topics/"+tt.topicID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
