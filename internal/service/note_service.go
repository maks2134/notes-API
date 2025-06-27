package service

import (
	"notes-api/internal/model"
	"notes-api/internal/repository"
)

type NoteService interface {
	CreateNote(note *model.Note) error
	GetNoteByID(id int64) (*model.Note, error)
	GetAllNotes() ([]*model.Note, error)
	UpdateNote(note *model.Note) error
	DeleteNote(id int64) error
}

type noteService struct {
	repo repository.NoteRepository
}

// NewNoteService возвращает новый экземпляр NoteService.
func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) CreateNote(note *model.Note) error {
	return s.repo.Create(note)
}

func (s *noteService) GetNoteByID(id int64) (*model.Note, error) {
	return s.repo.GetByID(id)
}

func (s *noteService) GetAllNotes() ([]*model.Note, error) {
	return s.repo.GetAll()
}

func (s *noteService) UpdateNote(note *model.Note) error {
	return s.repo.Update(note)
}

func (s *noteService) DeleteNote(id int64) error {
	return s.repo.Delete(id)
}
