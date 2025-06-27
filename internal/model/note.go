package model

import "time"

type Note struct {
	ID        int64     `json:"id" example:"1"`
	Title     string    `json:"title" example:"My First Note"`
	Content   string    `json:"content" example:"This is the content of my first note."`
	UserID    int64     `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
