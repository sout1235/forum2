package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	httpDelivery "github.com/sout1235/forum2/backend/forum-service/internal/delivery/http"
	"github.com/sout1235/forum2/backend/forum-service/internal/delivery/http/middleware"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
	"github.com/sout1235/forum2/backend/forum-service/internal/service"
	"github.com/sout1235/forum2/backend/forum-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	router *httpDelivery.Router
	token  string
)

func setupTestEnvironment() (*httpDelivery.Router, error) {
	// Initialize database connection
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=forum_test sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Initialize repositories
	topicRepo := repository.NewTopicRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	userRepo := repository.NewUserRepository(db, "http://localhost:8080")
	chatRepo := repository.NewChatRepository(db)

	// Initialize services and use cases
	topicService := service.NewTopicService(topicRepo, userRepo)
	commentUseCase := usecase.NewCommentUseCase(commentRepo, userRepo)

	// Initialize router
	router := httpDelivery.NewRouter(
		topicService,
		commentUseCase,
		userRepo,
		chatRepo,
		"http://localhost:8080",
		&middleware.AuthConfig{
			AuthServiceURL: "http://localhost:8080",
		},
	)

	return router, nil
}

func cleanup() error {
	// Add cleanup logic here if needed
	return nil
}

func TestMain(m *testing.M) {
	var err error
	router, err = setupTestEnvironment()
	if err != nil {
		log.Fatalf("Failed to setup test environment: %v", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err := cleanup(); err != nil {
		log.Printf("Failed to cleanup: %v", err)
	}

	os.Exit(code)
}

func TestCreateAndGetTopic(t *testing.T) {
	// Create topic
	topicData := map[string]interface{}{
		"title":       "Test Topic",
		"content":     "Test Content",
		"category_id": 1,
		"tags":        []string{"test", "integration"},
	}

	jsonData, _ := json.Marshal(topicData)
	req := httptest.NewRequest("POST", "/api/v1/topics", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	topicID := response["id"].(float64)

	// Get created topic
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/topics/%d", int(topicID)), nil)
	w = httptest.NewRecorder()
	router.Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var topic map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &topic)
	require.NoError(t, err)

	assert.Equal(t, "Test Topic", topic["title"])
	assert.Equal(t, "Test Content", topic["content"])
}

func TestCreateAndGetComment(t *testing.T) {
	// Create topic for comment
	topicData := map[string]interface{}{
		"title":       "Topic for Comment",
		"content":     "Topic Content",
		"category_id": 1,
		"tags":        []string{"test"},
	}

	jsonData, _ := json.Marshal(topicData)
	req := httptest.NewRequest("POST", "/api/v1/topics", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.Engine().ServeHTTP(w, req)

	var topicResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &topicResponse)
	topicID := topicResponse["id"].(float64)

	// Create comment
	commentData := map[string]interface{}{
		"content": "Test Comment",
	}

	jsonData, _ = json.Marshal(commentData)
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/v1/topics/%d/comments", int(topicID)), bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var commentResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &commentResponse)
	require.NoError(t, err)

	commentID := commentResponse["id"].(float64)

	// Get comment
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/comments/%d", int(commentID)), nil)
	w = httptest.NewRecorder()
	router.Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var comment map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &comment)
	require.NoError(t, err)

	assert.Equal(t, "Test Comment", comment["content"])
}

func TestTopicWithComments(t *testing.T) {
	// Create topic
	topicData := map[string]interface{}{
		"title":       "Topic with Comments",
		"content":     "Topic Content",
		"category_id": 1,
		"tags":        []string{"test"},
	}

	jsonData, _ := json.Marshal(topicData)
	req := httptest.NewRequest("POST", "/api/v1/topics", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.Engine().ServeHTTP(w, req)

	var topicResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &topicResponse)
	topicID := topicResponse["id"].(float64)

	// Create multiple comments
	comments := []string{"First comment", "Second comment", "Third comment"}
	for _, content := range comments {
		commentData := map[string]interface{}{
			"content": content,
		}

		jsonData, _ = json.Marshal(commentData)
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/v1/topics/%d/comments", int(topicID)), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.Engine().ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Get all comments for topic
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/topics/%d/comments", int(topicID)), nil)
	w = httptest.NewRecorder()
	router.Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var commentsResponse []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &commentsResponse)
	require.NoError(t, err)

	assert.Equal(t, len(comments), len(commentsResponse))
	for i, comment := range commentsResponse {
		assert.Equal(t, comments[i], comment["content"])
	}
}
