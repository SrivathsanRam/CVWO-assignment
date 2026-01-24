package models

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	TopicID   int       `json:"topic_id"`
	TopicName string    `json:"topic_name"` //Name of the topic (need to use SQL join on topics table)
	UserID	  int       `json:"user_id"`
	UserName  string    `json:"username"` //Name of the author (need to use SQL join on users table)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	TopicID  int    `json:"topic_id" binding:"required"`
	// Learning point: UserID is not included here as it will be derived from the authenticated user context
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	// TopicID is not included here intentionally
}