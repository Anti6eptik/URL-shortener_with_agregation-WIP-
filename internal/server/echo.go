package server

import (
	"url-shortener/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho(cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeot
	e.Server.IdleTimeout = cfg.IdleTimeout

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.HTTPErrorHandler = customHTTPErrorHandler

	return e
}

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	c.JSON(500, map[string]string{"error": "Internal server error"})
}
