package errors

import "net/http"

func InputValidationError(options ...Option) error {
	newOptions := append(options, WithStatus(http.StatusBadRequest), WithSkip())
	return New("input_validation", newOptions...)
}

func InternalValidationError(options ...Option) error {
	newOptions := append(options, WithStatus(http.StatusBadRequest), WithSkip())
	return New("internal_validation", newOptions...)
}

func InternalServerError(options ...Option) error {
	newOptions := append(options, WithSkip())
	return New("internal_server_error", newOptions...)
}
