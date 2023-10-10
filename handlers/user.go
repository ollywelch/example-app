package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ollywelch/example-app/types"
)

func (s *Server) GetUsers(ctx echo.Context) error {
	users := s.store.GetUsers()
	return ctx.JSON(http.StatusOK, users)
}

func (s *Server) PostUsers(ctx echo.Context) error {
	u := &types.UserCreate{}
	if err := ctx.Bind(u); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid JSON inputs")
	}
	if u.Username == "" || u.Password == "" {
		return ctx.JSON(http.StatusBadRequest, "username and password must both be specified")
	}
	user := s.store.CreateUser(*u)
	if user == nil {
		return ctx.JSON(http.StatusInternalServerError, "internal server error")
	}
	return ctx.JSON(http.StatusOK, *user)
}

func (s *Server) GetUsersMe(ctx echo.Context) error {
	username := userIDFromToken(ctx)
	user := s.store.GetUserByName(username)
	if user != nil {
		return ctx.JSON(http.StatusOK, user)
	}
	return fmt.Errorf("failed to find user with name %s in the DB", username)
}

func userIDFromToken(ctx echo.Context) string {
	token := ctx.Get("user").(*jwt.Token)
	user, _ := token.Claims.GetSubject()
	return user
}
