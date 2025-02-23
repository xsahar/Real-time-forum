package models

import (
    "time"
)

type Comment struct {
    ID        string    `json:"id"`
    PostID    string    `json:"postId"`
    UserID    string    `json:"userId"`
    Content   string    `json:"content" validate:"required,min=1"`
    CreatedAt time.Time `json:"createdAt"`
    Author    *User     `json:"author,omitempty"`
}