package main

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	RegisterHandlers(e, &Server{})
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":8080"))
}
