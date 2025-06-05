package service

import (
	"context"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment *entity.Comment) error
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
	GetCommentsByTopicID(ctx context.Context, topicID int64) ([]*entity.Comment, error)
	UpdateComment(ctx context.Context, comment *entity.Comment) error
	DeleteComment(ctx context.Context, id int64) error
	LikeComment(ctx context.Context, id int64) error
}

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetUsernameByID(ctx context.Context, id int64) (string, error)
}

type ChatService interface {
	SaveMessage(ctx context.Context, message *entity.ChatMessage) error
	GetRecentMessages(ctx context.Context, limit int) ([]*entity.ChatMessage, error)
	DeleteExpiredMessages(ctx context.Context) error
}
