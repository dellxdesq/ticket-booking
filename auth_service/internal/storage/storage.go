package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &Storage{DB: db}
	if err := store.InitDB(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Storage) InitDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		token TEXT
	);
	`
	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	return err
}
