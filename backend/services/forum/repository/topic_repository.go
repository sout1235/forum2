package repository

import (
	"backend/services/forum/models"
)

type TopicRepository struct {
	DB *DB
}

func NewTopicRepository(db *DB) *TopicRepository {
	return &TopicRepository{DB: db}
}

func (r *TopicRepository) GetAllTopics() ([]models.Topic, error) {
	rows, err := r.DB.Conn.Query(`SELECT id, title, content, author_id, category_id, views, created_at, updated_at FROM topics ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topics := []models.Topic{}
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.AuthorID, &t.CategoryID, &t.Views, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func (r *TopicRepository) GetTopicByID(id int64) (*models.Topic, error) {
	row := r.DB.Conn.QueryRow(`SELECT id, title, content, author_id, category_id, views, created_at, updated_at FROM topics WHERE id = $1`, id)
	var t models.Topic
	if err := row.Scan(&t.ID, &t.Title, &t.Content, &t.AuthorID, &t.CategoryID, &t.Views, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TopicRepository) CreateTopic(t *models.Topic) error {
	return r.DB.Conn.QueryRow(
		`INSERT INTO topics (title, content, author_id, category_id, views) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`,
		t.Title, t.Content, t.AuthorID, t.CategoryID, t.Views,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TopicRepository) UpdateTopic(t *models.Topic) error {
	_, err := r.DB.Conn.Exec(
		`UPDATE topics SET title=$1, content=$2, category_id=$3, updated_at=NOW() WHERE id=$4`,
		t.Title, t.Content, t.CategoryID, t.ID,
	)
	return err
}

func (r *TopicRepository) DeleteTopic(id int64) error {
	_, err := r.DB.Conn.Exec(`DELETE FROM topics WHERE id = $1`, id)
	return err
}
