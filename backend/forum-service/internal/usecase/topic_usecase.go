package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
)

type TopicUseCase struct {
	topicRepo   repository.TopicRepository
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
}

func NewTopicUseCase(topicRepo repository.TopicRepository, commentRepo repository.CommentRepository, userRepo repository.UserRepository) *TopicUseCase {
	return &TopicUseCase{
		topicRepo:   topicRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

func (uc *TopicUseCase) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	if topic.Title == "" || topic.Content == "" {
		return errors.New("title and content are required")
	}

	now := time.Now()
	topic.CreatedAt = now
	topic.UpdatedAt = now
	topic.Views = 0

	return uc.topicRepo.CreateTopic(ctx, topic)
}

func (uc *TopicUseCase) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	return uc.topicRepo.GetTopicByID(ctx, id)
}

func (uc *TopicUseCase) GetAllTopics(ctx context.Context) ([]*entity.Topic, error) {
	topics, err := uc.topicRepo.GetAllTopics(ctx)
	if err != nil {
		return nil, err
	}

	// Получаем имена пользователей для каждого топика
	for _, topic := range topics {
		username, err := uc.userRepo.GetUsernameByID(ctx, topic.AuthorID)
		if err != nil {
			log.Printf("Error getting username for user %d: %v", topic.AuthorID, err)
			username = fmt.Sprintf("User_%d", topic.AuthorID)
		}
		topic.Author = &entity.User{
			ID:       topic.AuthorID,
			Username: username,
		}
	}

	return topics, nil
}

func (uc *TopicUseCase) UpdateTopic(ctx context.Context, topic *entity.Topic) error {
	if topic.Title == "" || topic.Content == "" {
		return errors.New("title and content are required")
	}

	topic.UpdatedAt = time.Now()
	return uc.topicRepo.UpdateTopic(ctx, topic)
}

func (uc *TopicUseCase) DeleteTopic(ctx context.Context, id int64) error {
	return uc.topicRepo.DeleteTopic(ctx, id)
}

func (uc *TopicUseCase) IncrementViews(ctx context.Context, id int64) error {
	topic, err := uc.topicRepo.GetTopicByID(ctx, id)
	if err != nil {
		return err
	}

	topic.Views++
	return uc.topicRepo.UpdateTopic(ctx, topic)
}

func (uc *TopicUseCase) UpdateCommentCount(ctx context.Context, topicID int64) error {
	return uc.topicRepo.UpdateCommentCount(ctx, topicID)
}
