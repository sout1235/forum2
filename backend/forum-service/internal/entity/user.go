package entity

type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Avatar   string `json:"avatar" db:"avatar"`
}
