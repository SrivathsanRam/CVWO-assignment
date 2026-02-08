package posts

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
	handleCreateOp = "posts.HandleCreate"
	handleListOp   = "posts.HandleList"
	handleGetOp    = "posts.HandleGet"
	handleUpdateOp = "posts.HandleUpdate"
	handleDeleteOp = "posts.HandleDelete"

	msgCreatePostOK = "Successfully created post"
	msgListPostsOK  = "Successfully listed posts"
	msgGetPostOK    = "Successfully retrieved post"
	msgUpdatePostOK = "Successfully updated post"
	msgDeletePostOK = "Post deleted successfully"

	errGetDB              = "failed to retrieve database"
	errDecodeCreate       = "failed to decode create post request"
	errInvalidCreate      = "invalid create post request"
	errCreatePost         = "failed to create post"
	errListPosts          = "failed to retrieve posts"
	errDecodeUpdate       = "failed to decode update post request"
	errParsePostID        = "failed to parse post id"
	errFindPost           = "failed to retrieve post"
	errUpdatePost         = "failed to update post"
	errDeletePost         = "failed to delete post"
	errMarshalPost        = "failed to encode post response"
	errMarshalPosts       = "failed to encode posts response"
	errMarshalDelete      = "failed to encode delete response"
	errUserID             = "failed to retrieve user id"
	errUnauthorizedUpdate = "unauthorized to update post"
	errUnauthorizedDelete = "unauthorized to delete post"
	errPostNotFound       = "post not found"
)

// HandleCreate creates a new post (authenticated).
func HandleCreate(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		return nil, errors.Wrapf(errors.New("missing user_id"), "%s: %s", handleCreateOp, errUserID)
	}

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errDecodeCreate)
	}
	if req.Title == "" || req.Content == "" || req.TopicID == 0 {
		return nil, errors.Wrapf(errors.New("title, content, and topic_id are required"), "%s: %s", handleCreateOp, errInvalidCreate)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errGetDB)
	}

	repo := dataaccess.NewPostRepository(db)
	created, err := repo.Create(&models.Post{
		Title:   req.Title,
		Content: req.Content,
		TopicID: req.TopicID,
		UserID:  userID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errCreatePost)
	}

	data, err := json.Marshal(created)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errMarshalPost)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgCreatePostOK},
	}, nil
}

// HandleList returns all posts.
func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errGetDB)
	}

	repo := dataaccess.NewPostRepository(db)

	var posts []models.Post
	if topicID := r.URL.Query().Get("topic_id"); topicID != "" {
		id, err := strconv.Atoi(topicID)
		if err != nil {
			return nil, errors.Wrapf(err, "%s: failed to parse topic_id", handleListOp)
		}
		posts, err = repo.ListByTopic(id)
	} else {
		posts, err = repo.ListAll()
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errListPosts)
	}

	data, err := json.Marshal(posts)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errMarshalPosts)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgListPostsOK},
	}, nil
}

// HandleGet returns a single post by ID.
func HandleGet(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		idParam = r.URL.Query().Get("id")
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errParsePostID)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errGetDB)
	}

	repo := dataaccess.NewPostRepository(db)
	post, err := repo.FindByID(id)
	if err != nil {
		if stdErrors.Is(err, dataaccess.ErrorPostNotFound) {
			return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errPostNotFound)
		}
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errFindPost)
	}

	data, err := json.Marshal(post)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errMarshalPost)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgGetPostOK},
	}, nil
}

// HandleUpdate updates a post (authenticated + owner).
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
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errParsePostID)
	}

	var req models.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errDecodeUpdate)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errGetDB)
	}

	repo := dataaccess.NewPostRepository(db)
	updated, err := repo.Update(id, userID, req.Title, req.Content)
	if err != nil {
		switch {
		case stdErrors.Is(err, dataaccess.ErrorPostNotFound):
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errPostNotFound)
		case stdErrors.Is(err, dataaccess.ErrorUnauthorized):
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errUnauthorizedUpdate)
		default:
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errUpdatePost)
		}
	}

	data, err := json.Marshal(updated)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errMarshalPost)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgUpdatePostOK},
	}, nil
}

// HandleDelete deletes a post (authenticated + owner).
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
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errParsePostID)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errGetDB)
	}

	repo := dataaccess.NewPostRepository(db)
	if err := repo.Delete(id, userID); err != nil {
		switch {
		case stdErrors.Is(err, dataaccess.ErrorPostNotFound):
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errPostNotFound)
		case stdErrors.Is(err, dataaccess.ErrorUnauthorized):
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errUnauthorizedDelete)
		default:
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errDeletePost)
		}
	}

	data, err := json.Marshal(map[string]string{"message": msgDeletePostOK})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errMarshalDelete)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgDeletePostOK},
	}, nil
}
