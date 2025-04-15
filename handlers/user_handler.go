package handlers

import (
	"github.com/gofiber/fiber/v2"
	"server/services"
)

// UserHandler interface handles user requests.
type UserHandler interface {
	Register() fiber.Handler
	Login() fiber.Handler
	Refresh() fiber.Handler
}

// DefaultUserHandler interface is the default implementation of [UserHandler]
type DefaultUserHandler struct {
	userService services.UserService
}

func (h *DefaultUserHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		panic("implement me")
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
