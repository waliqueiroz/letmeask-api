package errors

import (
	"fmt"
	"net/http"
)

type ResourceNotFoundError struct {
	Message string
}

func NewResourceNotFoundError(message ...string) *ResourceNotFoundError {
	defaultMessage := "Not found"

	err := &ResourceNotFoundError{
		Message: defaultMessage,
	}

	if len(message) > 0 {
		err.Message = fmt.Sprintf("%s: %s", defaultMessage, message[0])
	}

	return err
}

func (err *ResourceNotFoundError) Error() string {
	return err.Message
}

func (*ResourceNotFoundError) Code() int {
	return http.StatusNotFound
}
