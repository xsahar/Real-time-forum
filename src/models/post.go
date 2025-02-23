package models

import (
    "time"
)

type Post struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    Title     string    `json:"title" validate:"required,min=3"`
    Content   string    `json:"content" validate:"required,min=10"`
    Category  string    `json:"category" validate:"required"`
    CreatedAt time.Time `json:"createdAt"`
    Comments  []Comment `json:"comments,omitempty"`
    Author    *User     `json:"author,omitempty"` // For joining user information
}