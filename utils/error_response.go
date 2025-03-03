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

func HandleErrorResponse(w http.ResponseWriter, response *ErrorResponse) bool {
	if response == nil {
		return false
	}

	w.WriteHeader(response.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "There was an unknwon error", http.StatusInternalServerError)
	}

	return false
}

func InternalServerError() *ErrorResponse {
	return &ErrorResponse{
		Message:    "Internal Server Error",
		StatusCode: http.StatusInternalServerError,
	}
}

func InvalidJson() *ErrorResponse {
	return &ErrorResponse{
		Message:    "Invalid Json",
		StatusCode: http.StatusBadRequest,
	}
}
