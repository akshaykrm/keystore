package user

import "time"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Email    string
	Password string
	Name     string
}

type UpdateUserInput struct {
	Email string
	Name  string
}
