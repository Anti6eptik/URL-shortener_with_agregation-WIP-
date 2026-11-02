package middleware

import (
	"github.com/labstack/echo/v4"
)

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get("role") != role {
				return c.JSON(403, "forbidden")
			}
			return next(c)
		}
	}
}
