package utils

import "net/http"

// Response is a struct used to respond to a http request.
type Response[T any] struct {
	// The status code of the response.
	StatusCode int
	// The data of the response.
	Data T
}

var (
	InternalServerErrorResponse = Response[string]{
		StatusCode: http.StatusInternalServerError,
		Data:       "Internal Server Error",
	}

	UnauthorizedResponse = Response[string]{
		StatusCode: http.StatusUnauthorized,
		Data:       "Unauthorized",
	}
)
