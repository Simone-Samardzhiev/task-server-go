package task

import (
	"github.com/google/uuid"
	"time"
)

// Task type struct hold information about a task.
type Task struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Type          int       `json:"type"`
	DueDate       time.Time `json:"due_date"`
	DateCompleted NullTime  `json:"date_completed"`
	DateDeleted   NullTime  `json:"date_deleted"`
}

// DataTask is used receive new task from http request.
type DataTask struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        int       `json:"type"`
	DueDate     time.Time `json:"due_date"`
}
