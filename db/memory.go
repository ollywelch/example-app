package db

import (
	"sync"

	"github.com/google/uuid"
	"github.com/ollywelch/example-app/types"
)

type InMemoryStore struct{}

var (
	users = []types.UserRead{
		{
			Username: "Olly",
			Password: "password",
			Id:       uuid.New(),
		},
	}
	mu = sync.Mutex{}
)

func (s *InMemoryStore) GetUsers() []types.UserRead {
	mu.Lock()
	defer mu.Unlock()
	return users
}

func (s *InMemoryStore) GetUserByName(name string) *types.UserRead {
	mu.Lock()
	defer mu.Unlock()
	for _, user := range users {
		if user.Username == name {
			return &user
		}
	}
	return nil
}

func (s *InMemoryStore) GetUserByID(id uuid.UUID) *types.UserRead {
	mu.Lock()
	defer mu.Unlock()
	for _, user := range users {
		if user.Id == id {
			return &user
		}
	}
	return nil
}

func (s *InMemoryStore) CreateUser(u types.UserCreate) *types.UserRead {
	if exists := s.GetUserByName(u.Username); exists != nil {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	user := types.UserRead{Id: uuid.New(), Username: u.Username, Password: u.Password}
	users = append(users, user)
	return &user
}
