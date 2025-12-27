package models

import "time"

type POST struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	TopicID   int       `json:"topic_id"`
	AuthorID  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostInput struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	TopicID  int    `json:"topic_id" binding:"required"`
	AuthorID int    `json:"author_id" binding:"required"`
}