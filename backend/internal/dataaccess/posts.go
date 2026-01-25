package dataaccess

import (
	"database/sql"
	"errors"

	//"github.com/SrivathsanRam/CVWO_project/backend/internal/database"
	"time"

	"github.com/SrivathsanRam/CVWO_project/backend/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func scanPostRow(row *sql.Row) (*models.Post, error) {
	var post models.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.TopicID, &post.TopicName, &post.UserID, &post.UserName, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

const getPosts = `
SELECT p.id, p.title, p.content, p.topic_id, t.title, p.user_id, u.username, p.created_at, p.updated_at
FROM posts p
JOIN users u ON p.user_id = u.id
JOIN topics t ON p.topic_id = t.id
`

func (r *PostRepository) Create(t *models.Post) (*models.Post, error) {
	// Insert and return the inserted post with username using JOIN
	query := `
	WITH inserted AS (
	INSERT INTO posts (title, content, topic_id, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id
)
` + getPosts + `
WHERE p.id = (SELECT id FROM inserted)
`
	return scanPostRow(r.db.QueryRow(query, t.Title, t.Content, t.TopicID, t.UserID))
}

func (r *PostRepository) Update(id, userID int, title, content string) (*models.Post, error) {
	// Update and return the updated post with username using JOIN
	query := `
UPDATE posts
SET title = $1, content = $2, updated_at = $3
WHERE id = $4 AND user_id = $5
RETURNING id
`
	var updatedID int
	err := r.db.QueryRow(query, title, content, time.Now(), id, userID).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Could be "not found" OR "exists but not owner". Distinguish with one extra check.
			p, findErr := r.FindByID(id)
			if findErr != nil {
				return nil, ErrorPostNotFound
			}
			if p.UserID != userID {
				return nil, ErrorUnauthorized
			}
			return nil, ErrorPostNotFound // fallback
		}
		return nil, err
	}

	// Return full post (with username)
	return r.FindByID(updatedID)
}

func (r *PostRepository) Delete(id, userID int) error {
	result, err := r.db.Exec("DELETE FROM posts WHERE id=$1 AND user_id=$2", id, userID)
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
			return ErrorPostNotFound
		}
		if t.UserID != userID {
			return ErrorUnauthorized
		}
		return ErrorPostNotFound
	}
	return nil
}

func (r *PostRepository) FindByTitle(title string) (*models.Post, error) {
	row := r.db.QueryRow(getPosts+" WHERE p.title=$1", title)
	return scanPostRow(row)
}

func (r *PostRepository) FindByID(id int) (*models.Post, error) {
	row := r.db.QueryRow(getPosts+" WHERE p.id=$1", id)
	return scanPostRow(row)
}

func (r *PostRepository) ListAll() ([]models.Post, error) {
	rows, err := r.db.Query(getPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.TopicID, &post.TopicName, &post.UserID, &post.UserName, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

// ListByTopic returns all posts for a specific topic
func (r *PostRepository) ListByTopic(topicID int) ([]models.Post, error) {
	rows, err := r.db.Query(getPosts+" WHERE p.topic_id = $1 ORDER BY p.created_at DESC", topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.TopicID, &post.TopicName, &post.UserID, &post.UserName, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}
