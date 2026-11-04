package handler

import (
	"errors"
	"url-shortener/internal/dto"
	"url-shortener/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

var BadRequestError = map[string]string{"error": "invalid request body"}

func (h *UserHandler) HandleError(c echo.Context, err error) error {
	var serviceErr *service.ServiceError

	if errors.As(err, &serviceErr) {
		if serviceErr.Kind == service.KindInternal {
			c.Logger().Errorf("Internal error: %v", serviceErr)
			return c.JSON(serviceErr.HTTPStatus(), dto.ErrorResponse{Error: "internal server error", Kind: string(serviceErr.Kind)})
		}
		return c.JSON(serviceErr.HTTPStatus(), dto.ErrorResponse{Error: serviceErr.Message, Kind: string(serviceErr.Kind)})
	}

	c.Logger().Errorf("Unexpected error: %v", err)
	return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
}
