package domain

type ErrorType string

const (
	ErrInvalidInput ErrorType = "invalid_input"
	ErrNotFound     ErrorType = "not_found"
	ErrConflict     ErrorType = "conflict"
	ErrUnauthorized ErrorType = "unauthorized"
	ErrInternal     ErrorType = "internal"
)

type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) Unwrap() error {
	return e.Err
}

func NewError(tp ErrorType, message string) Error {
	return Error{
		Type:    tp,
		Message: message,
	}
}

func WrapError(tp ErrorType, message string, err error) Error {
	return Error{
		Type:    tp,
		Message: message,
		Err:     err,
	}
}
