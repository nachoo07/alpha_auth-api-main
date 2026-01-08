package errors

import (
	"context"
	baseError "errors"
)

func HTTPOutput(ctx context.Context, e error) error {
	var newErr *Error
	if baseError.As(e, &newErr) {
		return newErr
	}

	return InternalServerError(WithBaseError(e))
}
