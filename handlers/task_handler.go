package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"server/auth/tokens"
	"server/models"
	"server/services"
	"server/utils"
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
		claims, ok := c.Locals(tokens.JWTClaimsKey).(*tokens.Token)
		if !ok {
			utils.HandleErrorResponse(c, utils.InternalServerErrorResponse())
		}

		tasks, err := h.taskService.GetTasks(c.Context(), *claims)
		if !utils.HandleErrorResponse(c, err) {
			return nil
		}

		return c.JSON(tasks)
	}
}

func (h *DefaultTaskHandler) AddTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(tokens.JWTClaimsKey).(*tokens.Token)
		if !ok {
			utils.HandleErrorResponse(c, utils.InternalServerErrorResponse())
		}

		var task models.NewTaskPayload
		if err := c.BodyParser(&task); err != nil {
			return err
		}

		newTask, err := h.taskService.AddTask(c.Context(), *claims, &task)
		if !utils.HandleErrorResponse(c, err) {
			return nil
		}

		return c.JSON(newTask)
	}
}

func (h *DefaultTaskHandler) UpdateTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var task models.TaskPayload
		if err := c.BodyParser(&task); err != nil {
			return err
		}

		err := h.taskService.UpdateTask(c.Context(), &task)
		if !utils.HandleErrorResponse(c, err) {
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func (h *DefaultTaskHandler) DeleteTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			utils.HandleErrorResponse(c, utils.NewErrorResponse("Invalid uuid", fiber.StatusBadRequest))
			return nil
		}

		errorResponse := h.taskService.DeleteTask(c.Context(), parsedId)
		if !utils.HandleErrorResponse(c, errorResponse) {
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}

func NewDefaultTaskHandler(taskService services.TaskService) *DefaultTaskHandler {
	return &DefaultTaskHandler{taskService}
}
