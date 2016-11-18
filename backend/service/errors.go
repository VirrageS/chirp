package service

import (
	"errors"
)

type Error struct {
	Code int
	Err  error
}

var UnexpectedError = &Error{
	Code: 500,
	Err:  errors.New("Unexpected error."),
}

func (err Error) Error() string {
	return err.Err.Error()
}
