package entity

import "time"

type Topic struct {
	ID           int64      `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Content      string     `json:"content" db:"content"`
	AuthorID     int64      `json:"author_id" db:"author_id"`
	CategoryID   int64      `json:"category_id" db:"category_id"`
	Views        int        `json:"views" db:"views"`
	CommentCount int        `json:"comment_count" db:"comment_count"`
	Comments     []*Comment `json:"comments,omitempty"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	Author       *User      `json:"author,omitempty"`
}
