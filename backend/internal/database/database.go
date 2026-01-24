package database

import (
	"database/sql"
	"os"
	//_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

var instance *Database

func GetDB() (*Database, error) {
	if instance != nil {
		return instance, nil
	}

	connectionString := os.Getenv("DATABASE_URL") // get connection string from .env file
	db, err := sql.Open("postgres", connectionString)

	// Handling potential connection errors
	if err != nil {
		return nil, err
	}
	// Verifying the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// Initialising database instance
	instance = &Database{db: db}
	return instance, nil

}
