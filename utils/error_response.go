package utils

import "net/http"

// ErrorResponse is a response returned when the payload is empty.
type ErrorResponse struct {
	// The message that should be returned
	Message string `json:"message"`
	// The status code of the error.
	StatusCode int `json:"status_code"`
}

func NewErrorResponse(message string, statusCode int) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		StatusCode: statusCode,
	}
}

func HandleErrorResponse(w http.ResponseWriter, response *ErrorResponse) bool {
	if response == nil {
		return false
	}

	w.WriteHeader(response.StatusCode)
	_, err := w.Write([]byte(response.Message))
	if err != nil {
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
	}
	return false
}
