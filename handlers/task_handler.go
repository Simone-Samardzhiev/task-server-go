package handlers

import (
	"encoding/json"
	"net/http"
	"server/auth/tokens"
	"server/services"
	"server/utils"
)

// TaskHandler handles tasks request.
type TaskHandler interface {
	// GetTasks will return all tasks of a user.
	GetTasks() http.HandlerFunc
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

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
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

func NewDefaultTaskHandler(taskService services.TaskService) *DefaultTaskHandler {
	return &DefaultTaskHandler{taskService}
}
