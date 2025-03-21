package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"server/auth/tokens"
	"server/models"
	"server/services"
	"server/utils"
)

// TaskHandler handles tasks request.
type TaskHandler interface {
	// GetTasks will return all tasks of a user.
	GetTasks() http.HandlerFunc

	// AddTask will add a new task.
	AddTask() http.HandlerFunc

	// UpdateTask will update an existing task.
	UpdateTask() http.HandlerFunc

	// DeleteTask will delete an existing task.
	DeleteTask() http.HandlerFunc
}

// DefaultTaskHandler is the default implementation of [TaskHandler]
type DefaultTaskHandler struct {
	taskService services.TaskService
}

func (h *DefaultTaskHandler) GetTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(tokens.JWTClaimsKey).(tokens.Token)
		if !ok {
			utils.HandleErrorResponse(w, utils.InternalServerError())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		tasks, errResponse := h.taskService.GetTasks(r.Context(), claims)
		if errResponse != nil {
			utils.HandleErrorResponse(w, errResponse)
			return
		}

		err := json.NewEncoder(w).Encode(tasks)
		if err != nil {
			utils.HandleErrorResponse(w, utils.InternalServerError())
			return
		}
	}
}

func (h *DefaultTaskHandler) AddTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(tokens.JWTClaimsKey).(tokens.Token)
		if !ok {
			utils.HandleErrorResponse(w, utils.InternalServerError())
			return
		}

		var taskPayload models.NewTaskPayload
		err := json.NewDecoder(r.Body).Decode(&taskPayload)
		if err != nil {
			utils.HandleErrorResponse(w, utils.InvalidJson())
		}

		if utils.HandlePayload(w, &taskPayload) {
			return
		}

		task, errResponse := h.taskService.AddTask(r.Context(), claims, &taskPayload)
		if errResponse != nil {
			utils.HandleErrorResponse(w, errResponse)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			utils.HandleErrorResponse(w, utils.InternalServerError())
		}
	}
}

func (h *DefaultTaskHandler) UpdateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskPayload models.TaskPayload
		err := json.NewDecoder(r.Body).Decode(&taskPayload)
		if err != nil {
			utils.HandleErrorResponse(w, utils.InvalidJson())
		}

		if utils.HandlePayload(w, &taskPayload) {
			return
		}

		errorResponse := h.taskService.UpdateTask(r.Context(), &taskPayload)
		if errorResponse != nil {
			utils.HandleErrorResponse(w, errorResponse)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *DefaultTaskHandler) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			utils.HandleErrorResponse(w, utils.NewErrorResponse("Invalid id", http.StatusBadRequest))
			return
		}

		errorResponse := h.taskService.DeleteTask(r.Context(), parsedId)
		if errorResponse != nil {
			utils.HandleErrorResponse(w, errorResponse)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func NewDefaultTaskHandler(taskService services.TaskService) *DefaultTaskHandler {
	return &DefaultTaskHandler{taskService}
}
