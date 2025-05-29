package models

type Topic struct {
	ID         int64  `db:"id" json:"id"`
	Title      string `db:"title" json:"title"`
	Content    string `db:"content" json:"content"`
	AuthorID   int64  `db:"author_id" json:"author_id"`
	CategoryID int64  `db:"category_id" json:"category_id"`
	Views      int    `db:"views" json:"views"`
	CreatedAt  string `db:"created_at" json:"created_at"`
	UpdatedAt  string `db:"updated_at" json:"updated_at"`
}
