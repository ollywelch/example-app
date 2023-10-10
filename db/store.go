package db

import (
	"github.com/google/uuid"
	"github.com/ollywelch/example-app/types"
)

type Store interface {
	GetUsers() []types.UserRead
	GetUserByName(string) *types.UserRead
	GetUserByID(uuid.UUID) *types.UserRead
	CreateUser(types.UserCreate) *types.UserRead
}
