package handlers

import (
	"github.com/ollywelch/example-app/db"
)

type Server struct {
	JWTSecret []byte
	store     db.Store
}

func NewServer() *Server {
	return &Server{
		JWTSecret: []byte("supersecret"),
		store:     &db.InMemoryStore{},
	}
}
