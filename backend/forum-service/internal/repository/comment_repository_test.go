package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
)

func newTestCommentRepo(t *testing.T) (CommentRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	repo := &commentRepository{db: db}
	return repo, mock, func() { db.Close() }
}

func TestCommentRepository_GetCommentsByTopic(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	// Создаем тестовые данные
	now := time.Now()
	expectedComments := []*entity.Comment{
		{
			ID:        1,
			Content:   "First comment",
			AuthorID:  1,
			TopicID:   1,
			ParentID:  nil,
			Likes:     0,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Content:   "Second comment",
			AuthorID:  2,
			TopicID:   1,
			ParentID:  nil,
			Likes:     0,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Ожидаем, что будет выполнен запрос на получение комментариев
	rows := sqlmock.NewRows([]string{"id", "content", "author_id", "topic_id", "parent_id", "likes", "created_at", "updated_at"})
	for _, comment := range expectedComments {
		rows.AddRow(comment.ID, comment.Content, comment.AuthorID, comment.TopicID, comment.ParentID, comment.Likes, comment.CreatedAt, comment.UpdatedAt)
	}

	mock.ExpectQuery(`SELECT id, content, author_id, topic_id, parent_id, likes, created_at, updated_at FROM comments WHERE topic_id = \$1 ORDER BY created_at ASC`).
		WithArgs(1).
		WillReturnRows(rows)

	// Вызываем тестируемый метод
	comments, err := repo.GetCommentsByTopic(context.Background(), 1)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, expectedComments[0].Content, comments[0].Content)
	assert.Equal(t, expectedComments[1].Content, comments[1].Content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepository_GetCommentByID(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	// Создаем тестовые данные
	now := time.Now()
	expectedComment := &entity.Comment{
		ID:        1,
		Content:   "Test comment",
		AuthorID:  1,
		TopicID:   1,
		ParentID:  nil,
		Likes:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Ожидаем, что будет выполнен запрос на получение комментария
	rows := sqlmock.NewRows([]string{"id", "content", "author_id", "topic_id", "parent_id", "likes", "created_at", "updated_at"}).
		AddRow(expectedComment.ID, expectedComment.Content, expectedComment.AuthorID, expectedComment.TopicID, expectedComment.ParentID, expectedComment.Likes, expectedComment.CreatedAt, expectedComment.UpdatedAt)

	mock.ExpectQuery(`SELECT id, content, author_id, topic_id, parent_id, likes, created_at, updated_at FROM comments WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	// Вызываем тестируемый метод
	comment, err := repo.GetCommentByID(context.Background(), 1)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, expectedComment.Content, comment.Content)
	assert.Equal(t, expectedComment.AuthorID, comment.AuthorID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepository_CreateComment(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	// Создаем тестовый комментарий
	now := time.Now()
	comment := &entity.Comment{
		Content:   "Test comment",
		AuthorID:  1,
		TopicID:   1,
		ParentID:  nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Ожидаем, что будет выполнен запрос на вставку
	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(comment.Content, comment.AuthorID, comment.TopicID, comment.ParentID, comment.CreatedAt, comment.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Вызываем тестируемый метод
	err := repo.CreateComment(context.Background(), comment)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, int64(1), comment.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepository_DeleteComment(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	// Ожидаем, что будет выполнен запрос на удаление комментария
	mock.ExpectExec(`DELETE FROM comments WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Вызываем тестируемый метод
	err := repo.DeleteComment(context.Background(), 1)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepository_LikeComment(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	// Ожидаем, что будет выполнен запрос на увеличение количества лайков
	mock.ExpectExec(`UPDATE comments SET likes = likes \+ 1 WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Вызываем тестируемый метод
	err := repo.LikeComment(context.Background(), 1)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepository_GetCommentByID_NotFound(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, content, author_id, topic_id, parent_id, likes, created_at, updated_at
		FROM comments
		WHERE id = $1
	`)).
		WithArgs(int64(1)).
		WillReturnError(sql.ErrNoRows)

	comment, err := repo.GetCommentByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, comment)
}

func TestCommentRepository_GetCommentByID_ScanError(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, content, author_id, topic_id, parent_id, likes, created_at, updated_at
		FROM comments
		WHERE id = $1
	`)).
		WithArgs(int64(1)).
		WillReturnError(errors.New("scan error"))

	comment, err := repo.GetCommentByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, comment)
}

func TestCommentRepository_CreateComment_QueryRowError(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	comment := &entity.Comment{
		Content:  "test comment",
		AuthorID: 1,
		TopicID:  1,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO comments (content, author_id, topic_id, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`)).
		WithArgs(comment.Content, comment.AuthorID, comment.TopicID, comment.ParentID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(assert.AnError)

	err := repo.CreateComment(context.Background(), comment)
	assert.Error(t, err)
}

func TestCommentRepository_LikeComment_ExecError(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE comments
		SET likes = likes + 1
		WHERE id = $1
	`)).
		WithArgs(int64(1)).
		WillReturnError(assert.AnError)

	err := repo.LikeComment(context.Background(), 1)
	assert.Error(t, err)
}

func TestCommentRepository_GetCommentsByTopic_ScanError(t *testing.T) {
	repo, mock, closeFn := newTestCommentRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, content, author_id, topic_id, parent_id, likes, created_at, updated_at
		FROM comments
		WHERE topic_id = $1
		ORDER BY created_at ASC
	`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "author_id", "topic_id", "parent_id", "likes", "created_at", "updated_at"}).
			AddRow(nil, nil, nil, nil, nil, nil, nil, nil))

	comments, err := repo.GetCommentsByTopic(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, comments)
}
