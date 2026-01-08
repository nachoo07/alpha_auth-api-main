package session

type TokenExpiredError struct {
	Message string
}

func (e TokenExpiredError) Error() string {
	return e.Message
}
