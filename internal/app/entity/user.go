package entity

import (
	"net/http"
	"time"
)

const STORE_FILE string = `./data/users.json`

type User struct {
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}

type UserList map[string]User

type UserStore struct {
	Increment int      `json:"increment"`
	List      UserList `json:"list"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error {
	return nil
}

type CreateUserRequest struct {
	DisplayName string `json:"display_name" validate:"required"`
	Email       string `json:"email" validate:"email,required"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error { return nil }
