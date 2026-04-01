package model

import "time"

type TodoItem struct {
	ID        string    `json:"id"        dynamodbav:"id"`
	UserID    string    `json:"userId"    dynamodbav:"userId"`
	Title     string    `json:"title"     dynamodbav:"title"`
	Done      bool      `json:"done" 	  dynamodbav:"done"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

type CreateTodoRequest struct {
	UserID string `json:"userId"`
	Title  string `json:"title"`
}

type UpdateTodoRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}
