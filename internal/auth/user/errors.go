package user

type ResourceNotFoundError struct {
	Message string
}

func (e ResourceNotFoundError) Error() string {
	return e.Message
}

type UserNotVerifyError struct {
	Message string
}

func (e UserNotVerifyError) Error() string {
	return e.Message
}

type InvalidCredentialsError struct {
	Message string
}

func (e InvalidCredentialsError) Error() string {
	return e.Message
}
