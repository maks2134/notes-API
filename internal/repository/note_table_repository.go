package repository

import (
	"database/sql"
	"fmt"
	"notes-api/internal/model"
)

type NoteTableRepository interface {
	Create(tx *sql.Tx, table *model.NoteTable) error
	CreateColumns(tx *sql.Tx, tableID int64, columns []string) error
	AddRow(tableID int64, cells []string) (*model.TableRow, error)
	GetTablesByNoteID(noteID int64) ([]*model.NoteTable, error)
	BeginTx() (*sql.Tx, error)
}

type PostgresNoteTableRepository struct {
	db *sql.DB
}

func NewPostgresNoteTableRepository(db *sql.DB) NoteTableRepository {
	return &PostgresNoteTableRepository{db: db}
}

func (r *PostgresNoteTableRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *PostgresNoteTableRepository) Create(tx *sql.Tx, table *model.NoteTable) error {
	query := `INSERT INTO note_tables (note_id, title) VALUES ($1, $2) RETURNING id, created_at;`
	return tx.QueryRow(query, table.NoteID, table.Title).Scan(&table.ID, &table.CreatedAt)
}

func (r *PostgresNoteTableRepository) CreateColumns(tx *sql.Tx, tableID int64, columns []string) error {
	stmt, err := tx.Prepare(`INSERT INTO table_columns (table_id, name, "position") VALUES ($1, $2, $3);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, name := range columns {
		if _, err := stmt.Exec(tableID, name, i); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresNoteTableRepository) AddRow(tableID int64, cells []string) (*model.TableRow, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	colsQuery := `SELECT id FROM table_columns WHERE table_id = $1 ORDER BY "position" ASC;`
	rows, err := tx.Query(colsQuery, tableID)
	if err != nil {
		return nil, err
	}
	var columnIDs []int64
	for rows.Next() {
		var colID int64
		if err := rows.Scan(&colID); err != nil {
			return nil, err
		}
		columnIDs = append(columnIDs, colID)
	}
	rows.Close()

	if len(cells) != len(columnIDs) {
		return nil, fmt.Errorf("количество ячеек (%d) не соответствует количеству колонок (%d)", len(cells), len(columnIDs))
	}

	row := &model.TableRow{TableID: tableID}
	rowQuery := `INSERT INTO table_rows (table_id, "position") VALUES ($1, (SELECT COALESCE(MAX("position"), -1) + 1 FROM table_rows WHERE table_id = $1)) RETURNING id, "position";`
	err = tx.QueryRow(rowQuery, tableID).Scan(&row.ID, &row.Position)
	if err != nil {
		return nil, err
	}

	cellStmt, err := tx.Prepare(`INSERT INTO table_cells (row_id, column_id, content) VALUES ($1, $2, $3) RETURNING id, column_id, content;`)
	if err != nil {
		return nil, err
	}
	defer cellStmt.Close()

	for i, content := range cells {
		cell := &model.TableCell{}
		colID := columnIDs[i]
		err := cellStmt.QueryRow(row.ID, colID, content).Scan(&cell.ID, &cell.ColumnID, &cell.Content)
		if err != nil {
			return nil, err
		}
		row.Cells = append(row.Cells, cell)
	}

	return row, tx.Commit()
}

func (r *PostgresNoteTableRepository) GetTablesByNoteID(noteID int64) ([]*model.NoteTable, error) {
	tablesQuery := `SELECT id, title, created_at FROM note_tables WHERE note_id = $1;`
	rows, err := r.db.Query(tablesQuery, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*model.NoteTable
	tableMap := make(map[int64]*model.NoteTable)

	for rows.Next() {
		table := &model.NoteTable{NoteID: noteID}
		if err := rows.Scan(&table.ID, &table.Title, &table.CreatedAt); err != nil {
			return nil, err
		}
		tables = append(tables, table)
		tableMap[table.ID] = table
	}

	if len(tables) == 0 {
		return nil, nil
	}

	for _, table := range tables {
		cols, err := r.db.Query(`SELECT id, name, "position" FROM table_columns WHERE table_id = $1 ORDER BY "position" ASC;`, table.ID)
		if err != nil {
			return nil, err
		}
		defer cols.Close()

		for cols.Next() {
			col := &model.TableColumn{TableID: table.ID}
			if err := cols.Scan(&col.ID, &col.Name, &col.Position); err != nil {
				return nil, err
			}
			table.Columns = append(table.Columns, col)
		}

		rowsData, err := r.db.Query(`SELECT id, "position" FROM table_rows WHERE table_id = $1 ORDER BY "position" ASC;`, table.ID)
		if err != nil {
			return nil, err
		}
		defer rowsData.Close()

		for rowsData.Next() {
			row := &model.TableRow{TableID: table.ID}
			if err := rowsData.Scan(&row.ID, &row.Position); err != nil {
				return nil, err
			}

			cellsData, err := r.db.Query(`SELECT id, column_id, content FROM table_cells WHERE row_id = $1;`, row.ID)
			if err != nil {
				return nil, err
			}
			defer cellsData.Close()

			for cellsData.Next() {
				cell := &model.TableCell{RowID: row.ID}
				if err := cellsData.Scan(&cell.ID, &cell.ColumnID, &cell.Content); err != nil {
					return nil, err
				}
				row.Cells = append(row.Cells, cell)
			}
			table.Rows = append(table.Rows, row)
		}
	}

	return tables, nil
}
