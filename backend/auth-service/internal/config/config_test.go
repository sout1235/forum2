package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Сохраняем текущие значения переменных окружения
	envVars := map[string]string{
		"AUTH_DB_HOST":     os.Getenv("AUTH_DB_HOST"),
		"AUTH_DB_PORT":     os.Getenv("AUTH_DB_PORT"),
		"AUTH_DB_USER":     os.Getenv("AUTH_DB_USER"),
		"AUTH_DB_PASSWORD": os.Getenv("AUTH_DB_PASSWORD"),
		"AUTH_DB_NAME":     os.Getenv("AUTH_DB_NAME"),
		"AUTH_HTTP_PORT":   os.Getenv("AUTH_HTTP_PORT"),
		"AUTH_GRPC_PORT":   os.Getenv("AUTH_GRPC_PORT"),
		"JWT_SECRET":       os.Getenv("JWT_SECRET"),
		"INTERNAL_TOKEN":   os.Getenv("INTERNAL_TOKEN"),
	}

	// Устанавливаем тестовые значения
	os.Setenv("AUTH_DB_HOST", "test-host")
	os.Setenv("AUTH_DB_PORT", "5433")
	os.Setenv("AUTH_DB_USER", "test-user")
	os.Setenv("AUTH_DB_PASSWORD", "test-password")
	os.Setenv("AUTH_DB_NAME", "test-db")
	os.Setenv("AUTH_HTTP_PORT", "8081")
	os.Setenv("AUTH_GRPC_PORT", "50052")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("INTERNAL_TOKEN", "test-token")

	// Восстанавливаем переменные окружения после теста
	defer func() {
		for k, v := range envVars {
			os.Setenv(k, v)
		}
	}()

	config := NewConfig()

	assert.Equal(t, "test-host", config.DBHost)
	assert.Equal(t, "5433", config.DBPort)
	assert.Equal(t, "test-user", config.DBUser)
	assert.Equal(t, "test-password", config.DBPassword)
	assert.Equal(t, "test-db", config.DBName)
	assert.Equal(t, "8081", config.HTTPPort)
	assert.Equal(t, "50052", config.GRPCPort)
	assert.Equal(t, "test-secret", config.JWTSecret)
	assert.Equal(t, "test-token", config.InternalToken)
}
