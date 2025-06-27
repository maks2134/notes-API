package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"notes-api/internal/model"
	"time"
)

type PostgresNoteRepository struct {
	db *sql.DB
}

func NewPostgresNoteRepository(db *sql.DB) *PostgresNoteRepository {
	return &PostgresNoteRepository{db: db}
}

func (r *PostgresNoteRepository) Create(note *model.Note) error {
	query := `INSERT INTO notes (title, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	now := time.Now()
	err := r.db.QueryRow(query, note.Title, note.Content, note.UserID, now, now).Scan(&note.ID)
	if err != nil {
		return err
	}
	note.CreatedAt = now
	note.UpdatedAt = now
	return nil
}

func (r *PostgresNoteRepository) GetByID(id int64, userID int64) (*model.Note, error) {
	query := `SELECT id, title, content, user_id, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2;`
	note := new(model.Note)
	err := r.db.QueryRow(query, id, userID).Scan(&note.ID, &note.Title, &note.Content, &note.UserID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("заметка не найдена")
		}
		return nil, err
	}
	return note, nil
}

func (r *PostgresNoteRepository) GetAll(userID int64) ([]*model.Note, error) {
	query := `SELECT id, title, content, user_id, created_at, updated_at FROM notes WHERE user_id = $1;`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	notes := make([]*model.Note, 0)
	for rows.Next() {
		note := new(model.Note)
		if err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.UserID, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *PostgresNoteRepository) Update(note *model.Note, userID int64) error {
	query := `UPDATE notes SET title = $1, content = $2, updated_at = $3 WHERE id = $4 AND user_id = $5;`
	now := time.Now()
	res, err := r.db.Exec(query, note.Title, note.Content, now, note.ID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("заметка не найдена или у вас нет прав на её изменение")
	}
	note.UpdatedAt = now
	return nil
}

func (r *PostgresNoteRepository) Delete(id int64, userID int64) error {
	query := `DELETE FROM notes WHERE id = $1 AND user_id = $2;`
	res, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("заметка не найдена или у вас нет прав на её удаление")
	}
	return nil
}
