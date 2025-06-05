package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sout1235/forum2/backend/auth-service/internal/entity"
	"github.com/stretchr/testify/assert"
)

func newTestRepo(t *testing.T) (*UserRepositoryImpl, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	repo := &UserRepositoryImpl{db: db}
	return repo, mock, func() { db.Close() }
}

func TestUserRepositoryImpl_ExistsByUsername(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`)).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsByUsername(context.Background(), "testuser")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	user := &entity.User{
		Username:     "testuser",
		PasswordHash: "hash",
		Email:        "test@example.com",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (username, password_hash, email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`)).
		WithArgs(user.Username, user.PasswordHash, user.Email, user.Role, user.CreatedAt, user.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
}

func TestUserRepositoryImpl_GetByID(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, password_hash, email, role, created_at, updated_at 
		FROM users WHERE id = $1`)).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "email", "role", "created_at", "updated_at"}).
			AddRow("1", "testuser", "hash", "test@example.com", "user", time.Now(), time.Now()))

	user, err := repo.GetByID(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "testuser", user.Username)
}

func TestUserRepositoryImpl_Update(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	user := &entity.User{ID: "1", Username: "updateduser", PasswordHash: "newhash"}
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET username = $1, password_hash = $2 WHERE id = $3`)).
		WithArgs(user.Username, user.PasswordHash, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), user)
	assert.NoError(t, err)
}

func TestUserRepositoryImpl_DeleteUser(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(context.Background(), "1")
	assert.NoError(t, err)
}

func TestUserRepositoryImpl_GetByEmail(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, password_hash, email, role, created_at, updated_at FROM users WHERE email = $1`)).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "email", "role", "created_at", "updated_at"}).
			AddRow("1", "testuser", "hash", "test@example.com", "user", time.Now(), time.Now()))

	user, err := repo.GetByEmail(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUserRepositoryImpl_GetByUsername(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, password_hash FROM users WHERE username = $1`)).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash"}).
			AddRow("1", "testuser", "hash"))

	user, err := repo.GetByUsername(context.Background(), "testuser")
	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "testuser", user.Username)
}

func TestUserRepositoryImpl_GetByUsername_NotFound(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, password_hash FROM users WHERE username = $1`)).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByUsername(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepositoryImpl_ExistsByEmail(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`)).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsByEmail(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepositoryImpl_ExistsByEmail_NotFound(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`)).
		WithArgs("nonexistent@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := repo.ExistsByEmail(context.Background(), "nonexistent@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepositoryImpl_Create_Error(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	user := &entity.User{
		Username:     "testuser",
		PasswordHash: "hash",
		Email:        "test@example.com",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (username, password_hash, email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`)).
		WithArgs(user.Username, user.PasswordHash, user.Email, user.Role, user.CreatedAt, user.UpdatedAt).
		WillReturnError(sql.ErrConnDone)

	err := repo.Create(context.Background(), user)
	assert.Error(t, err)
}

func TestUserRepositoryImpl_GetByID_Error(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, password_hash, email, role, created_at, updated_at 
		FROM users WHERE id = $1`)).
		WithArgs("1").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(context.Background(), "1")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepositoryImpl_Update_Error(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	user := &entity.User{ID: "1", Username: "updateduser", PasswordHash: "newhash"}
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET username = $1, password_hash = $2 WHERE id = $3`)).
		WithArgs(user.Username, user.PasswordHash, user.ID).
		WillReturnError(sql.ErrConnDone)

	err := repo.Update(context.Background(), user)
	assert.Error(t, err)
}

func TestUserRepositoryImpl_DeleteUser_Error(t *testing.T) {
	repo, mock, closeFn := newTestRepo(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs("1").
		WillReturnError(sql.ErrConnDone)

	err := repo.DeleteUser(context.Background(), "1")
	assert.Error(t, err)
}
