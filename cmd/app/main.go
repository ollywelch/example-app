package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ollywelch/example-app/handlers"
)

func init() {
	metadataDir := os.Getenv("METADATA_DIR")
	if metadataDir == "" {
		metadataDir = "."
	}
	if err := godotenv.Load(fmt.Sprintf("%s/.env", metadataDir)); err != nil {
		log.Printf("warning, failed to read .env file: %+v\n", err)
	}
}

func main() {
	e := echo.New()

	// *** MIDDLEWARE ***
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// *** HANDLERS ***
	server := handlers.NewServer()

	jwtConfig := echojwt.Config{
		SigningKey: server.JWTSecret,
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/login" || (c.Request().Method == http.MethodPost && c.Path() == "/users")
		},
	}
	e.Use(echojwt.WithConfig(jwtConfig))

	e.POST("/login", server.PostLogin)

	userRoutes := e.Group("/users")
	userRoutes.GET("", server.GetUsers)
	userRoutes.GET("/me", server.GetUsersMe)
	userRoutes.POST("", server.PostUsers)

	e.Logger.Fatal(e.Start(":8080"))
}
