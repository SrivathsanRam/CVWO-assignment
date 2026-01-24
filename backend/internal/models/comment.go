package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"username"` //Name of the author (need to use SQL join on users table)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	PostID   int    `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	// TopicID is not included here we will get it from post context
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}