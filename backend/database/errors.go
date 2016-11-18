package database

import (
	"errors"
)

const UniqueConstraintViolationCode = "23505"

var NoRowsError = errors.New("No rows matching given query were found.")
var DatabaseError = errors.New("Database error.")
var UserAlreadyExistsError = errors.New("User with given username or email already exists.")
