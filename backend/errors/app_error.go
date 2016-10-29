package errors

import (
	"errors"
)

type AppError struct {
	Code int
	Err  error
}

var UnexpectedError = &AppError{
	Code: 500,
	Err:  errors.New("Unexpected error."),
}

func (err AppError) Error() string {
	return err.Err.Error()
}
