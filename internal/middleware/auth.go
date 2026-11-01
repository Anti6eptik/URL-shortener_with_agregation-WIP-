package middleware

import (
	"goshop/internal/handler"
	service "goshop/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	TokenService *service.TokenService
}

func NewAuthMiddleware(tokenService *service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{TokenService: tokenService}
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := handler.GetAccessTokenCookie(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		claims, err := m.TokenService.ValidateAccessToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		return next(c)
	}
}
