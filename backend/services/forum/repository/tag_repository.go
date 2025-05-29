package repository

import (
	"backend/services/forum/models"
)

type TagRepository struct {
	DB *DB
}

func NewTagRepository(db *DB) *TagRepository {
	return &TagRepository{DB: db}
}

func (r *TagRepository) GetAllTags() ([]models.Tag, error) {
	rows, err := r.DB.Conn.Query(`SELECT id, name, created_at FROM tags ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *TagRepository) GetTagByID(id int64) (*models.Tag, error) {
	row := r.DB.Conn.QueryRow(`SELECT id, name, created_at FROM tags WHERE id = $1`, id)
	var t models.Tag
	if err := row.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TagRepository) CreateTag(t *models.Tag) error {
	return r.DB.Conn.QueryRow(
		`INSERT INTO tags (name) VALUES ($1) RETURNING id, created_at`,
		t.Name,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TagRepository) UpdateTag(t *models.Tag) error {
	_, err := r.DB.Conn.Exec(
		`UPDATE tags SET name=$1 WHERE id=$2`,
		t.Name, t.ID,
	)
	return err
}

func (r *TagRepository) DeleteTag(id int64) error {
	_, err := r.DB.Conn.Exec(`DELETE FROM tags WHERE id = $1`, id)
	return err
}
