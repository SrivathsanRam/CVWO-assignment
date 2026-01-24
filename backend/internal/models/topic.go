package models

import "time"

type Topic struct {
	ID        	int       	`json:"id"`
	Title      	string    	`json:"title"`
	Description string    	`json:"description"`
	UserID 		int       	`json:"user_id"` //ID of the creator
	UserName 	string		`json:"username"` //Name of the creator (need to use SQL join on users table)
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

type CreateTopicRequest struct {
	Title 		string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type UpdateTopicRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
}
