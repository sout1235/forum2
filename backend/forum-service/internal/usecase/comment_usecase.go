package usecase

import (
	"context"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
)

type CommentUseCase interface {
	GetCommentsByTopicID(ctx context.Context, topicID int64) ([]*entity.Comment, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
	CreateComment(ctx context.Context, comment *entity.Comment) error
	DeleteComment(ctx context.Context, id int64) error
	LikeComment(ctx context.Context, id int64) error
}

type commentUseCase struct {
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
}

func NewCommentUseCase(commentRepo repository.CommentRepository, userRepo repository.UserRepository) CommentUseCase {
	return &commentUseCase{
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

func (uc *commentUseCase) GetCommentsByTopicID(ctx context.Context, topicID int64) ([]*entity.Comment, error) {
	comments, err := uc.commentRepo.GetCommentsByTopic(ctx, topicID)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if comment.AuthorID > 0 && (comment.Author == nil || comment.Author.Username == "") {
			author, err := uc.userRepo.GetUserByID(ctx, comment.AuthorID)
			if err != nil {
				return nil, err
			}
			comment.Author = author
		}
	}

	return comments, nil
}

func (uc *commentUseCase) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	comment, err := uc.commentRepo.GetCommentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if comment.AuthorID > 0 && (comment.Author == nil || comment.Author.Username == "") {
		author, err := uc.userRepo.GetUserByID(ctx, comment.AuthorID)
		if err != nil {
			return nil, err
		}
		comment.Author = author
	}

	return comment, nil
}

func (uc *commentUseCase) CreateComment(ctx context.Context, comment *entity.Comment) error {
	err := uc.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return err
	}

	if comment.AuthorID > 0 && (comment.Author == nil || comment.Author.Username == "") {
		author, err := uc.userRepo.GetUserByID(ctx, comment.AuthorID)
		if err != nil {
			return err
		}
		comment.Author = author
	}

	return nil
}

func (uc *commentUseCase) DeleteComment(ctx context.Context, id int64) error {
	return uc.commentRepo.DeleteComment(ctx, id)
}

func (uc *commentUseCase) LikeComment(ctx context.Context, id int64) error {
	return uc.commentRepo.LikeComment(ctx, id)
}
