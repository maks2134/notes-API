package service

import (
	"errors"
	"notes-api/internal/model"
	"notes-api/internal/repository"
)

type ChecklistItemService interface {
	Create(item *model.ChecklistItem, userID int64) error
	Update(item *model.ChecklistItem, itemID int64, userID int64) error
	Delete(itemID int64, userID int64) error
}
type checklistItemService struct {
	itemRepo repository.ChecklistItemRepository
	noteRepo repository.NoteRepository // Нужен для проверки прав доступа
}

func NewChecklistItemService(itemRepo repository.ChecklistItemRepository, noteRepo repository.NoteRepository) ChecklistItemService {
	return &checklistItemService{
		itemRepo: itemRepo,
		noteRepo: noteRepo,
	}
}

// Проверка, что пользователь владеет заметкой, к которой относится чек-лист
func (s *checklistItemService) checkNoteOwnership(noteID, userID int64) error {
	note, err := s.noteRepo.GetByID(noteID, userID) // Используем GetByID, который уже проверяет userID
	if err != nil || note == nil {
		return errors.New("заметка не найдена или у вас нет к ней доступа")
	}
	return nil
}

func (s *checklistItemService) Create(item *model.ChecklistItem, userID int64) error {
	if err := s.checkNoteOwnership(item.NoteID, userID); err != nil {
		return err
	}
	return s.itemRepo.Create(item)
}

func (s *checklistItemService) Update(itemData *model.ChecklistItem, itemID, userID int64) error {
	existingItem, err := s.itemRepo.GetByID(itemID)
	if err != nil {
		return err
	}

	// 2. Проверить права доступа к заметке, к которой он относится
	if err := s.checkNoteOwnership(existingItem.NoteID, userID); err != nil {
		return err
	}

	// 3. Обновить поля и сохранить
	existingItem.Text = itemData.Text
	existingItem.Completed = itemData.Completed
	return s.itemRepo.Update(existingItem)
}

func (s *checklistItemService) Delete(itemID, userID int64) error {
	// 1. Найти существующий элемент
	existingItem, err := s.itemRepo.GetByID(itemID)
	if err != nil {
		return err
	}

	// 2. Проверить права доступа
	if err := s.checkNoteOwnership(existingItem.NoteID, userID); err != nil {
		return err
	}

	// 3. Удалить
	return s.itemRepo.Delete(itemID)
}
