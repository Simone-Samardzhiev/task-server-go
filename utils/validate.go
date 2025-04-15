package utils

import "github.com/gofiber/fiber/v2"

// ValidatablePayload interface is used for payload that needs to be validated.
type ValidatablePayload interface {
	// ValidatePayload return error response if the payload is invalid with
	// specified message. If the payload is valid it must return nil.
	ValidatePayload() *ErrorResponse
}

// HandlePayload used to check the payload. The function
// calls [HandleErrorResponse] with the error from validation and return the result.
func HandlePayload(c *fiber.Ctx, payload ValidatablePayload) bool {
	errorResponse := payload.ValidatePayload()
	return HandleErrorResponse(c, errorResponse)
}
