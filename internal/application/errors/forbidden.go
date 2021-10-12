package errors

import "fmt"

type ForbiddenError struct {
	Message string
}

func NewForbiddenError(message ...string) *ForbiddenError {
	defaultMessage := "Forbidden"

	err := &ForbiddenError{
		Message: defaultMessage,
	}

	if len(message) > 0 {
		err.Message = fmt.Sprintf("%s: %s", defaultMessage, message[0])
	}

	return err
}

func (err *ForbiddenError) Error() string {
	return err.Message
}
