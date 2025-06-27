package service

import (
	"notes-api/internal/model"
	"notes-api/internal/repository"
)

type NoteService interface {
	CreateNote(note *model.Note) error
	GetNoteByID(id int64, userID int64) (*model.Note, error)
	GetAllNotes(userID int64) ([]*model.Note, error)
	UpdateNote(note *model.Note, userID int64) error
	DeleteNote(id int64, userID int64) error
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) CreateNote(note *model.Note) error {
	return s.repo.Create(note)
}

func (s *noteService) GetNoteByID(id int64, userID int64) (*model.Note, error) {
	return s.repo.GetByID(id, userID)
}

func (s *noteService) GetAllNotes(userID int64) ([]*model.Note, error) {
	return s.repo.GetAll(userID)
}

func (s *noteService) UpdateNote(note *model.Note, userID int64) error {
	return s.repo.Update(note, userID)
}

func (s *noteService) DeleteNote(id int64, userID int64) error {
	return s.repo.Delete(id, userID)
}
