package dataaccess

import (
	"database/sql"
	"errors"
	//"github.com/SrivathsanRam/CVWO_project/backend/internal/database"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/models"
)


type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func scanUserRow(row *sql.Row) (*models.User, error) {
	var user models.User
	err := row.Scan(&user.ID, &user.UserName, &user.CreatedAt); if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

const getUsers = `SELECT id, username, created_at FROM users`

func (r *UserRepository) Create(username string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := r.FindByUsername(username)
	if err == nil {
		return existingUser, nil // User exists, return it
	}
	if !errors.Is(err, ErrorUserNotFound) {
		return nil, err // Some other error occurred
	}

	
	// Insert new user
	var id int
	err = r.db.QueryRow(
		`INSERT INTO users (username) VALUES ($1) RETURNING id`, username,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	row := r.db.QueryRow(getUsers+" WHERE username=$1", username)
	return scanUserRow(row)
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	row := r.db.QueryRow(getUsers+" WHERE id=$1", id)
	return scanUserRow(row)
}

func (r *UserRepository) ListAll() ([]models.User, error) {
	rows, err := r.db.Query(getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.UserName, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}	