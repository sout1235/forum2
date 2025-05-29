package repository

import (
	"backend/services/forum/models"
)

type CategoryRepository struct {
	DB *DB
}

func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.DB.Conn.Query(`SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(id int64) (*models.Category, error) {
	row := r.DB.Conn.QueryRow(`SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`, id)
	var c models.Category
	if err := row.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) CreateCategory(c *models.Category) error {
	return r.DB.Conn.QueryRow(
		`INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		c.Name, c.Description,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CategoryRepository) UpdateCategory(c *models.Category) error {
	_, err := r.DB.Conn.Exec(
		`UPDATE categories SET name=$1, description=$2, updated_at=NOW() WHERE id=$3`,
		c.Name, c.Description, c.ID,
	)
	return err
}

func (r *CategoryRepository) DeleteCategory(id int64) error {
	_, err := r.DB.Conn.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}
