package service

import (
	"context"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
)

type TopicService interface {
	CreateTopic(ctx context.Context, topic *entity.Topic) error
	GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error)
	GetAllTopics(ctx context.Context) ([]*entity.Topic, error)
	UpdateTopic(ctx context.Context, topic *entity.Topic) error
	DeleteTopic(ctx context.Context, id int64) error
	UpdateCommentCount(ctx context.Context, topicID int64) error
}

type topicService struct {
	topicRepo repository.TopicRepository
	userRepo  repository.UserRepository
}

// NewTopicService creates a new instance of TopicService
func NewTopicService(topicRepo repository.TopicRepository, userRepo repository.UserRepository) TopicService {
	return &topicService{
		topicRepo: topicRepo,
		userRepo:  userRepo,
	}
}

func (s *topicService) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	return s.topicRepo.CreateTopic(ctx, topic)
}

func (s *topicService) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	topic, err := s.topicRepo.GetTopicByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get author information
	if topic.AuthorID > 0 {
		author, err := s.userRepo.GetUserByID(ctx, topic.AuthorID)
		if err != nil {
			return nil, err
		}
		topic.Author = author
	}

	return topic, nil
}

func (s *topicService) GetAllTopics(ctx context.Context) ([]*entity.Topic, error) {
	topics, err := s.topicRepo.GetAllTopics(ctx)
	if err != nil {
		return nil, err
	}

	// Get author information for each topic that doesn't have it
	for _, topic := range topics {
		if topic.AuthorID > 0 && (topic.Author == nil || topic.Author.Username == "") {
			author, err := s.userRepo.GetUserByID(ctx, topic.AuthorID)
			if err != nil {
				return nil, err
			}
			topic.Author = author
		}
	}

	return topics, nil
}

func (s *topicService) UpdateTopic(ctx context.Context, topic *entity.Topic) error {
	return s.topicRepo.UpdateTopic(ctx, topic)
}

func (s *topicService) DeleteTopic(ctx context.Context, id int64) error {
	return s.topicRepo.DeleteTopic(ctx, id)
}

func (s *topicService) UpdateCommentCount(ctx context.Context, topicID int64) error {
	return s.topicRepo.UpdateCommentCount(ctx, topicID)
}
