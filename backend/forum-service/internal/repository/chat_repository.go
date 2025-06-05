package repository

import (
	"context"
	"database/sql"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, message *entity.ChatMessage) error
	GetRecentMessages(ctx context.Context, limit int) ([]*entity.ChatMessage, error)
	DeleteExpiredMessages(ctx context.Context) error
}

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) SaveMessage(ctx context.Context, message *entity.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (content, author_id, author_username, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return r.db.QueryRowContext(
		ctx,
		query,
		message.Content,
		message.AuthorID,
		message.AuthorUsername,
		message.CreatedAt,
		message.ExpiresAt,
	).Scan(&message.ID)
}

func (r *chatRepository) GetRecentMessages(ctx context.Context, limit int) ([]*entity.ChatMessage, error) {
	query := `
		SELECT id, content, author_id, author_username, created_at, expires_at
		FROM chat_messages
		WHERE expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entity.ChatMessage
	for rows.Next() {
		msg := &entity.ChatMessage{}
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.AuthorID,
			&msg.AuthorUsername,
			&msg.CreatedAt,
			&msg.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *chatRepository) DeleteExpiredMessages(ctx context.Context) error {
	query := `DELETE FROM chat_messages WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
