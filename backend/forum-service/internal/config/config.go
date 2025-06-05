package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DatabaseURL    string
	HTTPPort       string
	GRPCPort       string
	AuthURL        string
	AuthGRPCURL    string
	AuthServiceURL string
}

func NewConfig() *Config {
	cfg := &Config{
		DBHost:         getEnv("FORUM_DB_HOST", "localhost"),
		DBPort:         getEnv("FORUM_DB_PORT", "5432"),
		DBUser:         getEnv("FORUM_DB_USER", "postgres"),
		DBPassword:     getEnv("FORUM_DB_PASSWORD", "postgres"),
		DBName:         getEnv("FORUM_DB_NAME", "forum"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		HTTPPort:       getEnv("FORUM_HTTP_PORT", "8081"),
		GRPCPort:       getEnv("FORUM_GRPC_PORT", "50052"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
	}

	// Если DATABASE_URL не указан, формируем его из отдельных параметров
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
