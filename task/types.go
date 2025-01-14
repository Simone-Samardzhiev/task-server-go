package task

import (
	"github.com/google/uuid"
	"time"
)

// Task type struct hold information about a task.
type Task struct {
	Id            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Type          int        `json:"type"`
	DueDate       time.Time  `json:"due_date"`
	DateCompleted *time.Time `json:"date_completed"`
	DateDeleted   *time.Time `json:"date_deleted"`
	UserId        uuid.UUID
}

// NewTask will create a [Task]
func NewTask(id uuid.UUID, Type int, description string, name string, dueDate time.Time, dateDeleted *time.Time, dateCompleted *time.Time) *Task {
	return &Task{
		Id:            id,
		Type:          Type,
		Description:   description,
		Name:          name,
		DueDate:       dueDate,
		DateDeleted:   dateDeleted,
		DateCompleted: dateCompleted,
	}
}

// DataTask is used receive new task from http request.
type DataTask struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        int       `json:"type"`
	DueDate     time.Time `json:"due_date"`
}
