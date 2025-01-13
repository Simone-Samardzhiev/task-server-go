package response

// SuccessResponse type struct represent successful response with data and status code.
type SuccessResponse[T any] struct {
	Data       T
	StatusCode int
}

// NewSuccessResponse returns new SuccessResponse with data and status code.
func NewSuccessResponse[T any](data *T, statusCode *int) *SuccessResponse[T] {
	return &SuccessResponse[T]{*data, *statusCode}
}
