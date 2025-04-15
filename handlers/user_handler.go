package handlers

import (
	"github.com/gofiber/fiber/v2"
	"server/auth/tokens"
	"server/models"
	"server/services"
	"server/utils"
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

		if !utils.HandlePayload(c, &payload) {
			return nil
		}

		err := h.userService.Register(c.Context(), payload)
		if !utils.HandleErrorResponse(c, err) {
			return nil
		}

		c.Status(fiber.StatusBadRequest)
		return nil
	}
}

func (h *DefaultUserHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload models.LoginPayload
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		tokenGroup, err := h.userService.Login(c.Context(), payload)
		if !utils.HandleErrorResponse(c, err) {
			return nil
		}

		return c.JSON(tokenGroup)
	}
}

func (h *DefaultUserHandler) Refresh() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(tokens.JWTClaimsKey).(*tokens.Token)
		if !ok {
			utils.HandleErrorResponse(c, utils.InternalServerErrorResponse())
		}

		tokenGroup, err := h.userService.Refresh(c.Context(), *claims)
		if !utils.HandleErrorResponse(c, err) {
			return nil
		}
		return c.JSON(tokenGroup)
	}
}

func NewDefaultUserHandler(userRepository services.UserService) *DefaultUserHandler {
	return &DefaultUserHandler{
		userService: userRepository,
	}
}
