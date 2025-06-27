package model

import "time"

type ChecklistItem struct {
	ID        int64     `json:"id" example:"101"`
	Text      string    `json:"text" example:"Купить молоко"`
	Completed bool      `json:"completed" example:"false"`
	Style     TextStyle `json:"style,omitempty" example:"italic"` // <-- НОВОЕ ПОЛЕ
	NoteID    int64     `json:"note_id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
