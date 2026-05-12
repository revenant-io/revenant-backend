package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserPublic struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
