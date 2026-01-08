// Package error is responsible for centralize and organize
// business error logic with a specific code and message.
package errors

import (
	"encoding/json"
	baseerror "errors"
	"fmt"
)

const MaxStackDepth = 50

type (
	StackTracer interface {
		Stack() []byte
	}

	StatusCoder interface {
		StatusCode() int
	}

	Err interface {
		error
		StackTrace() []uintptr
		Stack() []byte
		Unwrap() error
	}

	Error struct {
		Code     string `json:"code"`
		Message  string `json:"message"`
		Status   int    `json:"status"`
		canRetry bool
		Err      error `json:"-"`
		skip     int
	}
)

type Option func(*Error)

func New(code string, options ...Option) error {
	newErr := &Error{
		Code:     code,
		Status:   500,
		canRetry: true,
	}

	for _, o := range options {
		o(newErr)
	}

	return newErr
}

func (e Error) Error() string {
	if e.Message == "" {
		return e.Code
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.Code,
		Message: e.Message,
	})
}

func (e Error) StatusCode() int {
	return e.Status
}

func (e Error) As(target any) bool {
	if targetErr, ok := target.(*Error); ok {
		return targetErr.Code == e.Code
	}

	return baseerror.As(e.Unwrap(), target)
}

func (e Error) Is(target error) bool {
	if targetErr, ok := target.(*Error); ok {
		return targetErr.Code == e.Code
	}

	return baseerror.Is(e.Unwrap(), target)
}

func Code(baseErr error) string {
	var err Error
	if baseerror.As(baseErr, &err) {
		return err.Code
	}

	return ""
}
