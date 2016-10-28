package errors

import "errors"

type AppError struct {
	Code int
	Err  error
}

var UnexpectedError = &AppError{
	Code: 500,
	Err:  errors.New("Unexpected error."),
}

func (error AppError) Error() string {
	return error.Err.Error()
}
