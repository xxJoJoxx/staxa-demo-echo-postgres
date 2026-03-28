package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func HealthCheck(db *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := db.Ping(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "error", "detail": "database unreachable"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}
