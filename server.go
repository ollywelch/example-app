package main

import "github.com/labstack/echo/v4"

type Server struct{}

// (GET /ping)
func (s *Server) GetPing(ctx echo.Context) error {
	value := "pong"
	return ctx.JSON(200, ValueResponse{Value: &value})
}
