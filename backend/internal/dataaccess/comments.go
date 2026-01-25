package dataaccess

import (
	"database/sql"
	"errors"

	//"github.com/SrivathsanRam/CVWO_project/backend/internal/database"
	"time"

	"github.com/SrivathsanRam/CVWO_project/backend/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func scanCommentRow(row *sql.Row) (*models.Comment, error) {
	var comment models.Comment
	err := row.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.UserName, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorCommentNotFound
		}
		return nil, err
	}
	return &comment, nil
}

const getComments = `
SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.created_at, c.updated_at
FROM comments c
JOIN users u ON c.user_id = u.id
`

func (r *CommentRepository) Create(t *models.Comment) (*models.Comment, error) {
	// Insert and return the inserted comment with username using JOIN
	query := `
	WITH inserted AS (
	INSERT INTO comments (content, post_id, user_id)
	VALUES ($1, $2, $3)
	RETURNING id
)
` + getComments + `
WHERE c.id = (SELECT id FROM inserted)
`
	return scanCommentRow(r.db.QueryRow(query, t.Content, t.PostID, t.UserID))
}

func (r *CommentRepository) Update(id, userID int, content string) (*models.Comment, error) {
	// Update and return the updated comment with username using JOIN
	query := `
UPDATE comments
SET content = $1, updated_at = $2
WHERE id = $3 AND user_id = $4
RETURNING id
`
	var updatedID int
	err := r.db.QueryRow(query, content, time.Now(), id, userID).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Could be "not found" OR "exists but not owner". Distinguish with one extra check.
			t, findErr := r.FindByID(id)
			if findErr != nil {
				return nil, ErrorCommentNotFound
			}
			if t.UserID != userID {
				return nil, ErrorUnauthorized
			}
			return nil, ErrorCommentNotFound // fallback
		}
		return nil, err
	}

	// Return full comment (with username)
	return r.FindByID(updatedID)
}

func (r *CommentRepository) Delete(id, userID int) error {
	result, err := r.db.Exec("DELETE FROM comments WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Could be not found or exists but user is not the owner
		c, findErr := r.FindByID(id)
		if findErr != nil {
			return ErrorCommentNotFound
		}
		if c.UserID != userID {
			return ErrorUnauthorized
		}
		return ErrorCommentNotFound
	}
	return nil
}

func (r *CommentRepository) FindByID(id int) (*models.Comment, error) {
	row := r.db.QueryRow(getComments+" WHERE id=$1", id)
	return scanCommentRow(row)
}

func (r *CommentRepository) ListAll() ([]models.Comment, error) {
	rows, err := r.db.Query(getComments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.UserName, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

// ListByPost returns all comments for a specific post
func (r *CommentRepository) ListByPost(postID int) ([]models.Comment, error) {
	rows, err := r.db.Query(getComments+" WHERE c.post_id = $1 ORDER BY c.created_at ASC", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.UserName, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}
