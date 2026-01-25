package dataaccess

import (
	"database/sql"
	"errors"

	//"github.com/SrivathsanRam/CVWO_project/backend/internal/database"
	"time"

	"github.com/SrivathsanRam/CVWO_project/backend/internal/models"
)

type TopicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

func scanTopicRow(row *sql.Row) (*models.Topic, error) {
	var topic models.Topic
	err := row.Scan(&topic.ID, &topic.Title, &topic.Description, &topic.UserID, &topic.UserName, &topic.CreatedAt, &topic.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorTopicNotFound
		}
		return nil, err
	}
	return &topic, nil
}

const getTopics = `
SELECT t.id, t.title, t.description, t.user_id, u.username, t.created_at, t.updated_at
FROM topics t
JOIN users u ON t.user_id = u.id
`

func (r *TopicRepository) Create(t *models.Topic) (*models.Topic, error) {
	// Insert and return the inserted topic with username using JOIN
	query := `
	WITH inserted AS (
	INSERT INTO topics (title, description, user_id)
	VALUES ($1, $2, $3)
	RETURNING id
)
` + getTopics + `
WHERE t.id = (SELECT id FROM inserted)
`
	return scanTopicRow(r.db.QueryRow(query, t.Title, t.Description, t.UserID))
}

func (r *TopicRepository) Update(id, userID int, title, description string) (*models.Topic, error) {
	// Update and return the updated topic with username using JOIN
	query := `
UPDATE topics
SET title = $1, description = $2, updated_at = $3
WHERE id = $4 AND user_id = $5
RETURNING id
`
	var updatedID int
	err := r.db.QueryRow(query, title, description, time.Now(), id, userID).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Could be "not found" OR "exists but not owner". Distinguish with one extra check.
			t, findErr := r.FindByID(id)
			if findErr != nil {
				return nil, ErrorTopicNotFound
			}
			if t.UserID != userID {
				return nil, ErrorUnauthorized
			}
			return nil, ErrorTopicNotFound // fallback
		}
		return nil, err
	}

	// Return full topic (with username)
	return r.FindByID(updatedID)
}

func (r *TopicRepository) Delete(id, userID int) error {
	result, err := r.db.Exec("DELETE FROM topics WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Could be not found or exists but user is not the owner
		t, findErr := r.FindByID(id)
		if findErr != nil {
			return ErrorTopicNotFound
		}
		if t.UserID != userID {
			return ErrorUnauthorized
		}
		return ErrorTopicNotFound
	}
	return nil
}

func (r *TopicRepository) FindByTitle(title string) (*models.Topic, error) {
	row := r.db.QueryRow(getTopics+" WHERE t.title=$1", title)
	return scanTopicRow(row)
}

func (r *TopicRepository) FindByID(id int) (*models.Topic, error) {
	row := r.db.QueryRow(getTopics+" WHERE t.id=$1", id)
	return scanTopicRow(row)
}

func (r *TopicRepository) ListAll() ([]models.Topic, error) {
	rows, err := r.db.Query(getTopics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var topic models.Topic
		if err := rows.Scan(&topic.ID, &topic.Title, &topic.Description, &topic.UserID, &topic.UserName, &topic.CreatedAt, &topic.UpdatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, rows.Err()
}
