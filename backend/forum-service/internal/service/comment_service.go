package service

import (
	"context"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
)

type commentService struct {
	commentRepo repository.CommentRepository
	topicRepo   repository.TopicRepository
}

// NewCommentService creates a new instance of CommentService
func NewCommentService(commentRepo repository.CommentRepository, topicRepo repository.TopicRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		topicRepo:   topicRepo,
	}
}

func (s *commentService) CreateComment(ctx context.Context, comment *entity.Comment) error {
	if err := s.commentRepo.CreateComment(ctx, comment); err != nil {
		return err
	}
	return s.topicRepo.UpdateCommentCount(ctx, comment.TopicID)
}

func (s *commentService) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	return s.commentRepo.GetCommentByID(ctx, id)
}

func (s *commentService) GetCommentsByTopicID(ctx context.Context, topicID int64) ([]*entity.Comment, error) {
	return s.commentRepo.GetCommentsByTopic(ctx, topicID)
}

func (s *commentService) DeleteComment(ctx context.Context, id int64) error {
	comment, err := s.commentRepo.GetCommentByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.commentRepo.DeleteComment(ctx, id); err != nil {
		return err
	}
	return s.topicRepo.UpdateCommentCount(ctx, comment.TopicID)
}

func (s *commentService) LikeComment(ctx context.Context, id int64) error {
	return s.commentRepo.LikeComment(ctx, id)
}

func (s *commentService) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	return s.commentRepo.UpdateComment(ctx, comment)
}
