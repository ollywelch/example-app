package main

import "sync"

type Store interface {
	GetUsers() []UserResponse
	GetUserByName(string) *UserResponse
	GetUserByID(int) *UserResponse
	CreateUser(NewUser) *UserResponse
}

type InMemoryStore struct{}

var (
	users = []UserResponse{
		{
			Id:       1,
			Username: "Olly",
			Password: "password",
		},
	}
	counter = 1
	mu      = sync.Mutex{}
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

func (s *InMemoryStore) CreateUser(u NewUser) *UserResponse {
	if exists := s.GetUserByName(u.Username); exists != nil {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	user := UserResponse{Id: counter + 1, Username: u.Username, Password: u.Password}
	users = append(users, user)
	counter++
	return &user
}
