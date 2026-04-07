package dto

import (
	"goshop/internal/domain"
	"time"
)

type RegisterRequest struct {
	Email    string `param:"email" query:"email" json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `param:"email" query:"email" json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	Email        string `param:"email" query:"email" json:"email"`
	Old_Password string `json:"old_password"`
	New_Password string `json:"new_password"`
}

type GetAllRequest struct {
	Page  int `param:"page" query:"page"`
	Limit int `param:"limit" query:"limit"`
}

type DeleteRequest struct {
	ID int `param:"id" query:"id"`
}

type GetAllResponse struct {
	Users      []domain.User `json:"users"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	TotalPages int           `json:"total_pages"`
}

type UserResponse struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Kind  string `json:"kind,omitempty"`
}
