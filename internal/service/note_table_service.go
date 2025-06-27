package service

import (
	"errors"
	"notes-api/internal/model"
	"notes-api/internal/repository"
)

type NoteTableService interface {
	CreateTable(req *model.CreateNoteTableRequest, noteID, userID int64) (*model.NoteTable, error)
	AddRow(req *model.AddTableRowRequest, tableID, userID int64) (*model.TableRow, error)
}

type noteTableServiceImpl struct {
	tableRepo repository.NoteTableRepository
	noteRepo  repository.NoteRepository
}

func NewNoteTableService(tableRepo repository.NoteTableRepository, noteRepo repository.NoteRepository) NoteTableService {
	return &noteTableServiceImpl{tableRepo: tableRepo, noteRepo: noteRepo}
}

func (s *noteTableServiceImpl) checkNoteOwnership(noteID, userID int64) error {
	note, err := s.noteRepo.GetByID(noteID, userID)
	if err != nil || note == nil {
		return errors.New("заметка не найдена или у вас нет к ней доступа")
	}
	return nil
}

func (s *noteTableServiceImpl) CreateTable(req *model.CreateNoteTableRequest, noteID, userID int64) (*model.NoteTable, error) {
	if err := s.checkNoteOwnership(noteID, userID); err != nil {
		return nil, err
	}

	if len(req.Columns) == 0 {
		return nil, errors.New("таблица должна иметь хотя бы одну колонку")
	}

	tx, err := s.tableRepo.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	table := &model.NoteTable{
		NoteID: noteID,
		Title:  req.Title,
	}

	if err := s.tableRepo.Create(tx, table); err != nil {
		return nil, err
	}

	if err := s.tableRepo.CreateColumns(tx, table.ID, req.Columns); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	for i, colName := range req.Columns {
		table.Columns = append(table.Columns, &model.TableColumn{Name: colName, Position: i})
	}

	return table, nil
}

func (s *noteTableServiceImpl) AddRow(req *model.AddTableRowRequest, tableID, userID int64) (*model.TableRow, error) {
	return s.tableRepo.AddRow(tableID, req.Cells)
}
