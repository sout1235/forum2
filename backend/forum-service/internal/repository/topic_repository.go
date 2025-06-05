package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
)

type TopicRepository interface {
	CreateTopic(ctx context.Context, topic *entity.Topic) error
	GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error)
	GetAllTopics(ctx context.Context) ([]*entity.Topic, error)
	UpdateTopic(ctx context.Context, topic *entity.Topic) error
	DeleteTopic(ctx context.Context, id int64) error
	UpdateCommentCount(ctx context.Context, topicID int64) error
}

type topicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	log.Printf("Creating new topic: Title=%s, AuthorID=%d, CategoryID=%d",
		topic.Title, topic.AuthorID, topic.CategoryID)

	err := r.db.QueryRowContext(ctx,
		`INSERT INTO topics (title, content, author_id, category_id, views, comment_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		topic.Title, topic.Content, topic.AuthorID, topic.CategoryID, topic.Views, 0, topic.CreatedAt, topic.UpdatedAt,
	).Scan(&topic.ID)

	if err != nil {
		log.Printf("Error creating topic: %v", err)
		return fmt.Errorf("failed to create topic: %w", err)
	}

	log.Printf("Successfully created topic with ID=%d", topic.ID)
	return nil
}

func (r *topicRepository) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	topic := &entity.Topic{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, content, author_id, category_id, views, comment_count, created_at, updated_at
		FROM topics
		WHERE id = $1`,
		id,
	).Scan(
		&topic.ID, &topic.Title, &topic.Content, &topic.AuthorID, &topic.CategoryID,
		&topic.Views, &topic.CommentCount, &topic.CreatedAt, &topic.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}
	return topic, nil
}

func (r *topicRepository) GetAllTopics(ctx context.Context) ([]*entity.Topic, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, content, author_id, category_id, views, comment_count, created_at, updated_at
		FROM topics
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query topics: %w", err)
	}
	defer rows.Close()

	var topics []*entity.Topic
	for rows.Next() {
		topic := &entity.Topic{}
		err := rows.Scan(
			&topic.ID,
			&topic.Title,
			&topic.Content,
			&topic.AuthorID,
			&topic.CategoryID,
			&topic.Views,
			&topic.CommentCount,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating topics: %w", err)
	}

	return topics, nil
}

func (r *topicRepository) UpdateTopic(ctx context.Context, topic *entity.Topic) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE topics SET title = $1, content = $2, category_id = $3, updated_at = NOW() 
		WHERE id = $4`,
		topic.Title, topic.Content, topic.CategoryID, topic.ID,
	)
	return err
}

func (r *topicRepository) DeleteTopic(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM topics WHERE id = $1`, id)
	return err
}

func (r *topicRepository) UpdateCommentCount(ctx context.Context, topicID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE topics 
		SET comment_count = (
			SELECT COUNT(*) 
			FROM comments 
			WHERE topic_id = $1
		)
		WHERE id = $1`, topicID)
	return err
}

func (r *topicRepository) GetDB() *sql.DB {
	return r.db
}
