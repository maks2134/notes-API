package repository

import (
	"database/sql"
	"errors"
	"notes-api/internal/model"
	"time"
)

type PostgresChecklistItemRepository struct {
	db *sql.DB
}

func NewPostgresChecklistItemRepository(db *sql.DB) *PostgresChecklistItemRepository {
	return &PostgresChecklistItemRepository{db: db}
}

func (r *PostgresChecklistItemRepository) Create(item *model.ChecklistItem) error {
	query := `INSERT INTO checklist_items (text, note_id) VALUES ($1, $2) RETURNING id, completed, created_at, updated_at;`
	now := time.Now()
	err := r.db.QueryRow(query, item.Text, item.NoteID).Scan(&item.ID, &item.Completed, &item.CreatedAt, &item.UpdatedAt)
	item.UpdatedAt = now
	return err
}

func (r *PostgresChecklistItemRepository) GetByNoteID(noteID int64) ([]*model.ChecklistItem, error) {
	query := `SELECT id, text, completed, note_id, created_at, updated_at FROM checklist_items WHERE note_id = $1 ORDER BY created_at ASC;`
	rows, err := r.db.Query(query, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.ChecklistItem
	for rows.Next() {
		item := new(model.ChecklistItem)
		if err := rows.Scan(&item.ID, &item.Text, &item.Completed, &item.NoteID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *PostgresChecklistItemRepository) GetByID(itemID int64) (*model.ChecklistItem, error) {
	query := `SELECT id, text, completed, note_id, created_at, updated_at FROM checklist_items WHERE id = $1;`
	item := new(model.ChecklistItem)
	err := r.db.QueryRow(query, itemID).Scan(&item.ID, &item.Text, &item.Completed, &item.NoteID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("элемент чек-листа не найден")
		}
		return nil, err
	}
	return item, nil
}

func (r *PostgresChecklistItemRepository) Update(item *model.ChecklistItem) error {
	query := `UPDATE checklist_items SET text = $1, completed = $2, updated_at = $3 WHERE id = $4;`
	item.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, item.Text, item.Completed, item.UpdatedAt, item.ID)
	return err
}

func (r *PostgresChecklistItemRepository) Delete(itemID int64) error {
	query := `DELETE FROM checklist_items WHERE id = $1;`
	_, err := r.db.Exec(query, itemID)
	return err
}
