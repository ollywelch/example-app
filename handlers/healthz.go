package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetHealthz(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
