package errors

import "errors"

var NoResultsError = errors.New("Not found.")
var UnexpectedError = errors.New("Internal server error.")

var UserAlreadyExistsError = errors.New("User with given username or email already exists.")

var ForbiddenError = errors.New("User is not allowed to modify this resource.")
var InvalidCredentialsError = errors.New("Invalid email or password.")

var NotExistingUserAuthenticatingError = errors.New("User authenticating with auth token of a user that does not exist.")
