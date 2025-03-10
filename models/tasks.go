package models

import (
	"github.com/google/uuid"
	"net/http"
	"server/utils"
	"time"
)

// NewTaskPayload stores tasks information.
type NewTaskPayload struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Date        time.Time `json:"date"`
}

func (t *NewTaskPayload) ValidatePayload() *utils.ErrorResponse {
	if t.Name == "" {
		return utils.NewErrorResponse("Name cannot be empty", http.StatusBadRequest)
	}

	if t.Description == "" {
		return utils.NewErrorResponse("Description cannot be empty", http.StatusBadRequest)
	}

	if t.Priority == "" {
		return utils.NewErrorResponse("Priority cannot be empty", http.StatusBadRequest)
	}

	return nil
}

// TaskPayload stores task information with an id created by the server.
type TaskPayload struct {
	Id uuid.UUID `json:"id"`
	NewTaskPayload
}

func (t *TaskPayload) ValidatePayload() *utils.ErrorResponse {
	if t.Name == "" {
		return utils.NewErrorResponse("Name cannot be empty", http.StatusBadRequest)
	}

	if t.Description == "" {
		return utils.NewErrorResponse("Description cannot be empty", http.StatusBadRequest)
	}

	if t.Priority == "" {
		return utils.NewErrorResponse("Priority cannot be empty", http.StatusBadRequest)
	}

	if t.Id == uuid.Nil {
		return utils.NewErrorResponse("Id cannot be empty", http.StatusBadRequest)
	}
	return nil
}
