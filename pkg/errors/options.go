package errors

func WithBaseError(berror error) Option {
	return func(e *Error) {
		if berror == nil {
			return
		}

		e.Err = berror
		if e.Message == "" {
			e.Message = berror.Error()
		}
	}
}

func WithMessage(message string) Option {
	return func(e *Error) {
		e.Message = message
	}
}

func WithStatus(status int) Option {
	return func(e *Error) {
		e.Status = status
	}
}

func WithSkip() Option {
	return func(e *Error) {
		e.skip++
	}
}

func NotRetryable() Option {
	return func(e *Error) {
		e.canRetry = false
	}
}
