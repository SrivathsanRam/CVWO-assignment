package users

import (
	"encoding/json"
	"net/http"

	"github.com/SrivathsanRam/CVWO_project/backend/internal/api"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/dataaccess"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/database"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/models"
	"github.com/pkg/errors"
)

const (
	handleListOp           = "users.HandleList"
	handleLoginOp          = "users.HandleLogin"
	handleLogoutOp         = "users.HandleLogout"
	handleGetCurrentUserOp = "users.HandleGetCurrentUser"

	msgListUsersOK       = "Successfully listed users"
	msgLoginOK           = "Login successful"
	msgLogoutOK          = "Logout successful"
	msgGetCurrentUserOK  = "Successfully retrieved user"

	errGetDB              = "failed to retrieve database"
	errList               = "failed to retrieve users"
	errMarshalList        = "failed to encode users response"
	errDecodeLogin        = "failed to decode login request"
	errInvalidLogin       = "invalid login request"
	errLogin              = "failed to login"
	errMarshalLogin       = "failed to encode login response"
	errMarshalLogout      = "failed to encode logout response"
	errUserID             = "failed to retrieve user id"
	errFindUser           = "failed to retrieve user"
	errMarshalCurrentUser = "failed to encode user response"
)

// HandleList handles the listing of all users.
func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errGetDB)
	}

	repo := dataaccess.NewUserRepository(db)
	users, err := repo.ListAll()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errList)
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleListOp, errMarshalList)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgListUsersOK},
	}, nil
}

// Login handler
func HandleLogin(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleLoginOp, errGetDB)
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleLoginOp, errDecodeLogin)
	}
	
	if req.UserName == "" {
		return nil, errors.Wrapf(errors.New("username is required"), "%s: %s", handleLoginOp, errInvalidLogin)
	}

	repo := dataaccess.NewUserRepository(db)
	user, err := repo.Create(req.UserName)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleLoginOp, errLogin)
	}

	data, err := json.Marshal(models.LoginResponse{
		User:    *user,
		Message: msgLoginOK,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleLoginOp, errMarshalLogin)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgLoginOK},
	}, nil
}

// Logout handler
func HandleLogout(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	data, err := json.Marshal(map[string]string{"message": msgLogoutOK})
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleLogoutOp, errMarshalLogout)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgLogoutOK},
	}, nil
}

//getting current user from user_id in context
func HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		return nil, errors.Wrapf(errors.New("missing user_id"), "%s: %s", handleGetCurrentUserOp, errUserID)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetCurrentUserOp, errGetDB)
	}

	repo := dataaccess.NewUserRepository(db)
	user, err := repo.FindByID(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetCurrentUserOp, errFindUser)
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: %s", handleGetCurrentUserOp, errMarshalCurrentUser)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{msgGetCurrentUserOK},
	}, nil
}
