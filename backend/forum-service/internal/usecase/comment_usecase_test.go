package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) GetCommentsByTopic(ctx context.Context, topicID int64) ([]*entity.Comment, error) {
	args := m.Called(ctx, topicID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *MockCommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) DeleteComment(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) LikeComment(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsernameByID(ctx context.Context, id int64) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestCommentUseCase_GetCommentsByTopic(t *testing.T) {
	tests := []struct {
		name          string
		topicID       int64
		mockComments  []*entity.Comment
		mockError     error
		mockUsername  string
		mockUserError error
		expectedError error
	}{
		{
			name:    "success",
			topicID: 1,
			mockComments: []*entity.Comment{
				{ID: 1, TopicID: 1, AuthorID: 1, Content: "Test comment"},
			},
			mockUsername: "testuser",
		},
		{
			name:          "repository error",
			topicID:       1,
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
		{
			name:    "user repository error",
			topicID: 1,
			mockComments: []*entity.Comment{
				{ID: 1, TopicID: 1, AuthorID: 1, Content: "Test comment"},
			},
			mockUserError: errors.New("user error"),
			expectedError: errors.New("user error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewCommentUseCase(mockCommentRepo, mockUserRepo)

			mockCommentRepo.On("GetCommentsByTopic", mock.Anything, tt.topicID).
				Return(tt.mockComments, tt.mockError)

			if tt.mockError == nil {
				mockUserRepo.On("GetUsernameByID", mock.Anything, int64(1)).
					Return(tt.mockUsername, tt.mockUserError)
			}

			comments, err := uc.GetCommentsByTopicID(context.Background(), tt.topicID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, comments)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comments)
				if len(comments) > 0 {
					assert.Equal(t, tt.mockUsername, comments[0].Author.Username)
				}
			}

			mockCommentRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestCommentUseCase_GetCommentByID(t *testing.T) {
	tests := []struct {
		name          string
		commentID     int64
		mockComment   *entity.Comment
		mockError     error
		mockUsername  string
		mockUserError error
		expectedError error
	}{
		{
			name:         "success",
			commentID:    1,
			mockComment:  &entity.Comment{ID: 1, TopicID: 1, AuthorID: 1, Content: "Test comment"},
			mockUsername: "testuser",
		},
		{
			name:          "repository error",
			commentID:     1,
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
		{
			name:          "user repository error",
			commentID:     1,
			mockComment:   &entity.Comment{ID: 1, TopicID: 1, AuthorID: 1, Content: "Test comment"},
			mockUserError: errors.New("user error"),
			expectedError: errors.New("user error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewCommentUseCase(mockCommentRepo, mockUserRepo)

			mockCommentRepo.On("GetCommentByID", mock.Anything, tt.commentID).
				Return(tt.mockComment, tt.mockError)

			if tt.mockError == nil {
				mockUserRepo.On("GetUsernameByID", mock.Anything, int64(1)).
					Return(tt.mockUsername, tt.mockUserError)
			}

			comment, err := uc.GetCommentByID(context.Background(), tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tt.mockUsername, comment.Author.Username)
			}

			mockCommentRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestCommentUseCase_CreateComment(t *testing.T) {
	tests := []struct {
		name          string
		comment       *entity.Comment
		mockError     error
		mockUsername  string
		mockUserError error
		expectedError error
	}{
		{
			name: "success",
			comment: &entity.Comment{
				TopicID:  1,
				AuthorID: 1,
				Content:  "Test comment",
			},
			mockUsername: "testuser",
		},
		{
			name: "repository error",
			comment: &entity.Comment{
				TopicID:  1,
				AuthorID: 1,
				Content:  "Test comment",
			},
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
		{
			name: "user repository error",
			comment: &entity.Comment{
				TopicID:  1,
				AuthorID: 1,
				Content:  "Test comment",
			},
			mockUserError: errors.New("user error"),
			expectedError: errors.New("user error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewCommentUseCase(mockCommentRepo, mockUserRepo)

			mockCommentRepo.On("CreateComment", mock.Anything, tt.comment).
				Return(tt.mockError)

			if tt.mockError == nil {
				mockUserRepo.On("GetUsernameByID", mock.Anything, int64(1)).
					Return(tt.mockUsername, tt.mockUserError)
			}

			err := uc.CreateComment(context.Background(), tt.comment)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUsername, tt.comment.Author.Username)
			}

			mockCommentRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestCommentUseCase_DeleteComment(t *testing.T) {
	tests := []struct {
		name          string
		commentID     int64
		mockError     error
		expectedError error
	}{
		{
			name:      "success",
			commentID: 1,
		},
		{
			name:          "repository error",
			commentID:     1,
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewCommentUseCase(mockCommentRepo, mockUserRepo)

			mockCommentRepo.On("DeleteComment", mock.Anything, tt.commentID).
				Return(tt.mockError)

			err := uc.DeleteComment(context.Background(), tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockCommentRepo.AssertExpectations(t)
		})
	}
}

func TestCommentUseCase_LikeComment(t *testing.T) {
	tests := []struct {
		name          string
		commentID     int64
		mockError     error
		expectedError error
	}{
		{
			name:      "success",
			commentID: 1,
		},
		{
			name:          "repository error",
			commentID:     1,
			mockError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentRepo := new(MockCommentRepository)
			mockUserRepo := new(MockUserRepository)
			uc := NewCommentUseCase(mockCommentRepo, mockUserRepo)

			mockCommentRepo.On("LikeComment", mock.Anything, tt.commentID).
				Return(tt.mockError)

			err := uc.LikeComment(context.Background(), tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockCommentRepo.AssertExpectations(t)
		})
	}
}
