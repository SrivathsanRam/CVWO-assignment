package database

import (
	"database/sql"
	"os"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)


var DB *sql.DB;

func GetDB() (*sql.DB, error) {
	if DB != nil {
		return DB, nil
	}

	connectionString := os.Getenv("DATABASE_URL") // get connection string from .env file
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	DB = db
	log.Println("Database connection established")
	return DB, nil
}

func CloseDB() error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		log.Println("Database connection closed")
	}
	return nil
}
