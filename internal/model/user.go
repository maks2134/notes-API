package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
