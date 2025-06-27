package repository

import "notes-api/internal/model"

type NoteRepository interface {
	Create(note *model.Note) error
	GetByID(id int64) (*model.Note, error)
	GetAll() ([]*model.Note, error)
	Update(note *model.Note) error
	Delete(id int64) error
}
