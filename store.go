package main

import "sync"

type Store interface {
	GetUsers() []UserResponse
	GetUserByName(string) *UserResponse
	GetUserByID(int) *UserResponse
}

type InMemoryStore struct{}

var (
	users = []UserResponse{
		{
			Id:       1,
			Username: "Olly",
		},
	}
	mu = sync.Mutex{}
)

func (s *InMemoryStore) GetUsers() []UserResponse {
	mu.Lock()
	defer mu.Unlock()
	return users
}

func (s *InMemoryStore) GetUserByName(name string) *UserResponse {
	mu.Lock()
	defer mu.Unlock()
	for _, user := range users {
		if user.Username == name {
			return &user
		}
	}
	return nil
}

func (s *InMemoryStore) GetUserByID(id int) *UserResponse {
	mu.Lock()
	defer mu.Unlock()
	for _, user := range users {
		if user.Id == id {
			return &user
		}
	}
	return nil
}
