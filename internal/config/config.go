package config

import (
	"fmt"
)

type Config struct {
	Port string
	DB   DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (db *DBConfig) GetPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.DBName,
	)
}

func LoadConfig() *Config {
	return &Config{
		Port: "8080",
		DB: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "notes",
			DBName:   "notes_api_db",
			SSLMode:  "disable",
		},
	}
}
