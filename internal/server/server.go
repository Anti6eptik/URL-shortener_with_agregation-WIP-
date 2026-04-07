package server

import (
	"context"
	"fmt"
	"goshop/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo   *echo.Echo
	pool   *pgxpool.Pool
	config *config.Config
}

func NewServer(e *echo.Echo, pool *pgxpool.Pool, cfg *config.Config) *Server {
	return &Server{echo: e, pool: pool, config: cfg}
}

func (s *Server) Start() error {

	serverErrors := make(chan error, 1)

	go func() {
		fmt.Printf("Server starting on %s", s.config.Port)
		serverErrors <- s.echo.Start(":" + s.config.Port)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		fmt.Printf("Recieved signal: %v, starting graceful shutdown\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.echo.Shutdown(ctx); err != nil {
			s.echo.Close()
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		s.pool.Close()
		fmt.Println("Database connection closed")
		fmt.Println("Server gracefully stopped")
	}
	return nil
}
