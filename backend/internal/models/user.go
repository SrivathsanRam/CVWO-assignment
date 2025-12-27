package models

import "fmt"

type User struct {
	ID   int    `json:"id"`
	UserName string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
}

func (user *User) Greet() string {
	return fmt.Sprintf("Hello, I am %s", user.UserName)
}
