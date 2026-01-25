package models

import "fmt"

type User struct {
	ID   int    `json:"id"`
	UserName string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
}

type LoginResponse struct {
	User    User   `json:"user"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

//TODO add signup request/response struct

func (user *User) Greet() string {
	return fmt.Sprintf("Hello, I am %s", user.UserName)
}
