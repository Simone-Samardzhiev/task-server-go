package services

import (
	"context"
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

func NewDefaultTaskService(taskRepository repositories.TaskRepository) *DefaultTaskService {
	return &DefaultTaskService{taskRepository}
}
