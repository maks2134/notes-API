-- =================================================================
-- SQL-скрипт для инициализации базы данных "notes_api_db"
-- Версия: 1.3 (с поддержкой вложенных пользовательских таблиц, без админов)
-- =================================================================

-- Удаление существующих объектов в обратном порядке зависимостей
DROP TABLE IF EXISTS table_cells;
DROP TABLE IF EXISTS table_rows;
DROP TABLE IF EXISTS table_columns;
DROP TABLE IF EXISTS note_tables;
DROP TABLE IF EXISTS checklist_items;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS text_style;

CREATE TYPE text_style AS ENUM (
    'normal',
    'bold',
    'italic'
    );

CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       username VARCHAR(255) NOT NULL UNIQUE,
                       password TEXT NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE notes (
                       id BIGSERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       content TEXT,
                       style text_style NOT NULL DEFAULT 'normal',
                       user_id BIGINT NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE checklist_items (
                                 id BIGSERIAL PRIMARY KEY,
                                 text TEXT NOT NULL,
                                 completed BOOLEAN NOT NULL DEFAULT FALSE,
                                 style text_style NOT NULL DEFAULT 'normal',
                                 note_id BIGINT NOT NULL,
                                 created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                 CONSTRAINT fk_note FOREIGN KEY(note_id) REFERENCES notes(id) ON DELETE CASCADE
);

CREATE TABLE note_tables (
                             id BIGSERIAL PRIMARY KEY,
                             note_id BIGINT NOT NULL,
                             title VARCHAR(255) NOT NULL,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                             CONSTRAINT fk_note_table FOREIGN KEY(note_id) REFERENCES notes(id) ON DELETE CASCADE
);

CREATE TABLE table_columns (
                               id BIGSERIAL PRIMARY KEY,
                               table_id BIGINT NOT NULL,
                               name VARCHAR(255) NOT NULL,
                               position INT NOT NULL,
                               CONSTRAINT fk_table_column FOREIGN KEY(table_id) REFERENCES note_tables(id) ON DELETE CASCADE,
                               UNIQUE (table_id, name),
                               UNIQUE (table_id, position)
);

CREATE TABLE table_rows (
                            id BIGSERIAL PRIMARY KEY,
                            table_id BIGINT NOT NULL,
                            position INT NOT NULL,
                            CONSTRAINT fk_table_row FOREIGN KEY(table_id) REFERENCES note_tables(id) ON DELETE CASCADE
);

CREATE TABLE table_cells (
                             id BIGSERIAL PRIMARY KEY,
                             row_id BIGINT NOT NULL,
                             column_id BIGINT NOT NULL,
                             content TEXT NOT NULL,
                             CONSTRAINT fk_cell_row FOREIGN KEY(row_id) REFERENCES table_rows(id) ON DELETE CASCADE,
                             CONSTRAINT fk_cell_column FOREIGN KEY(column_id) REFERENCES table_columns(id) ON DELETE CASCADE,
                             UNIQUE (row_id, column_id)
);