package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"goshop/internal/config"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	//db pool
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Println("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		fmt.Println("failed to ping database", "error", err)
		os.Exit(1)
	}
	fmt.Println("connected to database")

	//DI

	e := echo.New()

	errChan := make(chan error)
	address := ":" + cfg.Port

	go func() {
		fmt.Println("server is started")
		if err := e.Start(address); err != nil {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		fmt.Println("\nShutting down server...")

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			fmt.Printf("Error during shutdown: %v\n", err)
		}

	case err := <-errChan:
		fmt.Printf("Server error: %v\n", err)
	}

	fmt.Println("Server stopped")
}
