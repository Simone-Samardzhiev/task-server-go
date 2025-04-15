package handlers

import (
	"github.com/gofiber/fiber/v2"
	"server/models"
	"server/services"
)

// UserHandler interface handles user requests.
type UserHandler interface {
	// Register handler used to register by user.
	Register() fiber.Handler
	// Login handler used to log in by the user.
	Login() fiber.Handler
	// Refresh handler used to revalidate tokens.
	Refresh() fiber.Handler
}

// DefaultUserHandler interface is the default implementation of [UserHandler]
type DefaultUserHandler struct {
	userService services.UserService
}

func (h *DefaultUserHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload models.RegistrationsPayload
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

	}
}

func (h *DefaultUserHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		panic("implement me")
	}
}

func (h *DefaultUserHandler) Refresh() fiber.Handler {
	return func(c *fiber.Ctx) error {
		panic("implement me")
	}
}

func NewDefaultUserHandler(userRepository services.UserService) *DefaultUserHandler {
	return &DefaultUserHandler{
		userService: userRepository,
	}
}
