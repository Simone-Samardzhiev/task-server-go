package services

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"server/auth/tokens"
	"server/models"
	"server/repositories"
	"server/utils"
	"strconv"
)

// TaskService is the business login for tasks.
type TaskService interface {
	// GetTasks will return all of tasks
	GetTasks(ctx context.Context, token tokens.Token) ([]models.TaskPayload, *utils.ErrorResponse)

	// AddTask will add a new task and return the created one with an id.
	AddTask(ctx context.Context, token tokens.Token, taskPayload *models.NewTaskPayload) (*models.TaskPayload, *utils.ErrorResponse)
}

// DefaultTaskService is default implementation of [TaskService]
type DefaultTaskService struct {
	taskRepository repositories.TaskRepository
}

func (s *DefaultTaskService) GetTasks(ctx context.Context, token tokens.Token) ([]models.TaskPayload, *utils.ErrorResponse) {
	userId, err := strconv.Atoi(token.Subject)
	if err != nil {
		return nil, utils.InvalidToken()
	}

	tasks, err := s.taskRepository.GetTasks(ctx, userId)
	if err != nil {
		return nil, utils.InternalServerError()
	}

	return tasks, nil
}

func (s *DefaultTaskService) AddTask(ctx context.Context, token tokens.Token, taskPayload *models.NewTaskPayload) (*models.TaskPayload, *utils.ErrorResponse) {
	result, err := s.taskRepository.CheckPriority(ctx, taskPayload.Priority)
	if err != nil {
		return nil, utils.InternalServerError()
	}

	if !result {
		return nil, utils.NewErrorResponse("Invalid priority", http.StatusBadRequest)
	}

	task := models.TaskPayload{
		Id: uuid.New(),
		NewTaskPayload: models.NewTaskPayload{
			Name:        taskPayload.Name,
			Description: taskPayload.Description,
			Priority:    taskPayload.Priority,
			Date:        taskPayload.Date,
		},
	}

	userId, err := strconv.Atoi(token.Subject)
	if err != nil {
		return nil, utils.InvalidToken()
	}

	err = s.taskRepository.AddTask(ctx, &task, userId)
	if err != nil {
		return nil, utils.InternalServerError()
	}

	return &task, nil
}

func NewDefaultTaskService(taskRepository repositories.TaskRepository) *DefaultTaskService {
	return &DefaultTaskService{taskRepository}
}
