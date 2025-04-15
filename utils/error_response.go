package utils

import "github.com/gofiber/fiber/v2"

// ErrorResponse is the standard way of return error
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// NewErrorResponse creates new instance of [ErrorResponse]
func NewErrorResponse(message string, status int) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Status:  status,
	}
}

// HandleErrorResponse will return true if the error is  not nil
// and the function responded.
func HandleErrorResponse(c *fiber.Ctx, error *ErrorResponse) bool {
	if error == nil {
		return true
	}

	if err := c.Status(error.Status).JSON(error); err != nil {
		c.Status(fiber.StatusInternalServerError)
	}
	return false
}

// InternalServerErrorResponse is the standard error return when there is a server error.
func InternalServerErrorResponse() *ErrorResponse {
	return NewErrorResponse("Internal Server Error", 500)
}

// InvalidTokenErrorResponse is the standard error returned when the token is invalid.
func InvalidTokenErrorResponse() *ErrorResponse {
	return NewErrorResponse("Invalid Token", 401)
}
