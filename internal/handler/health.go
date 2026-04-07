package handler

import (
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	ready atomic.Value
}

func NewHealthHandler() *HealthHandler {
	h := &HealthHandler{}
	h.ready.Store(true)
	return h
}

func (h *HealthHandler) Check(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status":"alive"})
}

func (h *HealthHandler) Ready(c echo.Context) error{
	if h.ready.Load().(bool){
		return c.JSON(http.StatusOK, map[string]string{"status":"ready"})
	}
	return c.JSON(http.StatusServiceUnavailable,map[string]string{"status":"not ready"})
}

func (h *HealthHandler) SetReady(ready bool){
	h.ready.Store(ready)
}