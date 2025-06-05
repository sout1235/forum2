package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
)

func newTestChatRepo(t *testing.T) (ChatRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	repo := NewChatRepository(db)
	return repo, mock, func() { db.Close() }
}

func TestChatRepository_SaveMessage(t *testing.T) {
	repo, mock, closeFn := newTestChatRepo(t)
	defer closeFn()

	// Создаем тестовое сообщение
	now := time.Now()
	message := &entity.ChatMessage{
		Content:        "Test message",
		AuthorID:       1,
		AuthorUsername: "user1",
		CreatedAt:      now,
		ExpiresAt:      now.Add(24 * time.Hour),
	}

	// Ожидаем, что будет выполнен запрос на сохранение сообщения
	mock.ExpectQuery(`INSERT INTO chat_messages`).
		WithArgs(message.Content, message.AuthorID, message.AuthorUsername, message.CreatedAt, message.ExpiresAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Вызываем тестируемый метод
	err := repo.SaveMessage(context.Background(), message)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, int64(1), message.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetRecentMessages(t *testing.T) {
	repo, mock, closeFn := newTestChatRepo(t)
	defer closeFn()

	// Создаем тестовые сообщения
	now := time.Now()
	expectedMessages := []*entity.ChatMessage{
		{
			ID:             1,
			Content:        "Test message 1",
			AuthorID:       1,
			AuthorUsername: "user1",
			CreatedAt:      now,
			ExpiresAt:      now.Add(24 * time.Hour),
		},
		{
			ID:             2,
			Content:        "Test message 2",
			AuthorID:       2,
			AuthorUsername: "user2",
			CreatedAt:      now,
			ExpiresAt:      now.Add(24 * time.Hour),
		},
	}

	// Ожидаем, что будет выполнен запрос на получение сообщений
	rows := sqlmock.NewRows([]string{"id", "content", "author_id", "author_username", "created_at", "expires_at"})
	for _, msg := range expectedMessages {
		rows.AddRow(msg.ID, msg.Content, msg.AuthorID, msg.AuthorUsername, msg.CreatedAt, msg.ExpiresAt)
	}

	mock.ExpectQuery(`SELECT id, content, author_id, author_username, created_at, expires_at FROM chat_messages WHERE expires_at > CURRENT_TIMESTAMP ORDER BY created_at DESC LIMIT \$1`).
		WithArgs(10).
		WillReturnRows(rows)

	// Вызываем тестируемый метод
	messages, err := repo.GetRecentMessages(context.Background(), 10)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, expectedMessages[0].Content, messages[0].Content)
	assert.Equal(t, expectedMessages[1].Content, messages[1].Content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetRecentMessages_ScanError(t *testing.T) {
	repo, mock, closeFn := newTestChatRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, content, author_id, author_username, created_at, expires_at 
		FROM chat_messages 
		WHERE expires_at > CURRENT_TIMESTAMP 
		ORDER BY created_at DESC 
		LIMIT $1`)).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "author_id", "author_username", "created_at", "expires_at"}).
			AddRow(nil, nil, nil, nil, nil, nil))

	messages, err := repo.GetRecentMessages(context.Background(), 10)
	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_DeleteExpiredMessages(t *testing.T) {
	repo, mock, closeFn := newTestChatRepo(t)
	defer closeFn()

	// Ожидаем, что будет выполнен запрос на удаление устаревших сообщений
	mock.ExpectExec(`DELETE FROM chat_messages WHERE expires_at < CURRENT_TIMESTAMP`).
		WillReturnResult(sqlmock.NewResult(0, 2))

	// Вызываем тестируемый метод
	err := repo.DeleteExpiredMessages(context.Background())

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_DeleteExpiredMessages_ExecError(t *testing.T) {
	repo, mock, closeFn := newTestChatRepo(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM chat_messages WHERE expires_at < CURRENT_TIMESTAMP`)).
		WillReturnError(assert.AnError)

	err := repo.DeleteExpiredMessages(context.Background())
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
