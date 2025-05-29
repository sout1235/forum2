package models

type Comment struct {
	ID        int64  `db:"id" json:"id"`
	Content   string `db:"content" json:"content"`
	AuthorID  int64  `db:"author_id" json:"author_id"`
	TopicID   int64  `db:"topic_id" json:"topic_id"`
	ParentID  *int64 `db:"parent_id" json:"parent_id,omitempty"`
	Likes     int    `db:"likes" json:"likes"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
