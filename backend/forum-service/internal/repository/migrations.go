package repository

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	// Создаем таблицу topics, если она не существует
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS topics (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author_id BIGINT NOT NULL,
			category_id BIGINT,
			views INTEGER DEFAULT 0,
			comment_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("Error creating topics table: %v", err)
		return err
	}

	// Добавляем колонку comment_count, если она не существует
	_, err = db.Exec(`
		DO $$ 
		BEGIN 
			IF NOT EXISTS (
				SELECT 1 
				FROM information_schema.columns 
				WHERE table_name = 'topics' 
				AND column_name = 'comment_count'
			) THEN
				ALTER TABLE topics ADD COLUMN comment_count INTEGER DEFAULT 0;
			END IF;
		END $$;
	`)
	if err != nil {
		log.Printf("Error adding comment_count column: %v", err)
		return err
	}

	// Создаем таблицу comments, если она не существует
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id BIGSERIAL PRIMARY KEY,
			content TEXT NOT NULL,
			author_id BIGINT NOT NULL,
			topic_id BIGINT NOT NULL,
			parent_id BIGINT,
			likes INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE,
			FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Printf("Error creating comments table: %v", err)
		return err
	}

	// Создаем таблицу chat_messages, если она не существует
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			id SERIAL PRIMARY KEY,
			content TEXT NOT NULL,
			author_id BIGINT,
			author_username VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`)
	if err != nil {
		log.Printf("Error creating chat_messages table: %v", err)
		return err
	}

	// Создаем индексы для chat_messages
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at);
		CREATE INDEX IF NOT EXISTS idx_chat_messages_expires_at ON chat_messages(expires_at);
		CREATE INDEX IF NOT EXISTS idx_chat_messages_author_id ON chat_messages(author_id);
	`)
	if err != nil {
		log.Printf("Error creating chat_messages indexes: %v", err)
		return err
	}

	// Создаем функцию для автоматического удаления устаревших сообщений
	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION delete_expired_messages()
		RETURNS void AS $$
		BEGIN
			DELETE FROM chat_messages WHERE expires_at < CURRENT_TIMESTAMP;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		log.Printf("Error creating delete_expired_messages function: %v", err)
		return err
	}

	// Создаем триггер для автоматического удаления устаревших сообщений
	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION trigger_delete_expired_messages()
		RETURNS trigger AS $$
		BEGIN
			PERFORM delete_expired_messages();
			RETURN NULL;
		END;
		$$ LANGUAGE plpgsql;

		DROP TRIGGER IF EXISTS chat_messages_cleanup_trigger ON chat_messages;
		CREATE TRIGGER chat_messages_cleanup_trigger
		AFTER INSERT ON chat_messages
		EXECUTE FUNCTION trigger_delete_expired_messages();
	`)
	if err != nil {
		log.Printf("Error creating chat_messages cleanup trigger: %v", err)
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
