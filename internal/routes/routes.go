package routes

import (
	"goshop/internal/handler"
	"goshop/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, userHandler *handler.UserHandler, healthHandler *handler.HealthHandler, authMiddleware *middleware.AuthMiddleware) {
	e.GET("/health", healthHandler.Check)
	e.GET("/ready", healthHandler.Ready)

	api := e.Group("/api")

	api.POST("/register", userHandler.Register)
	api.POST("/login", userHandler.Login)

	authed := api.Group("/authed", authMiddleware.Auth)

	authed.PUT("/change_password", userHandler.ChangePassword)

	adminGroup := authed.Group("/admin", middleware.RequireRole("admin"))
	adminGroup.GET("/users", userHandler.GetAll)
	adminGroup.DELETE("/delete/:id", userHandler.Delete)
}
