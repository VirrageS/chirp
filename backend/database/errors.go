package database

import "errors"

const UniqueConstraintViolationCode = "23505"

var NoResults = errors.New("No results found in database.")
var DatabaseError = errors.New("Database error.")
var UserAlreadyExistsError = errors.New("User with given username or email already exists.")
