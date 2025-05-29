package config

import (
	"os"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	AuthURL     string
	AuthGRPCURL string
}

func NewConfig() *Config {
	return &Config{
		DBHost:      getEnv("FORUM_DB_HOST", "localhost"),
		DBPort:      getEnv("FORUM_DB_PORT", "5432"),
		DBUser:      getEnv("FORUM_DB_USER", "postgres"),
		DBPassword:  getEnv("FORUM_DB_PASSWORD", "postgres"),
		DBName:      getEnv("FORUM_DB_NAME", "forum"),
		AuthURL:     getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		AuthGRPCURL: getEnv("AUTH_GRPC_URL", "localhost:50051"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
