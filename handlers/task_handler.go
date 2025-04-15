package handlers

import (
	"github.com/gofiber/fiber/v2"
	"server/services"
)

// TaskHandler handles tasks request.
type TaskHandler interface {
	// GetTasks will return all tasks of a user.
	GetTasks() fiber.Handler

	// AddTask will add a new task.
	AddTask() fiber.Handler

	// UpdateTask will update an existing task.
	UpdateTask() fiber.Handler

	// DeleteTask will delete an existing task.
	DeleteTask() fiber.Handler
}

// DefaultTaskHandler is the default implementation of [TaskHandler]
type DefaultTaskHandler struct {
	taskService services.TaskService
}

func (h *DefaultTaskHandler) GetTasks() fiber.Handler {
	return func(c *fiber.Ctx) error {
		panic("implement me")
	}
}

func (h *DefaultTaskHandler) AddTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		panic("implement me")
	}
}

func (h *DefaultTaskHandler) UpdateTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		panic("implement me")
	}
}

func (h *DefaultTaskHandler) DeleteTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		panic("implement me")
	}
}

func NewDefaultTaskHandler(taskService services.TaskService) *DefaultTaskHandler {
	return &DefaultTaskHandler{taskService}
}
