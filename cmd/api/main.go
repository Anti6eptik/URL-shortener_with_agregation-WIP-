package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"

	"url-shortener/internal/config"
	"url-shortener/internal/database"
	"url-shortener/internal/handler"
	"url-shortener/internal/middleware"
	repository "url-shortener/internal/repositories/postgres"
	"url-shortener/internal/routes"
	"url-shortener/internal/server"
	service "url-shortener/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	c := dig.New()

	c.Provide(config.Load)
	c.Provide(func() context.Context { return context.Background() })
	c.Provide(database.NewPgxPool)
	c.Provide(server.NewEcho)
	c.Provide(repository.NewUserRepository)
	c.Provide(service.NewPasswordService)
	c.Provide(service.NewUserService)
	c.Provide(service.NewTokenService)
	c.Provide(middleware.NewAuthMiddleware)
	c.Provide(handler.NewUserHandler)
	c.Provide(handler.NewHealthHandler)
	c.Provide(server.NewServer)

	if err := c.Invoke(startServer); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}

}

func startServer(
	srv *server.Server,
	e *echo.Echo,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	authMiddleware *middleware.AuthMiddleware,
	cfg *config.Config,
) error {
	routes.Register(e, userHandler, healthHandler, authMiddleware)
	return srv.Start()
}
