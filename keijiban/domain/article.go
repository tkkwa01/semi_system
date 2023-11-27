package domain

import (
	"time"
)

type Article struct {
	ID        uint      `json:"id"`
	AuthorID  uint      `json:"author_id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
