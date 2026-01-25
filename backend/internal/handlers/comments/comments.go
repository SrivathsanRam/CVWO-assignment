package comments

import (
	"encoding/json"
	stdErrors "errors"
	"net/http"
	"strconv"

	"github.com/SrivathsanRam/CVWO_project/backend/internal/api"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/dataaccess"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/database"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

const (
	handleCreateOp = "comments.HandleCreate"
	handleListOp   = "comments.HandleList"
	handleGetOp    = "comments.HandleGet"
	handleUpdateOp = "comments.HandleUpdate"
	handleDeleteOp = "comments.HandleDelete"

	msgCreateCommentOK = "Successfully created comment"
	msgListCommentsOK  = "Successfully listed comments"
	msgGetCommentOK    = "Successfully retrieved comment"
	msgUpdateCommentOK = "Successfully updated comment"
	msgDeleteCommentOK = "Comment deleted successfully"

	errGetDB              = "failed to retrieve database"
	errDecodeCreate       = "failed to decode create comment request"
	errInvalidCreate      = "invalid create comment request"
	errCreateComment      = "failed to create comment"
	errListComments       = "failed to retrieve comments"
	errDecodeUpdate       = "failed to decode update comment request"
	errParseCommentID     = "failed to parse comment id"
	errFindComment        = "failed to retrieve comment"
	errUpdateComment      = "failed to update comment"
	errDeleteComment      = "failed to delete comment"
	errMarshalComment     = "failed to encode comment response"
	errMarshalComments    = "failed to encode comments response"
	errMarshalDelete      = "failed to encode delete response"
	errUserID             = "failed to retrieve user id"
	errUnauthorizedUpdate = "unauthorized to update comment"
	errUnauthorizedDelete = "unauthorized to delete comment"
	errCommentNotFound    = "comment not found"
)

// HandleCreate creates a new comment (authenticated).
func HandleCreate(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		return nil, errors.Wrapf(errors.New("missing user_id"), "%s: %s", handleCreateOp, errUserID)
	}

	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errDecodeCreate)
	}
	if req.PostID == 0 || req.Content == "" {
		return nil, errors.Wrapf(errors.New("post_id and content are required"), "%s: %s", handleCreateOp, errInvalidCreate)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errGetDB)
	}

	repo := dataaccess.NewCommentRepository(db)
	created, err := repo.Create(&models.Comment{
		Content: req.Content,
		PostID:  req.PostID,
		UserID:  userID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errCreateComment)
	}

	data, err := json.Marshal(created)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errMarshalComment)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgCreateCommentOK},
	}, nil
}

// HandleList returns all comments.
func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errGetDB)
	}

	repo := dataaccess.NewCommentRepository(db)
	comments, err := repo.ListAll()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errListComments)
	}

	data, err := json.Marshal(comments)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errMarshalComments)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgListCommentsOK},
	}, nil
}

// HandleGet returns a single comment by ID.
func HandleGet(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		idParam = r.URL.Query().Get("id")
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errParseCommentID)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errGetDB)
	}

	repo := dataaccess.NewCommentRepository(db)
	comment, err := repo.FindByID(id)
	if err != nil {
		if stdErrors.Is(err, dataaccess.ErrorCommentNotFound) {
			return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errCommentNotFound)
		}
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errFindComment)
	}

	data, err := json.Marshal(comment)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errMarshalComment)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgGetCommentOK},
	}, nil
}

// HandleUpdate updates a comment (authenticated + owner).
func HandleUpdate(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		return nil, errors.Wrapf(errors.New("missing user_id"), "%s: %s", handleUpdateOp, errUserID)
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		idParam = r.URL.Query().Get("id")
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errParseCommentID)
	}

	var req models.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errDecodeUpdate)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errGetDB)
	}

	repo := dataaccess.NewCommentRepository(db)
	updated, err := repo.Update(id, userID, req.Content)
	if err != nil {
		switch {
		case stdErrors.Is(err, dataaccess.ErrorCommentNotFound):
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errCommentNotFound)
		case stdErrors.Is(err, dataaccess.ErrorUnauthorized):
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errUnauthorizedUpdate)
		default:
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errUpdateComment)
		}
	}

	data, err := json.Marshal(updated)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errMarshalComment)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgUpdateCommentOK},
	}, nil
}

// HandleDelete deletes a comment (authenticated + owner).
func HandleDelete(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		return nil, errors.Wrapf(errors.New("missing user_id"), "%s: %s", handleDeleteOp, errUserID)
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		idParam = r.URL.Query().Get("id")
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errParseCommentID)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errGetDB)
	}

	repo := dataaccess.NewCommentRepository(db)
	if err := repo.Delete(id, userID); err != nil {
		switch {
		case stdErrors.Is(err, dataaccess.ErrorCommentNotFound):
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errCommentNotFound)
		case stdErrors.Is(err, dataaccess.ErrorUnauthorized):
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errUnauthorizedDelete)
		default:
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errDeleteComment)
		}
	}

	data, err := json.Marshal(map[string]string{"message": msgDeleteCommentOK})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errMarshalDelete)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgDeleteCommentOK},
	}, nil
}
