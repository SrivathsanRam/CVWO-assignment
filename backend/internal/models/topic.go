package models

import "time"

type Topic struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int       `json:"created_by"` //ID of the creator
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTopicInput struct {
	Name string `json:"name" binding:"required"`
}
