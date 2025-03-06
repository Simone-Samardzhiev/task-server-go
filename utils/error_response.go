package utils

import (
	"encoding/json"
	"net/http"
)

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

// HandleErrorResponse will return true if the [ErrorResponse] is not nil and
// it has responded otherwise false.
func HandleErrorResponse(w http.ResponseWriter, response *ErrorResponse) bool {
	if response == nil {
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "There was an unknown error", http.StatusInternalServerError)
		println("here")
		return true
	}

	return true
}

func InternalServerError() *ErrorResponse {
	return NewErrorResponse("Internal Server Error", http.StatusInternalServerError)
}

func InvalidJson() *ErrorResponse {
	return NewErrorResponse("Invalid Json", http.StatusBadRequest)
}

func InvalidToken() *ErrorResponse {
	return NewErrorResponse("Invalid token", http.StatusUnauthorized)
}
