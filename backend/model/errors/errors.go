package errors

import "errors"

var UnexpectedError = errors.New("Internal server error.")

var NoResultsError = errors.New("Not found.")
var UserAlreadyExistsError = errors.New("User with given username or email already exists.")

var ForbiddenError = errors.New("User is not allowed to modify this resource.")
var InvalidCredentialsError = errors.New("Invalid email or password.")

var NotExistingUserAuthenticatingError = errors.New("User authenticating with auth token of a user that does not exist.")

var NoUserAgentHeaderError = errors.New("User-Agent header is required in request for API authorization.")
