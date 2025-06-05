package service

import (
	"context"
	"time"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
)

type chatService struct {
	chatRepo repository.ChatRepository
}

// NewChatService creates a new instance of ChatService
func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (s *chatService) SaveMessage(ctx context.Context, message *entity.ChatMessage) error {
	if message.ExpiresAt.IsZero() {
		message.ExpiresAt = time.Now().Add(24 * time.Hour)
	}
	return s.chatRepo.SaveMessage(ctx, message)
}

func (s *chatService) GetRecentMessages(ctx context.Context, limit int) ([]*entity.ChatMessage, error) {
	return s.chatRepo.GetRecentMessages(ctx, limit)
}

func (s *chatService) DeleteExpiredMessages(ctx context.Context) error {
	return s.chatRepo.DeleteExpiredMessages(ctx)
}
