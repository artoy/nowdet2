package models

import "time"

type User struct {
	ID        string    `json:"id" spanner:"id"`
	Name      string    `json:"name" spanner:"name"`
	Email     string    `json:"email" spanner:"email"`
	CreatedAt time.Time `json:"created_at" spanner:"created_at"`
	UpdatedAt time.Time `json:"updated_at" spanner:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
