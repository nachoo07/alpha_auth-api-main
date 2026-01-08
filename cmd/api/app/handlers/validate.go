package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlphaCodinggroup/alpha_auth-api/pkg/errors"
	"github.com/AlphaCodinggroup/alpha_auth-api/pkg/validator"
)

type DTO interface {
	*loginDTO | *logoutDTO | *userDTO
}

func decodeAndValidate[T DTO](r *http.Request, body T) error {
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		return errors.InputValidationError(errors.WithBaseError(err))
	}

	if err := validator.Validate(body); err != nil {
		return errors.InputValidationError(errors.WithBaseError(err))
	}

	return nil
}
