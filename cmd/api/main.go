package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"

	"goshop/internal/config"
	"goshop/internal/database"
	"goshop/internal/handler"
	"goshop/internal/repository/postgres"
	"goshop/internal/routes"
	"goshop/internal/server"
	"goshop/internal/service"
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
	c.Provide(postgres.NewUserRepository)
	c.Provide(service.NewPasswordService)
	c.Provide(service.NewUserService)
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
	cfg *config.Config,
) error {
	routes.Register(e, userHandler, healthHandler)
	return srv.Start()
}
