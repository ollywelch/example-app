package types

import "github.com/google/uuid"

type UserRead struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
