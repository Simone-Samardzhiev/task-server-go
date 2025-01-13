package utils

import (
	"errors"
)

var (
	// InternalServerErr is returned when there was unknown error.
	InternalServerErr = errors.New("internal server error")

	// UnauthorizedErr is returned when user credentials are invalid.
	UnauthorizedErr = errors.New("unauthorized")

	// ConflictErr is return when the service cannot proceed with the requests due to conflict of data.
	ConflictErr = errors.New("conflict")
)

const (
	// InternalServerErrorMessage should be sent when there is an [InternalServerErr]
	InternalServerErrorMessage = "Internal Server Error"

	// UnauthorizedMessage should be sent when there is an [UnauthorizedErr]
	UnauthorizedMessage = "Unauthorized"
)
