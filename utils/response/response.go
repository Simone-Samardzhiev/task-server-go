package response

import (
	"net/http"
)

// Response type struct holds Ok only if Err is nil.
type Response[T any] struct {
	Ok  SuccessResponse[T]
	Err ErrorResponse
}

// NewResponseWithErr return Response where Ok is default value.
func NewResponseWithErr[T any](err ErrorResponse) Response[T] {
	return Response[T]{
		Ok:  SuccessResponse[T]{},
		Err: err,
	}
}

// NewResponseWithOk return Response where Err is nil.
func NewResponseWithOk[T any](ok *SuccessResponse[T]) Response[T] {
	return Response[T]{
		Ok:  *ok,
		Err: nil,
	}
}

func NewIntervalServerErrorResponse[T any]() Response[T] {
	return Response[T]{
		Ok: SuccessResponse[T]{},
		Err: NewErrorResponse(
			"Internal Server Error",
			http.StatusInternalServerError,
		),
	}
}

func NewUnauthorizedErrorResponse[T any]() Response[T] {
	return Response[T]{
		Ok: SuccessResponse[T]{},
		Err: NewErrorResponse(
			"Unauthorized,",
			http.StatusUnauthorized,
		),
	}
}
