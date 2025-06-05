package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
)

type CommentRepository interface {
	GetCommentsByTopic(ctx context.Context, topicID int64) ([]*entity.Comment, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
	CreateComment(ctx context.Context, comment *entity.Comment) error
	UpdateComment(ctx context.Context, comment *entity.Comment) error
	DeleteComment(ctx context.Context, id int64) error
	LikeComment(ctx context.Context, id int64) error
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) GetCommentsByTopic(ctx context.Context, topicID int64) ([]*entity.Comment, error) {
	query := `
		SELECT c.id, c.content, c.author_id, c.topic_id, c.parent_id, c.likes, c.created_at, c.updated_at,
		u.username, u.avatar
		FROM comments c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.topic_id = $1
		ORDER BY c.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		comment := &entity.Comment{
			Author: &entity.User{},
		}
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.AuthorID,
			&comment.TopicID,
			&comment.ParentID,
			&comment.Likes,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author.Username,
			&comment.Author.Avatar,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *commentRepository) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	query := `
		SELECT c.id, c.content, c.author_id, c.topic_id, c.parent_id, c.likes, c.created_at, c.updated_at,
		u.username, u.avatar
		FROM comments c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.id = $1
	`
	comment := &entity.Comment{
		Author: &entity.User{},
	}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.AuthorID,
		&comment.TopicID,
		&comment.ParentID,
		&comment.Likes,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Author.Username,
		&comment.Author.Avatar,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	return comment, nil
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *entity.Comment) error {
	query := `
		INSERT INTO comments (content, author_id, topic_id, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	err := r.db.QueryRowContext(ctx,
		query,
		comment.Content,
		comment.AuthorID,
		comment.TopicID,
		comment.ParentID,
		comment.CreatedAt,
		comment.UpdatedAt,
	).Scan(&comment.ID)

	if err != nil {
		log.Printf("Error creating comment: %v", err)
		return err
	}

	log.Printf("Comment created successfully with ID: %d for topic ID: %d", comment.ID, comment.TopicID)
	return nil
}

func (r *commentRepository) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	query := `
		UPDATE comments 
		SET content = $1, updated_at = $2
		WHERE id = $3
	`
	comment.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query, comment.Content, comment.UpdatedAt, comment.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, id int64) error {
	query := `DELETE FROM comments WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

func (r *commentRepository) LikeComment(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE comments
		SET likes = likes + 1
		WHERE id = $1
	`, id)
	return err
}
