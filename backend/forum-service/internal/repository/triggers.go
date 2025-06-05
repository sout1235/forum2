package repository

import (
	"database/sql"
	"log"
)

func CreateTriggers(db *sql.DB) error {
	// Создаем функцию для обновления счетчика комментариев
	_, err := db.Exec(`
		CREATE OR REPLACE FUNCTION update_comment_count()
		RETURNS TRIGGER AS $$
		BEGIN
			IF TG_OP = 'INSERT' THEN
				UPDATE topics SET comment_count = comment_count + 1 WHERE id = NEW.topic_id;
			ELSIF TG_OP = 'DELETE' THEN
				UPDATE topics SET comment_count = comment_count - 1 WHERE id = OLD.topic_id;
			END IF;
			RETURN NULL;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		log.Printf("Error creating update_comment_count function: %v", err)
		return err
	}

	// Создаем триггер для INSERT
	_, err = db.Exec(`
		DROP TRIGGER IF EXISTS update_comment_count_insert ON comments;
		CREATE TRIGGER update_comment_count_insert
		AFTER INSERT ON comments
		FOR EACH ROW
		EXECUTE FUNCTION update_comment_count();
	`)
	if err != nil {
		log.Printf("Error creating insert trigger: %v", err)
		return err
	}

	// Создаем триггер для DELETE
	_, err = db.Exec(`
		DROP TRIGGER IF EXISTS update_comment_count_delete ON comments;
		CREATE TRIGGER update_comment_count_delete
		AFTER DELETE ON comments
		FOR EACH ROW
		EXECUTE FUNCTION update_comment_count();
	`)
	if err != nil {
		log.Printf("Error creating delete trigger: %v", err)
		return err
	}

	log.Println("Triggers created successfully")
	return nil
}
