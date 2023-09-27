package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Server struct {
	jwtSecret []byte
	store     Store
}

func NewServer() *Server {
	return &Server{
		jwtSecret: []byte("supersecret"),
		store:     &InMemoryStore{},
	}
}

// (GET /ping)
func (s *Server) GetPing(ctx echo.Context) error {
	value := "pong"
	return ctx.JSON(200, ValueResponse{Value: &value})
}

func (s *Server) PostLogin(ctx echo.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	user := s.store.GetUserByName(username)

	if user == nil || password != "password" {
		status := http.StatusUnauthorized
		return ctx.JSON(status, ErrorResponse{Status: status, Message: "Invalid login details"})
	}

	claims := jwt.RegisteredClaims{
		Subject:   user.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return err
	}

	return ctx.JSON(200, LoginResponse{Token: t})
}

func (s *Server) GetUsersMe(ctx echo.Context) error {
	username := userIDFromToken(ctx)
	user := s.store.GetUserByName(username)
	if user != nil {
		return ctx.JSON(http.StatusOK, user)
	}
	return fmt.Errorf("failed to find user with name %s in the DB", username)
}

func (s *Server) GetUsers(ctx echo.Context) error {
	users := s.store.GetUsers()
	return ctx.JSON(http.StatusOK, users)
}

func userIDFromToken(ctx echo.Context) string {
	token := ctx.Get("user").(*jwt.Token)
	user, _ := token.Claims.GetSubject()
	return user
}
