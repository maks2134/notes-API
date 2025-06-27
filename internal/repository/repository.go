package repository

import "notes-api/internal/model"

type NoteRepository interface {
	Create(note *model.Note) error
	GetByID(id int64, userID int64) (*model.Note, error)
	GetAll(userID int64) ([]*model.Note, error)
	Update(note *model.Note, userID int64) error
	Delete(id int64, userID int64) error
}
