package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(dataSourceName string) (*Storage, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}
