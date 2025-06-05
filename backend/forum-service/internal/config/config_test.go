package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Сбрасываем переменные окружения перед тестом
	os.Unsetenv("FORUM_DB_HOST")
	os.Unsetenv("FORUM_DB_PORT")
	os.Unsetenv("FORUM_DB_USER")
	os.Unsetenv("FORUM_DB_PASSWORD")
	os.Unsetenv("FORUM_DB_NAME")
	os.Unsetenv("FORUM_HTTP_PORT")
	os.Unsetenv("FORUM_GRPC_PORT")
	os.Unsetenv("AUTH_SERVICE_URL")

	cfg := NewConfig()
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "postgres", cfg.DBPassword)
	assert.Equal(t, "forum", cfg.DBName)
	assert.Equal(t, "8081", cfg.HTTPPort)
	assert.Equal(t, "50052", cfg.GRPCPort)
	assert.Equal(t, "http://localhost:8080", cfg.AuthServiceURL)

	// Устанавливаем переменные окружения и проверяем, что они используются
	os.Setenv("FORUM_DB_HOST", "custom_host")
	os.Setenv("FORUM_DB_PORT", "1234")
	cfg = NewConfig()
	assert.Equal(t, "custom_host", cfg.DBHost)
	assert.Equal(t, "1234", cfg.DBPort)
}

func TestGetEnv(t *testing.T) {
	// Сбрасываем переменную окружения перед тестом
	os.Unsetenv("TEST_KEY")

	// Проверяем, что возвращается дефолтное значение, если переменная не задана
	assert.Equal(t, "default", getEnv("TEST_KEY", "default"))

	// Устанавливаем переменную окружения и проверяем, что её значение используется
	os.Setenv("TEST_KEY", "custom_value")
	assert.Equal(t, "custom_value", getEnv("TEST_KEY", "default"))
}
