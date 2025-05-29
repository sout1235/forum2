package repository

import (
	"backend/services/forum/models"
)

type CommentRepository struct {
	DB *DB
}

func NewCommentRepository(db *DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (r *CommentRepository) GetAllCommentsByTopic(topicID int64) ([]models.Comment, error) {
	rows, err := r.DB.Conn.Query(`SELECT id, content, author_id, topic_id, parent_id, likes, created_at, updated_at FROM comments WHERE topic_id = $1 ORDER BY created_at ASC`, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.AuthorID, &c.TopicID, &c.ParentID, &c.Likes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepository) GetCommentByID(id int64) (*models.Comment, error) {
	row := r.DB.Conn.QueryRow(`SELECT id, content, author_id, topic_id, parent_id, likes, created_at, updated_at FROM comments WHERE id = $1`, id)
	var c models.Comment
	if err := row.Scan(&c.ID, &c.Content, &c.AuthorID, &c.TopicID, &c.ParentID, &c.Likes, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CommentRepository) CreateComment(c *models.Comment) error {
	return r.DB.Conn.QueryRow(
		`INSERT INTO comments (content, author_id, topic_id, parent_id, likes) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`,
		c.Content, c.AuthorID, c.TopicID, c.ParentID, c.Likes,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CommentRepository) UpdateComment(c *models.Comment) error {
	_, err := r.DB.Conn.Exec(
		`UPDATE comments SET content=$1, parent_id=$2, likes=$3, updated_at=NOW() WHERE id=$4`,
		c.Content, c.ParentID, c.Likes, c.ID,
	)
	return err
}

func (r *CommentRepository) DeleteComment(id int64) error {
	_, err := r.DB.Conn.Exec(`DELETE FROM comments WHERE id = $1`, id)
	return err
}
