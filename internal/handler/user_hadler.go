package handler

import (
	"goshop/internal/dto"
	service "goshop/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(
	UserService *service.UserService,

) *UserHandler {
	return &UserHandler{
		UserService: UserService,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	user, token, err := h.UserService.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return h.HandleError(c, err)
	}

	SetAccessTokenCookie(c, token, h.UserService.TokenService.Cfg.AccessTokenExpiry)

	return c.JSON(http.StatusCreated, dto.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Created_at: user.Created_at,
	})

}

func (h *UserHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	user, token, err := h.UserService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return h.HandleError(c, err)
	}

	SetAccessTokenCookie(c, token, h.UserService.TokenService.Cfg.AccessTokenExpiry)

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Created_at: user.Created_at,
	})
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	token, err := h.UserService.ChangePassword(c.Request().Context(), req.Email, req.Old_Password, req.New_Password)
	if err != nil {
		return h.HandleError(c, err)
	}

	SetAccessTokenCookie(c, token, h.UserService.TokenService.Cfg.AccessTokenExpiry)

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	var req dto.GetAllRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	users, total, totalPages, err := h.UserService.GetAll(c.Request().Context(), req.Page, req.Limit)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.GetAllResponse{
		Users:      users,
		Total:      total,
		Page:       req.Page,
		TotalPages: totalPages,
	})
}

func (h *UserHandler) Delete(c echo.Context) error {
	var req dto.DeleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	err := h.UserService.Delete(c.Request().Context(), req.ID)
	if err != nil {
		return h.HandleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
