package handlers

import (
	"encoding/json"
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

func NewDefaultTaskHandler(taskService services.TaskService) *DefaultTaskHandler {
	return &DefaultTaskHandler{taskService}
}
