package model

import "time"

type NoteTable struct {
	ID        int64          `json:"id"`
	NoteID    int64          `json:"note_id"`
	Title     string         `json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	Columns   []*TableColumn `json:"columns"`
	Rows      []*TableRow    `json:"rows"`
}

type TableColumn struct {
	ID       int64  `json:"id"`
	TableID  int64  `json:"-"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type TableRow struct {
	ID       int64        `json:"id"`
	TableID  int64        `json:"-"`
	Position int          `json:"position"`
	Cells    []*TableCell `json:"cells"`
}

// TableCell содержит данные одной ячейки
type TableCell struct {
	ID       int64  `json:"id"`
	RowID    int64  `json:"-"`
	ColumnID int64  `json:"column_id"`
	Content  string `json:"content"`
}

// --- Структуры для API запросов ---

// CreateNoteTableRequest - запрос на создание новой таблицы
type CreateNoteTableRequest struct {
	Title   string   `json:"title" example:"Список задач"`
	Columns []string `json:"columns" example:"[\"Задача\",\"Срок\",\"Статус\"]"`
}

// AddTableRowRequest - запрос на добавление новой строки
type AddTableRowRequest struct {
	// Значения ячеек должны идти в том же порядке, что и колонки
	Cells []string `json:"cells" example:"[\"Реализовать API\",\"2024-12-31\",\"В процессе\"]"`
}
