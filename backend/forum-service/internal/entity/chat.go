package entity

import "time"

// ChatMessage represents a chat message in the system
type ChatMessage struct {
	ID             int64
	Content        string
	AuthorID       int64
	AuthorUsername string
	CreatedAt      time.Time
	ExpiresAt      time.Time
}
