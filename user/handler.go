package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/auth"
	"server/utils"
)

// Handler type interface will respond to http request to user.
type Handler interface {
	Register() http.Handler
	Login() http.Handler
	Refresh() http.Handler
}

// DefaultHandler type struct is default implementation of [Handler].
type DefaultHandler struct {
	service Service
}

func (d *DefaultHandler) GetJsonUserFromBody(r *http.Request) (*JsonUser, error) {
	decoder := json.NewDecoder(r.Body)
	var user JsonUser
	err := decoder.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *DefaultHandler) Register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user from the body.
		user, err := d.GetJsonUserFromBody(r)
		if err != nil {
			http.Error(w, utils.BadRequestErrorMessage, http.StatusBadRequest)
			return
		}

		// Register the user and handles errors.
		err = d.service.Register(user)
		if errors.Is(err, utils.ConflictErr) {
			http.Error(w, utils.ConflictErrorMessage, http.StatusConflict)
			return
		} else if errors.Is(err, utils.InternalServerErr) {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func (d *DefaultHandler) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user from the body.
		user, err := d.GetJsonUserFromBody(r)
		if err != nil {
			http.Error(w, utils.BadRequestErrorMessage, http.StatusBadRequest)
			return
		}

		// Login the user and handles the errors.
		group, err := d.service.Login(user)
		if errors.Is(err, utils.UnauthorizedErr) {
			http.Error(w, utils.UnauthorizedErrorMessage, http.StatusUnauthorized)
			return
		} else if errors.Is(err, utils.InternalServerErr) {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		// Encode the tokens data and send it.
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(group)
		if err != nil {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
		}
	})
}

func (d *DefaultHandler) Refresh() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from context.
		token := r.Context().Value(auth.TokenKey).(*auth.CustomClaims)

		// Refresh the tokens and handler errors.
		group, err := d.service.RefreshToken(token)
		if errors.Is(err, utils.UnauthorizedErr) {
			http.Error(w, utils.UnauthorizedErrorMessage, http.StatusUnauthorized)
			return
		} else if errors.Is(err, utils.InternalServerErr) {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		// Encode the tokens data and send it.
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(group)
		if err != nil {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
		}
	})
}
