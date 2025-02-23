package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname" validate:"required,min=3"`
	Age       int       `json:"age" validate:"required,min=13"`
	Gender    string    `json:"gender" validate:"required"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"` // Make sure password is included in JSON
	CreatedAt time.Time `json:"createdAt"`
}
