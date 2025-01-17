package task

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"server/auth"
	"server/utils"
)

// Handler type interface used to respond to http request for tasks.
type Handler interface {
	// GetTasks will respond with all tasks.
	GetTasks() http.Handler

	// AddTask will add a new task.
	AddTask() http.Handler

	// UpdateTask will update an existing task.
	UpdateTask() http.Handler

	// DeleteTask will delete an existing task.
	DeleteTask() http.Handler
}

type DefaultHandler struct {
	service Service
}

func NewDefaultHandler(service Service) *DefaultHandler {
	return &DefaultHandler{service: service}
}

func NewHandler(service Service) Handler {
	return &DefaultHandler{service: service}
}

func (d *DefaultHandler) GetTasks() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value(auth.TokenKey).(*auth.CustomClaims)
		if !ok {
			http.Error(w, utils.UnauthorizedErrorMessage, http.StatusUnauthorized)
			return
		}

		tasks, err := d.service.GetTasks(token)
		if err != nil {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tasks)
		if err != nil {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (d *DefaultHandler) AddTask() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value(auth.TokenKey).(*auth.CustomClaims)
		if !ok {
			http.Error(w, utils.UnauthorizedErrorMessage, http.StatusUnauthorized)
			return
		}

		var task DataTask
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, utils.BadRequestErrorMessage, http.StatusBadRequest)
			return
		}

		createdTask, err := d.service.AddTask(&task, token)
		if err != nil {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(createdTask)
		if err != nil {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (d *DefaultHandler) UpdateTask() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(auth.TokenKey).(*auth.CustomClaims)
		if !ok {
			http.Error(w, utils.UnauthorizedErrorMessage, http.StatusUnauthorized)
		}

		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, utils.BadRequestErrorMessage, http.StatusBadRequest)
		}

		err = d.service.UpdateTask(&task)
		if errors.Is(err, utils.NotFoundErr) {
			http.Error(w, utils.NotFoundErrorMessage, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (d *DefaultHandler) DeleteTask() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(auth.TokenKey).(*auth.CustomClaims)
		if !ok {
			http.Error(w, utils.UnauthorizedErrorMessage, http.StatusUnauthorized)
		}

		id, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, utils.BadRequestErrorMessage, http.StatusBadRequest)
		}

		err = d.service.DeleteTask(&id)
		if errors.Is(err, utils.NotFoundErr) {
			http.Error(w, utils.NotFoundErrorMessage, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, utils.InternalServerErrorMessage, http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})
}
