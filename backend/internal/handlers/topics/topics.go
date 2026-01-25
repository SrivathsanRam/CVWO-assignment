package topics

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
	handleCreateOp = "topics.HandleCreate"
	handleListOp   = "topics.HandleList"
	handleGetOp    = "topics.HandleGet"
	handleUpdateOp = "topics.HandleUpdate"
	handleDeleteOp = "topics.HandleDelete"

	msgCreateTopicOK = "Successfully created topic"
	msgListTopicsOK  = "Successfully listed topics"
	msgGetTopicOK    = "Successfully retrieved topic"
	msgUpdateTopicOK = "Successfully updated topic"
	msgDeleteTopicOK = "Topic deleted successfully"

	errGetDB              = "failed to retrieve database"
	errDecodeCreate       = "failed to decode create topic request"
	errInvalidCreate      = "invalid create topic request"
	errCreateTopic        = "failed to create topic"
	errListTopics         = "failed to retrieve topics"
	errDecodeUpdate       = "failed to decode update topic request"
	errParseTopicID       = "failed to parse topic id"
	errFindTopic          = "failed to retrieve topic"
	errUpdateTopic        = "failed to update topic"
	errDeleteTopic        = "failed to delete topic"
	errMarshalTopic       = "failed to encode topic response"
	errMarshalTopics      = "failed to encode topics response"
	errMarshalDelete      = "failed to encode delete response"
	errUserID             = "failed to retrieve user id"
	errUnauthorizedUpdate = "unauthorized to update topic"
	errUnauthorizedDelete = "unauthorized to delete topic"
	errTopicNotFound      = "topic not found"
)

// HandleCreate creates a new topic (authenticated).
func HandleCreate(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		return nil, errors.Wrapf(errors.New("missing user_id"), "%s: %s", handleCreateOp, errUserID)
	}

	var req models.CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errDecodeCreate)
	}
	if req.Title == "" || req.Description == "" {
		return nil, errors.Wrapf(errors.New("title and description are required"), "%s: %s", handleCreateOp, errInvalidCreate)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errGetDB)
	}

	repo := dataaccess.NewTopicRepository(db)
	created, err := repo.Create(&models.Topic{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errCreateTopic)
	}

	data, err := json.Marshal(created)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleCreateOp, errMarshalTopic)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgCreateTopicOK},
	}, nil
}

// HandleList returns all topics.
func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errGetDB)
	}

	repo := dataaccess.NewTopicRepository(db)
	topics, err := repo.ListAll()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errListTopics)
	}

	data, err := json.Marshal(topics)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errMarshalTopics)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgListTopicsOK},
	}, nil
}

// HandleGet returns a single topic by ID.
func HandleGet(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		idParam = r.URL.Query().Get("id")
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errParseTopicID)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errGetDB)
	}

	repo := dataaccess.NewTopicRepository(db)
	topic, err := repo.FindByID(id)
	if err != nil {
		if stdErrors.Is(err, dataaccess.ErrorTopicNotFound) {
			return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errTopicNotFound)
		}
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errFindTopic)
	}

	data, err := json.Marshal(topic)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetOp, errMarshalTopic)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgGetTopicOK},
	}, nil
}

// HandleUpdate updates a topic (authenticated + owner).
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
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errParseTopicID)
	}

	var req models.UpdateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errDecodeUpdate)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errGetDB)
	}

	repo := dataaccess.NewTopicRepository(db)
	updated, err := repo.Update(id, userID, req.Title, req.Description)
	if err != nil {
		switch {
		case stdErrors.Is(err, dataaccess.ErrorTopicNotFound):
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errTopicNotFound)
		case stdErrors.Is(err, dataaccess.ErrorUnauthorized):
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errUnauthorizedUpdate)
		default:
			return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errUpdateTopic)
		}
	}

	data, err := json.Marshal(updated)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleUpdateOp, errMarshalTopic)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgUpdateTopicOK},
	}, nil
}

// HandleDelete deletes a topic (authenticated + owner).
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
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errParseTopicID)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errGetDB)
	}

	repo := dataaccess.NewTopicRepository(db)
	if err := repo.Delete(id, userID); err != nil {
		switch {
		case stdErrors.Is(err, dataaccess.ErrorTopicNotFound):
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errTopicNotFound)
		case stdErrors.Is(err, dataaccess.ErrorUnauthorized):
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errUnauthorizedDelete)
		default:
			return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errDeleteTopic)
		}
	}

	data, err := json.Marshal(map[string]string{"message": msgDeleteTopicOK})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleDeleteOp, errMarshalDelete)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgDeleteTopicOK},
	}, nil
}
