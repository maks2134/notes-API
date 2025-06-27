package repository

import (
	"database/sql"
	"notes-api/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.Password).Scan(&user.ID)
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, password FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
