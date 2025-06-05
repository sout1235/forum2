package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTopicRepository struct {
	mock.Mock
}

func (m *MockTopicRepository) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicRepository) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Topic), args.Error(1)
}

func (m *MockTopicRepository) GetAllTopics(ctx context.Context) ([]*entity.Topic, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Topic), args.Error(1)
}

func (m *MockTopicRepository) UpdateTopic(ctx context.Context, topic *entity.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicRepository) DeleteTopic(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTopicRepository) UpdateCommentCount(ctx context.Context, topicID int64) error {
	args := m.Called(ctx, topicID)
	return args.Error(0)
}

func TestTopicUseCase_CreateTopic(t *testing.T) {
	tests := []struct {
		name          string
		topic         *entity.Topic
		mockError     error
		expectedError error
	}{
		{
			name: "success",
			topic: &entity.Topic{
				Title:   "Test Topic",
				Content: "Test Content",
			},
		},
		{
			name: "empty title",
			topic: &entity.Topic{
				Title:   "",
				Content: "Test Content",
			},
			expectedError: errors.New("title and content are required"),
		},
		{
			name: "empty content",
			topic: &entity.Topic{
				Title:   "Test Topic",
				Content: "",
			},
			expectedError: errors.New("title and content are required"),
		},
		{
			name: "repository error",
			topic: &entity.Topic{
				Title:   "Test Topic",
				Content: "Test Content",
			},
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTopicRepo := new(MockTopicRepository)
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewTopicUseCase(mockTopicRepo, mockCommentRepo, mockUserRepo)

			if tt.expectedError == nil || tt.expectedError.Error() != "title and content are required" {
				mockTopicRepo.On("CreateTopic", mock.Anything, mock.AnythingOfType("*entity.Topic")).
					Return(tt.mockError)
			}

			err := uc.CreateTopic(context.Background(), tt.topic)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.topic.CreatedAt)
				assert.NotZero(t, tt.topic.UpdatedAt)
				assert.Zero(t, tt.topic.Views)
			}

			mockTopicRepo.AssertExpectations(t)
		})
	}
}

func TestTopicUseCase_GetTopicByID(t *testing.T) {
	tests := []struct {
		name          string
		topicID       int64
		mockTopic     *entity.Topic
		mockError     error
		expectedError error
	}{
		{
			name:    "success",
			topicID: 1,
			mockTopic: &entity.Topic{
				ID:      1,
				Title:   "Test Topic",
				Content: "Test Content",
			},
		},
		{
			name:          "repository error",
			topicID:       1,
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTopicRepo := new(MockTopicRepository)
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewTopicUseCase(mockTopicRepo, mockCommentRepo, mockUserRepo)

			mockTopicRepo.On("GetTopicByID", mock.Anything, tt.topicID).
				Return(tt.mockTopic, tt.mockError)

			topic, err := uc.GetTopicByID(context.Background(), tt.topicID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, topic)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, topic)
				assert.Equal(t, tt.mockTopic.ID, topic.ID)
				assert.Equal(t, tt.mockTopic.Title, topic.Title)
				assert.Equal(t, tt.mockTopic.Content, topic.Content)
			}

			mockTopicRepo.AssertExpectations(t)
		})
	}
}

func TestTopicUseCase_GetAllTopics(t *testing.T) {
	tests := []struct {
		name          string
		mockTopics    []*entity.Topic
		mockError     error
		mockUsernames map[int64]string
		mockUserError error
		expectedError error
	}{
		{
			name: "success",
			mockTopics: []*entity.Topic{
				{ID: 1, AuthorID: 1, Title: "Topic 1"},
				{ID: 2, AuthorID: 2, Title: "Topic 2"},
			},
			mockUsernames: map[int64]string{
				1: "user1",
				2: "user2",
			},
		},
		{
			name:          "repository error",
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
		{
			name: "user repository error",
			mockTopics: []*entity.Topic{
				{ID: 1, AuthorID: 1, Title: "Topic 1"},
			},
			mockUserError: errors.New("user error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTopicRepo := new(MockTopicRepository)
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewTopicUseCase(mockTopicRepo, mockCommentRepo, mockUserRepo)

			mockTopicRepo.On("GetAllTopics", mock.Anything).
				Return(tt.mockTopics, tt.mockError)

			if tt.mockError == nil {
				for _, topic := range tt.mockTopics {
					if username, ok := tt.mockUsernames[topic.AuthorID]; ok {
						mockUserRepo.On("GetUsernameByID", mock.Anything, topic.AuthorID).
							Return(username, tt.mockUserError)
					} else if tt.mockUserError != nil {
						mockUserRepo.On("GetUsernameByID", mock.Anything, topic.AuthorID).
							Return("", tt.mockUserError)
					}
				}
			}

			if tt.name == "user repository error" {
				for _, topic := range tt.mockTopics {
					mockUserRepo.On("GetUsernameByID", mock.Anything, topic.AuthorID).
						Return("", tt.mockUserError)
				}
			}

			topics, err := uc.GetAllTopics(context.Background())

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, topics)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, topics)
				assert.Equal(t, len(tt.mockTopics), len(topics))

				for _, topic := range topics {
					if tt.mockUserError == nil {
						assert.Equal(t, tt.mockUsernames[topic.AuthorID], topic.Author.Username)
					} else {
						assert.Equal(t, "User_1", topic.Author.Username)
					}
				}
			}

			mockTopicRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestTopicUseCase_UpdateTopic(t *testing.T) {
	tests := []struct {
		name          string
		topic         *entity.Topic
		mockError     error
		expectedError error
	}{
		{
			name: "success",
			topic: &entity.Topic{
				ID:      1,
				Title:   "Updated Topic",
				Content: "Updated Content",
			},
		},
		{
			name: "empty title",
			topic: &entity.Topic{
				ID:      1,
				Title:   "",
				Content: "Updated Content",
			},
			expectedError: errors.New("title and content are required"),
		},
		{
			name: "empty content",
			topic: &entity.Topic{
				ID:      1,
				Title:   "Updated Topic",
				Content: "",
			},
			expectedError: errors.New("title and content are required"),
		},
		{
			name: "repository error",
			topic: &entity.Topic{
				ID:      1,
				Title:   "Updated Topic",
				Content: "Updated Content",
			},
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTopicRepo := new(MockTopicRepository)
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewTopicUseCase(mockTopicRepo, mockCommentRepo, mockUserRepo)

			if tt.expectedError == nil || tt.expectedError.Error() != "title and content are required" {
				mockTopicRepo.On("UpdateTopic", mock.Anything, mock.AnythingOfType("*entity.Topic")).
					Return(tt.mockError)
			}

			err := uc.UpdateTopic(context.Background(), tt.topic)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.topic.UpdatedAt)
			}

			mockTopicRepo.AssertExpectations(t)
		})
	}
}

func TestTopicUseCase_DeleteTopic(t *testing.T) {
	tests := []struct {
		name          string
		topicID       int64
		mockError     error
		expectedError error
	}{
		{
			name:    "success",
			topicID: 1,
		},
		{
			name:          "repository error",
			topicID:       1,
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTopicRepo := new(MockTopicRepository)
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewTopicUseCase(mockTopicRepo, mockCommentRepo, mockUserRepo)

			mockTopicRepo.On("DeleteTopic", mock.Anything, tt.topicID).
				Return(tt.mockError)

			err := uc.DeleteTopic(context.Background(), tt.topicID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockTopicRepo.AssertExpectations(t)
		})
	}
}

func TestTopicUseCase_IncrementViews(t *testing.T) {
	tests := []struct {
		name          string
		topicID       int64
		mockTopic     *entity.Topic
		mockError     error
		expectedError error
	}{
		{
			name:    "success",
			topicID: 1,
			mockTopic: &entity.Topic{
				ID:    1,
				Views: 5,
			},
		},
		{
			name:          "get topic error",
			topicID:       1,
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
		{
			name:    "update error",
			topicID: 1,
			mockTopic: &entity.Topic{
				ID:    1,
				Views: 5,
			},
			mockError:     errors.New("update error"),
			expectedError: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTopicRepo := new(MockTopicRepository)
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewTopicUseCase(mockTopicRepo, mockCommentRepo, mockUserRepo)

			mockTopicRepo.On("GetTopicByID", mock.Anything, tt.topicID).
				Return(tt.mockTopic, tt.mockError)

			if tt.mockError == nil {
				mockTopicRepo.On("UpdateTopic", mock.Anything, mock.AnythingOfType("*entity.Topic")).
					Return(tt.mockError)
			}

			oldViews := 0
			if tt.mockTopic != nil {
				oldViews = tt.mockTopic.Views
			}

			err := uc.IncrementViews(context.Background(), tt.topicID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, oldViews+1, tt.mockTopic.Views)
			}

			mockTopicRepo.AssertExpectations(t)
		})
	}
}
