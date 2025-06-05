package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	// Получаем параметры подключения к БД из переменных окружения
	dbHost := getEnv("AUTH_DB_HOST", "localhost")
	dbPort := getEnv("AUTH_DB_PORT", "5432")
	dbUser := getEnv("AUTH_DB_USER", "postgres")
	dbPassword := getEnv("AUTH_DB_PASSWORD", "postgres")
	dbName := getEnv("AUTH_DB_NAME", "auth")

	// Формируем строку подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Подключаемся к БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Читаем файлы миграций
	migrationsDir := "../migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Error reading migrations directory: %v", err)
	}

	// Применяем миграции
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				log.Fatalf("Error reading migration file %s: %v", file.Name(), err)
			}

			_, err = db.Exec(string(content))
			if err != nil {
				log.Fatalf("Error executing migration %s: %v", file.Name(), err)
			}

			log.Printf("Successfully applied migration: %s", file.Name())
		}
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
