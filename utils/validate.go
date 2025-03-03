package utils

import "net/http"

// ValidateResponse is a response returned when the payload is empty.
type ValidateResponse struct {
	// The message that should be returned
	Message string `json:"message"`
	// The status code of the error.
	StatusCode int `json:"status_code"`
}

// Payload is an interface used to validate a payload.
type Payload interface {
	// ValidatePayload is the method that will validate the payload.
	ValidatePayload() *ValidateResponse
}

// HandlePayload will take the check the payload.
// If the payload is not valid it will respond with the [ValidateResponse].
// It returns true if the handler have responded and false if the payload is valid
func HandlePayload(r *http.Request, w http.ResponseWriter, payload Payload) bool {
	response := payload.ValidatePayload()
	if response != nil {
		w.WriteHeader(response.StatusCode)
		_, err := w.Write([]byte(response.Message))
		if err != nil {
			http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		}
		return false
	}

	return true
}
