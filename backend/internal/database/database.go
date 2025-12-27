package database

import (
	"context"
	"database/sql"
	"time"
)

type Database struct {
	db *sql.DB
}


func GetDB() (*Database, error) {
	db, err := sql.Open("pgx", "your_connection_string")
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}
