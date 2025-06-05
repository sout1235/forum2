package entity

import "time"

type Comment struct {
	ID        int64     `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int64     `json:"author_id" db:"author_id"`
	TopicID   int64     `json:"topic_id" db:"topic_id"`
	ParentID  *int64    `json:"parent_id,omitempty" db:"parent_id"`
	Likes     int       `json:"likes" db:"likes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Author    *User     `json:"author" db:"-"`
}
