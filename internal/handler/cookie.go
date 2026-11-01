package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetAccessTokenCookie(c echo.Context, token string, expiresAfter time.Duration) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(expiresAfter),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)
}

func GetAccessTokenCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
