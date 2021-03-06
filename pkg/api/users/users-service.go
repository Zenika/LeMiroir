package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/allez-chauffe/marcel/pkg/api/auth"
	"github.com/allez-chauffe/marcel/pkg/api/commons"
	"github.com/allez-chauffe/marcel/pkg/db"
	"github.com/allez-chauffe/marcel/pkg/db/users"
)

type Service struct{}

func NewService() *Service {
	return new(Service)
}

type UserPayload struct {
	*users.User
	Password string `json:"password"`
}

var UserNotFoundErr = errors.New("User not found")

func (s *Service) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckPermissions(r, nil, "admin") {
		commons.WriteResponse(w, http.StatusForbidden, "")
		return
	}

	payload := s.getUserPayload(w, r)

	if payload.Login == "" || payload.DisplayName == "" || payload.Password == "" {
		commons.WriteResponse(w, http.StatusBadRequest, "Malformed request, missing required fields")
		return
	}

	u := payload.User

	if err := u.SetPassword(payload.Password); err != nil {
		commons.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Users().Insert(u); err != nil {
		commons.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	commons.WriteJsonResponse(w, u)
}

func (s Service) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	if !auth.CheckPermissions(r, []string{userID}, "admin") {
		commons.WriteResponse(w, http.StatusForbidden, "")
		return
	}

	payload := s.getUserPayload(w, r)

	db.Transactional(func(tx *db.Tx) error {
		savedUser, err := tx.Users().Get(userID)
		if err != nil {
			commons.WriteResponse(w, http.StatusInternalServerError, err.Error())
			return err
		}

		if savedUser == nil || savedUser.ID != payload.ID {
			commons.WriteResponse(w, http.StatusNotFound, "")
			return UserNotFoundErr
		}

		if payload.Password != "" {
			unchanged, err := savedUser.CheckPassword(payload.Password)
			if err != nil {
				commons.WriteResponse(w, http.StatusInternalServerError, err.Error())
				return err
			}

			if !unchanged {
				if err := savedUser.SetPassword(payload.Password); err != nil {
					commons.WriteResponse(w, http.StatusInternalServerError, err.Error())
					return err
				}

				savedUser.LastDisconnection = time.Now()
			}
		}

		savedUser.DisplayName = payload.DisplayName
		savedUser.Login = payload.Login

		if auth.CheckPermissions(r, nil, "admin") {
			savedUser.Role = payload.Role
		}

		if err := tx.Users().Update(savedUser); err != nil {
			commons.WriteResponse(w, http.StatusInternalServerError, err.Error())
			return err
		}

		return commons.WriteJsonResponse(w, savedUser)
	})
}

func (s *Service) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if !auth.CheckPermissions(r, nil, "admin") {
		commons.WriteResponse(w, http.StatusForbidden, "")
		return
	}

	users, err := db.Users().List()
	if err != nil {
		commons.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	commons.WriteJsonResponse(w, users)
}

func (s *Service) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	if !auth.CheckPermissions(r, []string{userID}, "admin") {
		commons.WriteResponse(w, http.StatusForbidden, "")
		return
	}

	err := db.Users().Delete(userID)
	if err != nil {
		commons.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	commons.WriteResponse(w, http.StatusNoContent, "")
}

func (s *Service) getUserPayload(w http.ResponseWriter, r *http.Request) *UserPayload {
	user := &UserPayload{
		User: users.New(),
	}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		commons.WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("Error while parsing JSON (%s)", err.Error()))
		return nil
	}

	return user
}
