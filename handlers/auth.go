package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ollywelch/example-app/types"
)

func (s *Server) PostLogin(ctx echo.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	user := s.store.GetUserByName(username)

	if user == nil || password != user.Password {
		status := http.StatusUnauthorized
		return ctx.JSON(status, types.ErrorResponse{Status: status, Message: "Invalid login details"})
	}

	claims := jwt.RegisteredClaims{
		Subject:   user.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return err
	}

	return ctx.JSON(200, types.LoginResponse{Token: t})
}
