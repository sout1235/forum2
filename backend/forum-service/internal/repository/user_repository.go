package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
)

type UserRepository interface {
	GetUsernameByID(ctx context.Context, id int64) (string, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
}

type UserRepositoryImpl struct {
	db      *sql.DB
	authURL string
}

func NewUserRepository(db *sql.DB, authURL string) UserRepository {
	return &UserRepositoryImpl{
		db:      db,
		authURL: authURL,
	}
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	log.Printf("Getting user data for ID: %d", id)
	log.Printf("Using auth service URL: %s", r.authURL)

	// Try to get from local database first
	user := &entity.User{ID: id}
	err := r.db.QueryRowContext(ctx, "SELECT username, avatar FROM users WHERE id = $1", id).Scan(&user.Username, &user.Avatar)
	if err == nil {
		log.Printf("Found user in local DB: %+v", user)
		return user, nil
	} else if err != sql.ErrNoRows {
		log.Printf("Error querying local DB: %v", err)
		return nil, err
	}
	log.Printf("User not found in local DB, trying auth-service")

	// If not found in local DB, try to get from auth-service
	url := fmt.Sprintf("%s/user/%d", r.authURL, id)
	log.Printf("Fetching user data from auth-service: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer internal-service-token")
	log.Printf("Request headers: %v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to auth-service: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Auth service response status: %d", resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Auth service response body: %s", string(body))

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("User not found in auth-service for ID: %d", id)
		return &entity.User{
			ID:       id,
			Username: fmt.Sprintf("User_%d", id),
			Avatar:   "",
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Auth-service returned non-OK status: %d", resp.StatusCode)
		return nil, errors.New("failed to get user data from auth-service")
	}

	var userData struct {
		ID       string `json:"user_id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&userData); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	// Save user to local database
	_, err = r.db.ExecContext(ctx, "INSERT INTO users (id, username, avatar) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET username = $2, avatar = $3",
		id, userData.Username, userData.Avatar)
	if err != nil {
		log.Printf("Error saving user to local DB: %v", err)
	}

	user = &entity.User{
		ID:       id,
		Username: userData.Username,
		Avatar:   userData.Avatar,
	}
	log.Printf("Found user from auth-service: %+v", user)
	return user, nil
}

func (r *UserRepositoryImpl) GetUsernameByID(ctx context.Context, id int64) (string, error) {
	log.Printf("Getting username for user ID: %d", id)
	log.Printf("Using auth service URL: %s", r.authURL)

	// Сначала пробуем получить из локальной базы
	var username string
	err := r.db.QueryRowContext(ctx, "SELECT username FROM users WHERE id = $1", id).Scan(&username)
	if err == nil {
		log.Printf("Found username in local DB: %s for user ID: %d", username, id)
		return username, nil
	} else if err != sql.ErrNoRows {
		log.Printf("Error querying local DB: %v", err)
		return "", err
	}
	log.Printf("User not found in local DB, trying auth-service")

	// Если не нашли в локальной базе, пробуем получить из auth-service
	url := fmt.Sprintf("%s/user/%d", r.authURL, id)
	log.Printf("Fetching user data from auth-service: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer internal-service-token")
	log.Printf("Request headers: %v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to auth-service: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("Auth service response status: %d", resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Auth service response body: %s", string(body))

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("User not found in auth-service for ID: %d", id)
		return fmt.Sprintf("User_%d", id), nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Auth-service returned non-OK status: %d", resp.StatusCode)
		return "", errors.New("failed to get user data from auth-service")
	}

	var user struct {
		ID       string `json:"user_id"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&user); err != nil {
		log.Printf("Error decoding response: %v", err)
		return "", err
	}

	// Сохраняем пользователя в локальной базе
	_, err = r.db.ExecContext(ctx, "INSERT INTO users (id, username) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET username = $2",
		id, user.Username)
	if err != nil {
		log.Printf("Error saving user to local DB: %v", err)
	}

	log.Printf("Found username from auth-service: %s for user ID: %d", user.Username, id)
	return user.Username, nil
}
