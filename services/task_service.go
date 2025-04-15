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

	// UpdateTask will update an existing task.
	UpdateTask(ctx context.Context, taskPayload *models.TaskPayload) *utils.ErrorResponse

	// DeleteTask will delete an existing task.
	DeleteTask(ctx context.Context, taskId uuid.UUID) *utils.ErrorResponse
}

// DefaultTaskService is default implementation of [TaskService]
type DefaultTaskService struct {
	taskRepository repositories.TaskRepository
}

func (s *DefaultTaskService) GetTasks(ctx context.Context, token tokens.Token) ([]models.TaskPayload, *utils.ErrorResponse) {
	userId, err := strconv.Atoi(token.Subject)
	if err != nil {
		return nil, utils.InvalidTokenErrorResponse()
	}

	tasks, err := s.taskRepository.GetTasks(ctx, userId)
	if err != nil {
		return nil, utils.InternalServerErrorResponse()
	}

	return tasks, nil
}

func (s *DefaultTaskService) AddTask(ctx context.Context, token tokens.Token, taskPayload *models.NewTaskPayload) (*models.TaskPayload, *utils.ErrorResponse) {
	result, err := s.taskRepository.CheckPriority(ctx, taskPayload.Priority)
	if err != nil {
		return nil, utils.InternalServerErrorResponse()
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
		return nil, utils.InvalidTokenErrorResponse()
	}

	err = s.taskRepository.AddTask(ctx, &task, userId)
	if err != nil {
		return nil, utils.InternalServerErrorResponse()
	}

	return &task, nil
}

func (s *DefaultTaskService) UpdateTask(ctx context.Context, taskPayload *models.TaskPayload) *utils.ErrorResponse {
	result, err := s.taskRepository.CheckPriority(ctx, taskPayload.Priority)
	if err != nil {
		return utils.InternalServerErrorResponse()
	}
	if !result {
		return utils.NewErrorResponse("Invalid priority", http.StatusBadRequest)
	}

	result, err = s.taskRepository.UpdateTask(ctx, taskPayload)
	if err != nil {
		return utils.InternalServerErrorResponse()
	}
	if !result {
		return utils.NewErrorResponse("Task not found", http.StatusNotFound)
	}

	return nil
}

func (s *DefaultTaskService) DeleteTask(ctx context.Context, taskId uuid.UUID) *utils.ErrorResponse {
	result, err := s.taskRepository.DeleteTask(ctx, taskId)
	if err != nil {
		return utils.InternalServerErrorResponse()
	}
	if !result {
		return utils.NewErrorResponse("Task not found", http.StatusNotFound)
	}

	return nil
}

func NewDefaultTaskService(taskRepository repositories.TaskRepository) *DefaultTaskService {
	return &DefaultTaskService{taskRepository}
}
