package errors

import (
	"fmt"
	"net/http"
)

type UnauthorizedError struct {
	Message string
}

func NewUnauthorizedError(message ...string) *UnauthorizedError {
	defaultMessage := "Unauthorized"

	err := &UnauthorizedError{
		Message: defaultMessage,
	}

	if len(message) > 0 {
		err.Message = fmt.Sprintf("%s: %s", defaultMessage, message[0])
	}

	return err
}

func (err *UnauthorizedError) Error() string {
	return err.Message
}

func (*UnauthorizedError) Code() int {
	return http.StatusUnauthorized
}
