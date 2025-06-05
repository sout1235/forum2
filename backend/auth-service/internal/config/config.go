package config

import (
	"os"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	HTTPPort      string
	GRPCPort      string
	JWTSecret     string
	InternalToken string
}

func NewConfig() *Config {
	return &Config{
		DBHost:        getEnv("AUTH_DB_HOST", "localhost"),
		DBPort:        getEnv("AUTH_DB_PORT", "5432"),
		DBUser:        getEnv("AUTH_DB_USER", "postgres"),
		DBPassword:    getEnv("AUTH_DB_PASSWORD", "postgres"),
		DBName:        getEnv("AUTH_DB_NAME", "auth"),
		HTTPPort:      getEnv("AUTH_HTTP_PORT", "8080"),
		GRPCPort:      getEnv("AUTH_GRPC_PORT", "50051"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		InternalToken: getEnv("INTERNAL_TOKEN", "internal-service-token"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
