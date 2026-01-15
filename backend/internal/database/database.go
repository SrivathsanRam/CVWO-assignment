package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func GetDB() (*sql.DB, error) {
	if DB != nil {
		return DB, nil
	}

	connectionString := os.Getenv("DATABASE_URL") // get connection string from .env file
	if connectionString == "" {
		return nil, fmt.Errorf("env vars error")
	}
	log.Printf("connecting to database")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("database connection success")
	return DB, nil
}

func CloseDB() error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return fmt.Errorf("failed to close: %w", err)
		}
		log.Println("database connection closed")
	}
	return nil
}
