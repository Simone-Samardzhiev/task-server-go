package utils

import "github.com/gofiber/fiber/v2"

// ValidatablePayload interface is used for payload that needs to be validated.
type ValidatablePayload interface {
	// ValidatePayload return error response if the payload is invalid with
	// specified message. If the payload is valid it must return nil.
	ValidatePayload() *ErrorResponse
}

// HandlePayload used to check the payload. If the payload is valid the result is true
// otherwise the function responds with [ErrorResponse] and return false.
func HandlePayload(c *fiber.Ctx, payload ValidatablePayload) bool {
	errorResponse := payload.ValidatePayload()
	if errorResponse == nil {
		return true
	}

	err := c.JSON(errorResponse)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
	}
	return false
}
