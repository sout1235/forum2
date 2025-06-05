package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
)

func newTestTopicRepo(t *testing.T) (TopicRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	repo := &topicRepository{db: db}
	return repo, mock, func() { db.Close() }
}

func TestTopicRepository_CreateTopic(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	now := time.Now()
	topic := &entity.Topic{
		Title:      "Test Topic",
		Content:    "Test Content",
		AuthorID:   1,
		CategoryID: 1,
		Views:      0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO topics (title, content, author_id, category_id, views, comment_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`)).
		WithArgs(topic.Title, topic.Content, topic.AuthorID, topic.CategoryID, topic.Views, 0, topic.CreatedAt, topic.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.CreateTopic(context.Background(), topic)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), topic.ID)
}

func TestTopicRepository_GetTopicByID(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	now := time.Now()
	expectedTopic := &entity.Topic{
		ID:         1,
		Title:      "Test Topic",
		Content:    "Test Content",
		AuthorID:   1,
		CategoryID: 1,
		Views:      0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title, content, author_id, category_id, views, created_at, updated_at 
		FROM topics WHERE id = $1
	`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "category_id", "views", "created_at", "updated_at"}).
			AddRow(expectedTopic.ID, expectedTopic.Title, expectedTopic.Content, expectedTopic.AuthorID, expectedTopic.CategoryID, expectedTopic.Views, expectedTopic.CreatedAt, expectedTopic.UpdatedAt))

	topic, err := repo.GetTopicByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedTopic.Title, topic.Title)
	assert.Equal(t, expectedTopic.Content, topic.Content)
}

func TestTopicRepository_GetTopicByID_NotFound(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title, content, author_id, category_id, views, created_at, updated_at 
		FROM topics WHERE id = $1
	`)).
		WithArgs(int64(1)).
		WillReturnError(sql.ErrNoRows)

	topic, err := repo.GetTopicByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, topic)
}

func TestTopicRepository_GetAllTopics(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT t.id, t.title, t.content, t.author_id, t.category_id, t.views, t.comment_count, t.created_at, t.updated_at
		FROM topics t
		ORDER BY t.created_at DESC
	`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "category_id", "views", "comment_count", "created_at", "updated_at"}).
			AddRow(1, "Test Topic 1", "Test Content 1", 1, 1, 0, 0, now, now).
			AddRow(2, "Test Topic 2", "Test Content 2", 1, 1, 0, 0, now, now))

	topics, err := repo.GetAllTopics(context.Background())
	assert.NoError(t, err)
	assert.Len(t, topics, 2)
	assert.Equal(t, int64(1), topics[0].ID)
	assert.Equal(t, "Test Topic 1", topics[0].Title)
	assert.Equal(t, int64(2), topics[1].ID)
	assert.Equal(t, "Test Topic 2", topics[1].Title)
}

func TestTopicRepository_GetAllTopics_ScanError(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT t.id, t.title, t.content, t.author_id, t.category_id, t.views, t.comment_count, t.created_at, t.updated_at
		FROM topics t
		ORDER BY t.created_at DESC
	`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "category_id", "views", "comment_count", "created_at", "updated_at"}).
			AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil))

	topics, err := repo.GetAllTopics(context.Background())
	assert.Error(t, err)
	assert.Nil(t, topics)
}

func TestTopicRepository_UpdateTopic(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	topic := &entity.Topic{
		ID:         1,
		Title:      "Updated Topic",
		Content:    "Updated Content",
		CategoryID: 2,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE topics SET title = $1, content = $2, category_id = $3, updated_at = NOW() 
		WHERE id = $4
	`)).
		WithArgs(topic.Title, topic.Content, topic.CategoryID, topic.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateTopic(context.Background(), topic)
	assert.NoError(t, err)
}

func TestTopicRepository_DeleteTopic(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM topics WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteTopic(context.Background(), 1)
	assert.NoError(t, err)
}

func TestTopicRepository_UpdateCommentCount(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE topics 
		SET comment_count = (
			SELECT COUNT(*) 
			FROM comments 
			WHERE topic_id = $1
		)
		WHERE id = $1
	`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateCommentCount(context.Background(), 1)
	assert.NoError(t, err)
}

func TestTopicRepository_UpdateTopic_ExecError(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	topic := &entity.Topic{
		ID:         1,
		Title:      "Updated Topic",
		Content:    "Updated Content",
		CategoryID: 2,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE topics SET title = $1, content = $2, category_id = $3, updated_at = NOW() 
		WHERE id = $4
	`)).
		WithArgs(topic.Title, topic.Content, topic.CategoryID, topic.ID).
		WillReturnError(assert.AnError)

	err := repo.UpdateTopic(context.Background(), topic)
	assert.Error(t, err)
}

func TestTopicRepository_DeleteTopic_ExecError(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM topics WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnError(assert.AnError)

	err := repo.DeleteTopic(context.Background(), 1)
	assert.Error(t, err)
}

func TestTopicRepository_UpdateCommentCount_ExecError(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE topics 
		SET comment_count = (
			SELECT COUNT(*) 
			FROM comments 
			WHERE topic_id = $1
		)
		WHERE id = $1
	`)).
		WithArgs(int64(1)).
		WillReturnError(assert.AnError)

	err := repo.UpdateCommentCount(context.Background(), 1)
	assert.Error(t, err)
}

func TestTopicRepository_CreateTopic_QueryRowError(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	topic := &entity.Topic{
		Title:      "Test Topic",
		Content:    "Test Content",
		AuthorID:   1,
		CategoryID: 1,
		Views:      0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO topics (title, content, author_id, category_id, views, comment_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`)).
		WithArgs(topic.Title, topic.Content, topic.AuthorID, topic.CategoryID, topic.Views, 0, topic.CreatedAt, topic.UpdatedAt).
		WillReturnError(assert.AnError)

	err := repo.CreateTopic(context.Background(), topic)
	assert.Error(t, err)
}

func TestTopicRepository_GetTopicByID_QueryRowError(t *testing.T) {
	repo, mock, closeFn := newTestTopicRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title, content, author_id, category_id, views, created_at, updated_at 
		FROM topics WHERE id = $1
	`)).
		WithArgs(int64(1)).
		WillReturnError(assert.AnError)

	topic, err := repo.GetTopicByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, topic)
}
