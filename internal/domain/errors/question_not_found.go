package errors

import "fmt"

type QuestionNotFoundError struct {
	Message string
}

func NewQuestionNotFoundError(message ...string) *QuestionNotFoundError {
	defaultMessage := "Not found"

	err := &QuestionNotFoundError{
		Message: defaultMessage,
	}

	if len(message) > 0 {
		err.Message = fmt.Sprintf("%s: %s", defaultMessage, message[0])
	}

	return err
}

func (err *QuestionNotFoundError) Error() string {
	return err.Message
}
