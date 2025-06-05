package service

import (
	"context"
	"testing"
	"time"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChatRepo struct {
	mock.Mock
}

func (m *mockChatRepo) SaveMessage(ctx context.Context, message *entity.ChatMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *mockChatRepo) GetRecentMessages(ctx context.Context, limit int) ([]*entity.ChatMessage, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ChatMessage), args.Error(1)
}

func (m *mockChatRepo) DeleteExpiredMessages(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestChatService_SaveMessage(t *testing.T) {
	mockRepo := new(mockChatRepo)
	service := NewChatService(mockRepo)

	message := &entity.ChatMessage{
		Content:        "Test Message",
		AuthorID:       1,
		AuthorUsername: "testuser",
		CreatedAt:      time.Now(),
	}

	mockRepo.On("SaveMessage", mock.Anything, message).Return(nil)

	err := service.SaveMessage(context.Background(), message)
	assert.NoError(t, err)
	assert.NotZero(t, message.ExpiresAt)
	mockRepo.AssertExpectations(t)
}

func TestChatService_SaveMessage_WithExpiresAt(t *testing.T) {
	mockRepo := new(mockChatRepo)
	service := NewChatService(mockRepo)

	expiresAt := time.Now().Add(2 * time.Hour)
	message := &entity.ChatMessage{
		Content:        "Test Message",
		AuthorID:       1,
		AuthorUsername: "testuser",
		CreatedAt:      time.Now(),
		ExpiresAt:      expiresAt,
	}

	mockRepo.On("SaveMessage", mock.Anything, message).Return(nil)

	err := service.SaveMessage(context.Background(), message)
	assert.NoError(t, err)
	assert.Equal(t, expiresAt, message.ExpiresAt)
	mockRepo.AssertExpectations(t)
}

func TestChatService_GetRecentMessages(t *testing.T) {
	mockRepo := new(mockChatRepo)
	service := NewChatService(mockRepo)

	expectedMessages := []*entity.ChatMessage{
		{
			ID:             1,
			Content:        "Test Message 1",
			AuthorID:       1,
			AuthorUsername: "testuser",
			CreatedAt:      time.Now(),
			ExpiresAt:      time.Now().Add(time.Hour),
		},
		{
			ID:             2,
			Content:        "Test Message 2",
			AuthorID:       2,
			AuthorUsername: "testuser2",
			CreatedAt:      time.Now(),
			ExpiresAt:      time.Now().Add(time.Hour),
		},
	}

	mockRepo.On("GetRecentMessages", mock.Anything, 10).Return(expectedMessages, nil)

	messages, err := service.GetRecentMessages(context.Background(), 10)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestChatService_DeleteExpiredMessages(t *testing.T) {
	mockRepo := new(mockChatRepo)
	service := NewChatService(mockRepo)

	mockRepo.On("DeleteExpiredMessages", mock.Anything).Return(nil)

	err := service.DeleteExpiredMessages(context.Background())
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
