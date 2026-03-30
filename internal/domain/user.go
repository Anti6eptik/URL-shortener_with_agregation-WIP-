package domain

import (
	"time"
)

type User struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Password_hash string    `json:"-"`
	Created_at    time.Time `json:"created_at"`
}

type UserResponse struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
}
