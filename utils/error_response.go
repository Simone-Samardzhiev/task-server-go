package utils

// ErrorResponse type interface is a conventional way to represent
// error response.
type ErrorResponse interface {
	Message() string
	StatusCode() int
}

// errorResponse type struct is implementation of ErrorResponse.
type errorResponse struct {
	message string
	status  int
}

func (e *errorResponse) Message() string {
	return e.message
}

func (e *errorResponse) StatusCode() int {
	return e.status
}

// NewErrorResponse will return an ErrorResponse with error message and status code.
func NewErrorResponse(message string, status int) ErrorResponse {
	return &errorResponse{message, status}
}
