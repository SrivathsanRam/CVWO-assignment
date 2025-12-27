package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	AuthorID  int       `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommentInput struct {
	PostID   int    `json:"post_id" binding:"required"`
	AuthorID int    `json:"author_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}