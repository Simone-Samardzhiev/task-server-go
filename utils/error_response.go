package utils

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
