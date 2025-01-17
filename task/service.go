package task

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"server/auth"
	"server/utils"
	"time"
)

// Service type interface will manage business.
type Service interface {
	// GetTasks used to return a slice of all tasks that belongs to a user.
	GetTasks(token *auth.CustomClaims) ([]Task, error)

	// AddTask used to add a new task and link it to a user.
	AddTask(task *DataTask, token *auth.CustomClaims) (*Task, error)

	// UpdateTask used to update task data.
	UpdateTask(task *Task) error

	// DeleteTask used to delete a task.
	DeleteTask(id *uuid.UUID) error
}

type DefaultService struct {
	repository Repository
}

func NewDefaultService(repository Repository) *DefaultService {
	return &DefaultService{repository}
}

// ParseSubjectId will return the id of the subject of the token.
func (d *DefaultService) ParseSubjectId(token *auth.CustomClaims) (*uuid.UUID, error) {
	id, err := uuid.Parse(token.Subject)
	return &id, err
}

func (d *DefaultService) GetTasks(token *auth.CustomClaims) ([]Task, error) {
	id, err := d.ParseSubjectId(token)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	tasks, err := d.repository.GetTasksByUserId(id)
	if err != nil {
		return nil, utils.InternalServerErr
	}

	return tasks, nil
}

func (d *DefaultService) AddTask(task *DataTask, token *auth.CustomClaims) (*Task, error) {
	id, err := d.ParseSubjectId(token)

	newTask := &Task{
		Id:          uuid.New(),
		Name:        task.Name,
		Description: task.Description,
		Type:        task.Type,
		DueDate:     task.DueDate,
		DateDeleted: NullTime{
			sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
		DateCompleted: NullTime{
			sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
	}

	err = d.repository.AddTask(newTask, id)
	if err != nil {
		fmt.Println(err)
		return nil, utils.InternalServerErr
	}
	return newTask, nil
}

func (d *DefaultService) UpdateTask(task *Task) error {
	updated, err := d.repository.UpdateTask(task)
	if err != nil {
		return utils.InternalServerErr
	}

	if !updated {
		return utils.NotFoundErr
	}

	return nil
}

func (d *DefaultService) DeleteTask(id *uuid.UUID) error {
	deleted, err := d.repository.DeleteTask(id)
	if err != nil {
		return utils.InternalServerErr
	}

	if !deleted {
		return utils.NotFoundErr
	}

	return nil
}
