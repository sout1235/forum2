package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestServer(handler gin.HandlerFunc) *httptest.Server {
	r := gin.Default()
	r.GET("/protected", handler)
	return httptest.NewServer(r)
}

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	cfg := &AuthConfig{AuthServiceURL: "http://auth"}
	mw := cfg.AuthMiddleware()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)

	mw(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "authorization header is required")
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	cfg := &AuthConfig{AuthServiceURL: "http://auth"}
	mw := cfg.AuthMiddleware()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	mw(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header format")
}

func TestAuthMiddleware_EmptyToken(t *testing.T) {
	cfg := &AuthConfig{AuthServiceURL: "http://auth"}
	mw := cfg.AuthMiddleware()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer ")

	mw(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "token is required")
}

func TestAuthMiddleware_AuthServiceError(t *testing.T) {
	// Создаем тестовый сервер, который всегда возвращает ошибку
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer authServer.Close()

	cfg := &AuthConfig{AuthServiceURL: authServer.URL}
	mw := cfg.AuthMiddleware()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer validtoken")

	mw(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Создаем тестовый сервер, который возвращает 401
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer authServer.Close()

	cfg := &AuthConfig{AuthServiceURL: authServer.URL}
	mw := cfg.AuthMiddleware()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer invalidtoken")

	mw(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAuthMiddleware_DecodeError(t *testing.T) {
	// Создаем тестовый сервер, который возвращает 200, но невалидный JSON
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("notjson"))
	}))
	defer authServer.Close()

	cfg := &AuthConfig{AuthServiceURL: authServer.URL}
	mw := cfg.AuthMiddleware()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer validtoken")

	mw(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to decode user data")
}

func TestAuthMiddleware_Success(t *testing.T) {
	// Создаем тестовый сервер, который возвращает 200 и валидный JSON
	userData := map[string]string{"user_id": "123", "username": "testuser"}
	respBody, _ := json.Marshal(userData)
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}))
	defer authServer.Close()

	cfg := &AuthConfig{AuthServiceURL: authServer.URL}
	mw := cfg.AuthMiddleware()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer validtoken")

	mw(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "123", c.GetString("user_id"))
	assert.Equal(t, "testuser", c.GetString("username"))
}
