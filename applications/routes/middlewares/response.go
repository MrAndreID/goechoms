package middlewares

import (
	"github.com/labstack/echo/v4"
)

func (cm *CustomMiddleware) NoCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")

		return next(c)
	}
}
