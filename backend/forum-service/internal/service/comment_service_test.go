package service

import (
	"context"
	"testing"
	"time"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCommentRepo struct {
	mock.Mock
}

func (m *mockCommentRepo) CreateComment(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *mockCommentRepo) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *mockCommentRepo) GetCommentsByTopic(ctx context.Context, topicID int64) ([]*entity.Comment, error) {
	args := m.Called(ctx, topicID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Comment), args.Error(1)
}

func (m *mockCommentRepo) DeleteComment(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockCommentRepo) LikeComment(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockCommentRepo) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

type mockTopicRepoForComment struct {
	mock.Mock
}

func (m *mockTopicRepoForComment) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *mockTopicRepoForComment) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Topic), args.Error(1)
}

func (m *mockTopicRepoForComment) GetAllTopics(ctx context.Context) ([]*entity.Topic, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Topic), args.Error(1)
}

func (m *mockTopicRepoForComment) UpdateTopic(ctx context.Context, topic *entity.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *mockTopicRepoForComment) DeleteTopic(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTopicRepoForComment) UpdateCommentCount(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCommentService_CreateComment(t *testing.T) {
	mockCommentRepo := new(mockCommentRepo)
	mockTopicRepo := new(mockTopicRepoForComment)
	service := NewCommentService(mockCommentRepo, mockTopicRepo)

	comment := &entity.Comment{
		Content:   "Test Comment",
		AuthorID:  1,
		TopicID:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Likes:     0,
		Author: &entity.User{
			ID:       1,
			Username: "testuser",
		},
	}

	mockCommentRepo.On("CreateComment", mock.Anything, comment).Return(nil)
	mockTopicRepo.On("UpdateCommentCount", mock.Anything, int64(1)).Return(nil)

	err := service.CreateComment(context.Background(), comment)
	assert.NoError(t, err)
	mockCommentRepo.AssertExpectations(t)
	mockTopicRepo.AssertExpectations(t)
}

func TestCommentService_GetCommentByID(t *testing.T) {
	mockCommentRepo := new(mockCommentRepo)
	mockTopicRepo := new(mockTopicRepoForComment)
	service := NewCommentService(mockCommentRepo, mockTopicRepo)

	expectedComment := &entity.Comment{
		ID:        1,
		Content:   "Test Comment",
		AuthorID:  1,
		TopicID:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Likes:     0,
		Author: &entity.User{
			ID:       1,
			Username: "testuser",
		},
	}

	mockCommentRepo.On("GetCommentByID", mock.Anything, int64(1)).Return(expectedComment, nil)

	comment, err := service.GetCommentByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_GetCommentsByTopic(t *testing.T) {
	mockCommentRepo := new(mockCommentRepo)
	mockTopicRepo := new(mockTopicRepoForComment)
	service := NewCommentService(mockCommentRepo, mockTopicRepo)

	expectedComments := []*entity.Comment{
		{
			ID:        1,
			Content:   "Test Comment 1",
			AuthorID:  1,
			TopicID:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Likes:     0,
			Author: &entity.User{
				ID:       1,
				Username: "testuser",
			},
		},
		{
			ID:        2,
			Content:   "Test Comment 2",
			AuthorID:  2,
			TopicID:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Likes:     0,
			Author: &entity.User{
				ID:       2,
				Username: "testuser2",
			},
		},
	}

	mockCommentRepo.On("GetCommentsByTopic", mock.Anything, int64(1)).Return(expectedComments, nil)

	comments, err := service.GetCommentsByTopicID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_DeleteComment(t *testing.T) {
	mockCommentRepo := new(mockCommentRepo)
	mockTopicRepo := new(mockTopicRepoForComment)
	service := NewCommentService(mockCommentRepo, mockTopicRepo)

	comment := &entity.Comment{
		ID:        1,
		Content:   "Test Comment",
		AuthorID:  1,
		TopicID:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Likes:     0,
		Author: &entity.User{
			ID:       1,
			Username: "testuser",
		},
	}

	mockCommentRepo.On("GetCommentByID", mock.Anything, int64(1)).Return(comment, nil)
	mockCommentRepo.On("DeleteComment", mock.Anything, int64(1)).Return(nil)
	mockTopicRepo.On("UpdateCommentCount", mock.Anything, int64(1)).Return(nil)

	err := service.DeleteComment(context.Background(), 1)
	assert.NoError(t, err)
	mockCommentRepo.AssertExpectations(t)
	mockTopicRepo.AssertExpectations(t)
}

func TestCommentService_LikeComment(t *testing.T) {
	mockCommentRepo := new(mockCommentRepo)
	mockTopicRepo := new(mockTopicRepoForComment)
	service := NewCommentService(mockCommentRepo, mockTopicRepo)

	mockCommentRepo.On("LikeComment", mock.Anything, int64(1)).Return(nil)

	err := service.LikeComment(context.Background(), 1)
	assert.NoError(t, err)
	mockCommentRepo.AssertExpectations(t)
}
