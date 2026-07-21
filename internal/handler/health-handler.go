package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	dbPingCheck func(ctx context.Context) error
}

func NewHealthHandler(dbPingCheck func(ctx context.Context) error) *HealthHandler {
	return &HealthHandler{dbPingCheck: dbPingCheck}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	if err := h.dbPingCheck(c.Request().Context()); err != nil {
		return c.JSON(http.StatusServiceUnavailable, echo.Map{
			"status": "DOWN",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "UP",
		"db":     "ok",
	})
}