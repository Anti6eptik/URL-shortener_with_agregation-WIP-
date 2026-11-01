package routes

import (
	"goshop/internal/handler"
	"goshop/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, userHandler *handler.UserHandler, healthHandler *handler.HealthHandler) {
	e.GET("/health", healthHandler.Check)
	e.GET("/ready", healthHandler.Ready)

	api := e.Group("/api")

	api.POST("/register", userHandler.Register)
	api.POST("/login", userHandler.Login)
	api.PUT("/change_password", userHandler.ChangePassword)

	adminGroup := api.Group("/admin", middleware.RequireRole("admin"))
	adminGroup.GET("/users", userHandler.GetAll)
	adminGroup.DELETE("/delete/:id", userHandler.Delete)
}
