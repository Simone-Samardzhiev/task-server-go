package utils

import "net/http"

// Payload is an interface used to validate a payload.
type Payload interface {
	// ValidatePayload is the method that will validate the payload.
	ValidatePayload() *ErrorResponse
}

// HandlePayload will check the payload.
// If the payload is not valid it will respond with the [ErrorResponse].
// It returns true if the handler have responded and false if the payload is valid
func HandlePayload(w http.ResponseWriter, payload Payload) bool {
	return HandleErrorResponse(w, payload.ValidatePayload())
}
