package main

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// *** MIDDLEWARE ***
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// *** HANDLERS ***
	server := NewServer()
	RegisterHandlersWithBaseURL(e, server, "/api")

	jwtConfig := echojwt.Config{
		SigningKey: server.jwtSecret,
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/api/login" || (c.Request().Method == http.MethodPost && c.Path() == "/api/users")
		},
	}
	e.Use(echojwt.WithConfig(jwtConfig))

	e.Logger.Fatal(e.Start(":8080"))
}
