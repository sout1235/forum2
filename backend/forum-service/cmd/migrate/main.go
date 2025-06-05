package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/sout1235/forum2/backend/forum-service/internal/config"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.NewConfig()

	// Подключаемся к базе данных
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Читаем файлы миграций
	migrationsDir := "../../migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Применяем миграции
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				log.Fatalf("Failed to read migration file %s: %v", file.Name(), err)
			}

			// Выполняем миграцию
			_, err = db.Exec(string(content))
			if err != nil {
				log.Fatalf("Failed to execute migration %s: %v", file.Name(), err)
			}

			log.Printf("Successfully applied migration: %s", file.Name())
		}
	}

	log.Println("All migrations completed successfully")
}
