package service

import (
	"context"
	"testing"
	"time"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTopicRepo struct {
	mock.Mock
}

func (m *mockTopicRepo) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *mockTopicRepo) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Topic), args.Error(1)
}

func (m *mockTopicRepo) GetAllTopics(ctx context.Context) ([]*entity.Topic, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Topic), args.Error(1)
}

func (m *mockTopicRepo) UpdateTopic(ctx context.Context, topic *entity.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *mockTopicRepo) DeleteTopic(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTopicRepo) UpdateCommentCount(ctx context.Context, topicID int64) error {
	args := m.Called(ctx, topicID)
	return args.Error(0)
}

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) GetUsernameByID(ctx context.Context, id int64) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestTopicService_CreateTopic(t *testing.T) {
	mockTopicRepo := new(mockTopicRepo)
	mockUserRepo := new(mockUserRepo)
	topicService := NewTopicService(mockTopicRepo, mockUserRepo)

	topic := &entity.Topic{
		Title:        "Test Topic",
		Content:      "Test Content",
		AuthorID:     1,
		CategoryID:   1,
		Views:        0,
		CommentCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockTopicRepo.On("CreateTopic", mock.Anything, topic).Return(nil)

	err := topicService.CreateTopic(context.Background(), topic)
	assert.NoError(t, err)
	mockTopicRepo.AssertExpectations(t)
}

func TestTopicService_GetTopicByID(t *testing.T) {
	mockTopicRepo := new(mockTopicRepo)
	mockUserRepo := new(mockUserRepo)
	topicService := NewTopicService(mockTopicRepo, mockUserRepo)

	expectedTopic := &entity.Topic{
		ID:           1,
		Title:        "Test Topic",
		Content:      "Test Content",
		AuthorID:     1,
		CategoryID:   1,
		Views:        0,
		CommentCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockTopicRepo.On("GetTopicByID", mock.Anything, int64(1)).Return(expectedTopic, nil)

	topic, err := topicService.GetTopicByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedTopic, topic)
	mockTopicRepo.AssertExpectations(t)
}

func TestTopicService_GetAllTopics(t *testing.T) {
	mockTopicRepo := new(mockTopicRepo)
	mockUserRepo := new(mockUserRepo)
	topicService := NewTopicService(mockTopicRepo, mockUserRepo)

	expectedTopics := []*entity.Topic{
		{
			ID:           1,
			Title:        "Test Topic 1",
			Content:      "Test Content 1",
			AuthorID:     1,
			CategoryID:   1,
			Views:        0,
			CommentCount: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           2,
			Title:        "Test Topic 2",
			Content:      "Test Content 2",
			AuthorID:     2,
			CategoryID:   1,
			Views:        0,
			CommentCount: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	mockTopicRepo.On("GetAllTopics", mock.Anything).Return(expectedTopics, nil)

	topics, err := topicService.GetAllTopics(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedTopics, topics)
	mockTopicRepo.AssertExpectations(t)
}

func TestTopicService_UpdateTopic(t *testing.T) {
	mockTopicRepo := new(mockTopicRepo)
	mockUserRepo := new(mockUserRepo)
	topicService := NewTopicService(mockTopicRepo, mockUserRepo)

	topic := &entity.Topic{
		ID:           1,
		Title:        "Updated Topic",
		Content:      "Updated Content",
		AuthorID:     1,
		CategoryID:   2,
		Views:        0,
		CommentCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockTopicRepo.On("UpdateTopic", mock.Anything, topic).Return(nil)

	err := topicService.UpdateTopic(context.Background(), topic)
	assert.NoError(t, err)
	mockTopicRepo.AssertExpectations(t)
}

func TestTopicService_DeleteTopic(t *testing.T) {
	mockTopicRepo := new(mockTopicRepo)
	mockUserRepo := new(mockUserRepo)
	topicService := NewTopicService(mockTopicRepo, mockUserRepo)

	mockTopicRepo.On("DeleteTopic", mock.Anything, int64(1)).Return(nil)

	err := topicService.DeleteTopic(context.Background(), 1)
	assert.NoError(t, err)
	mockTopicRepo.AssertExpectations(t)
}

func TestTopicService_UpdateCommentCount(t *testing.T) {
	mockTopicRepo := new(mockTopicRepo)
	mockUserRepo := new(mockUserRepo)
	topicService := NewTopicService(mockTopicRepo, mockUserRepo)

	mockTopicRepo.On("UpdateCommentCount", mock.Anything, int64(1)).Return(nil)

	err := topicService.UpdateCommentCount(context.Background(), 1)
	assert.NoError(t, err)
	mockTopicRepo.AssertExpectations(t)
}
