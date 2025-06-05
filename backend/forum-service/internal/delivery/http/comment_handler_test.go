package httpDelivery

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommentUseCase struct {
	mock.Mock
}

func (m *MockCommentUseCase) GetCommentsByTopicID(ctx context.Context, topicID int64) ([]*entity.Comment, error) {
	args := m.Called(ctx, topicID)
	return args.Get(0).([]*entity.Comment), args.Error(1)
}

func (m *MockCommentUseCase) GetCommentByID(ctx context.Context, commentID int64) (*entity.Comment, error) {
	args := m.Called(ctx, commentID)
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *MockCommentUseCase) CreateComment(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentUseCase) DeleteComment(ctx context.Context, commentID int64) error {
	args := m.Called(ctx, commentID)
	return args.Error(0)
}

func (m *MockCommentUseCase) LikeComment(ctx context.Context, commentID int64) error {
	args := m.Called(ctx, commentID)
	return args.Error(0)
}

func TestCommentHandler_GetAllCommentsByTopic(t *testing.T) {
	muc := new(MockCommentUseCase)
	h := NewCommentHandler(muc, nil)
	r, _ := setupTestRouter()
	r.GET("/topics/:id/comments", h.GetAllCommentsByTopic)
	timeNow := time.Now()
	comments := []*entity.Comment{{ID: 1, Content: "c1", AuthorID: 2, TopicID: 3, CreatedAt: timeNow, UpdatedAt: timeNow}}
	muc.On("GetCommentsByTopicID", mock.Anything, int64(3)).Return(comments, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics/3/comments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error case
	muc.On("GetCommentsByTopicID", mock.Anything, int64(99)).Return([]*entity.Comment{}, errors.New("fail"))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/topics/99/comments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// invalid id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/topics/bad/comments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_GetComment(t *testing.T) {
	muc := new(MockCommentUseCase)
	h := NewCommentHandler(muc, nil)
	r, _ := setupTestRouter()
	r.GET("/comments/:id", h.GetComment)
	timeNow := time.Now()
	comment := &entity.Comment{ID: 1, Content: "c1", AuthorID: 2, TopicID: 3, CreatedAt: timeNow, UpdatedAt: timeNow}
	muc.On("GetCommentByID", mock.Anything, int64(1)).Return(comment, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/comments/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error case
	muc.On("GetCommentByID", mock.Anything, int64(99)).Return(&entity.Comment{}, errors.New("fail"))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/comments/99", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// invalid id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/comments/bad", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_CreateComment(t *testing.T) {
	muc := new(MockCommentUseCase)
	h := NewCommentHandler(muc, nil)
	r, _ := setupTestRouter()
	r.POST("/topics/:id/comments", h.CreateComment)
	muc.On("CreateComment", mock.Anything, mock.AnythingOfType("*entity.Comment")).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/topics/3/comments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// error case
	muc.On("CreateComment", mock.Anything, mock.AnythingOfType("*entity.Comment")).Return(errors.New("fail"))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/topics/3/comments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_DeleteComment(t *testing.T) {
	muc := new(MockCommentUseCase)
	h := NewCommentHandler(muc, nil)
	r, _ := setupTestRouter()
	r.DELETE("/topics/:topic_id/comments/:id", h.DeleteComment)
	muc.On("DeleteComment", mock.Anything, int64(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/topics/3/comments/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error case
	muc.On("DeleteComment", mock.Anything, int64(99)).Return(errors.New("fail"))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/topics/3/comments/99", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// invalid id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/topics/3/comments/bad", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_LikeComment(t *testing.T) {
	muc := new(MockCommentUseCase)
	h := NewCommentHandler(muc, nil)
	r, _ := setupTestRouter()
	r.POST("/comments/:id/like", h.LikeComment)
	muc.On("LikeComment", mock.Anything, int64(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comments/1/like", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error case
	muc.On("LikeComment", mock.Anything, int64(99)).Return(errors.New("fail"))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/comments/99/like", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// invalid id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/comments/bad/like", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
