package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetUsernameByID_LocalDB(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Создаем тестовый сервер для auth-service
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"user_id":  "1",
			"username": "testuser",
		})
	}))
	defer server.Close()

	// Создаем репозиторий с моком
	repo := &UserRepositoryImpl{
		db:      db,
		authURL: server.URL,
	}

	// Ожидаем, что будет выполнен запрос к локальной базе
	mock.ExpectQuery("SELECT username FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("localuser"))

	// Вызываем тестируемый метод
	username, err := repo.GetUsernameByID(context.Background(), 1)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, "localuser", username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsernameByID_AuthService(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Создаем тестовый сервер для auth-service
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"user_id":  "1",
			"username": "authuser",
		})
	}))
	defer server.Close()

	// Создаем репозиторий с моком
	repo := &UserRepositoryImpl{
		db:      db,
		authURL: server.URL,
	}

	// Ожидаем, что запрос к локальной базе вернет ErrNoRows
	mock.ExpectQuery("SELECT username FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	// Ожидаем, что будет выполнен запрос на сохранение пользователя в локальной базе
	mock.ExpectExec("INSERT INTO users").
		WithArgs(1, "authuser").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызываем тестируемый метод
	username, err := repo.GetUsernameByID(context.Background(), 1)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, "authuser", username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsernameByID_AuthServiceError(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Создаем тестовый сервер для auth-service, который возвращает ошибку
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Создаем репозиторий с моком
	repo := &UserRepositoryImpl{
		db:      db,
		authURL: server.URL,
	}

	// Ожидаем, что запрос к локальной базе вернет ErrNoRows
	mock.ExpectQuery("SELECT username FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	// Вызываем тестируемый метод
	username, err := repo.GetUsernameByID(context.Background(), 1)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Empty(t, username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsernameByID_AuthServiceNotFound(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Создаем тестовый сервер для auth-service, который возвращает 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Создаем репозиторий с моком
	repo := &UserRepositoryImpl{
		db:      db,
		authURL: server.URL,
	}

	// Ожидаем, что запрос к локальной базе вернет ErrNoRows
	mock.ExpectQuery("SELECT username FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	// Вызываем тестируемый метод
	username, err := repo.GetUsernameByID(context.Background(), 1)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, "User_1", username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsernameByID_FromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	repo := NewUserRepository(db, "http://auth-service")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT username FROM users WHERE id = $1")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("testuser"))

	username, err := repo.GetUsernameByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsernameByID_FromAuthService(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// Не найдено в локальной базе
	mock.ExpectQuery(regexp.QuoteMeta("SELECT username FROM users WHERE id = $1")).
		WithArgs(int64(2)).
		WillReturnError(sql.ErrNoRows)
	// Ожидается вставка в базу
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id, username) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET username = $2")).
		WithArgs(int64(2), "remoteuser").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Мокаем внешний сервис
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"user_id":"2","username":"remoteuser"}`))
	}))
	defer ts.Close()

	repo := NewUserRepository(db, ts.URL)
	username, err := repo.GetUsernameByID(context.Background(), 2)
	assert.NoError(t, err)
	assert.Equal(t, "remoteuser", username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsernameByID_NotFoundAnywhere(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT username FROM users WHERE id = $1")).
		WithArgs(int64(3)).
		WillReturnError(sql.ErrNoRows)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer ts.Close()

	repo := NewUserRepository(db, ts.URL)
	username, err := repo.GetUsernameByID(context.Background(), 3)
	assert.NoError(t, err)
	assert.Equal(t, "User_3", username)
	assert.NoError(t, mock.ExpectationsWereMet())
}
