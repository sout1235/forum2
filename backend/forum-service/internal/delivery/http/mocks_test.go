package httpDelivery

import (
	"context"
	"time"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/mock"
)

// MockCommentRepository представляет мок для CommentRepository
type MockCommentRepository struct{}

func (m *MockCommentRepository) GetCommentsByTopic(_ context.Context, topicID int64) ([]*entity.Comment, error) {
	now := time.Now()
	return []*entity.Comment{
		{
			ID:        1,
			Content:   "Test comment",
			AuthorID:  42,
			TopicID:   topicID,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func (m *MockCommentRepository) GetCommentByID(_ context.Context, id int64) (*entity.Comment, error) {
	return nil, nil
}

func (m *MockCommentRepository) CreateComment(_ context.Context, c *entity.Comment) error {
	return nil
}

func (m *MockCommentRepository) UpdateComment(_ context.Context, c *entity.Comment) error {
	return nil
}

func (m *MockCommentRepository) DeleteComment(_ context.Context, id int64) error {
	return nil
}

func (m *MockCommentRepository) LikeComment(_ context.Context, id int64) error {
	return nil
}

// MockTopicRepository представляет мок для TopicRepository
type MockTopicRepository struct{}

func (m *MockTopicRepository) CreateTopic(_ context.Context, topic *entity.Topic) error {
	topic.ID = 1
	return nil
}

func (m *MockTopicRepository) GetTopicByID(_ context.Context, id int64) (*entity.Topic, error) {
	if id == 1 {
		return &entity.Topic{
			ID:        1,
			Title:     "Test Topic",
			Content:   "Test Content",
			AuthorID:  42,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	return nil, nil
}

func (m *MockTopicRepository) GetAllTopics(_ context.Context) ([]*entity.Topic, error) {
	return []*entity.Topic{
		{
			ID:        1,
			Title:     "Test Topic 1",
			Content:   "Test Content 1",
			AuthorID:  42,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "Test Topic 2",
			Content:   "Test Content 2",
			AuthorID:  43,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

func (m *MockTopicRepository) UpdateTopic(_ context.Context, topic *entity.Topic) error {
	return nil
}

func (m *MockTopicRepository) DeleteTopic(_ context.Context, id int64) error {
	return nil
}

func (m *MockTopicRepository) UpdateCommentCount(_ context.Context, topicID int64) error {
	return nil
}

// MockUserRepository представляет мок для UserRepository
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
