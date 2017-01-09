package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/model/errors"
)

var errorToCodeMap = map[error]int{
	errors.NoResultsError:                     http.StatusNotFound,
	errors.UnexpectedError:                    http.StatusInternalServerError,
	errors.UserAlreadyExistsError:             http.StatusConflict,
	errors.ForbiddenError:                     http.StatusForbidden,
	errors.InvalidCredentialsError:            http.StatusUnauthorized,
	errors.NotExistingUserAuthenticatingError: http.StatusBadRequest,
	errors.NoUserAgentHeaderError:             http.StatusBadRequest,
}

func getStatusCodeFromError(err error) int {
	if code, ok := errorToCodeMap[err]; ok {
		return code
	}

	log.WithError(err).Error("Api received an unexpected error type.")
	return http.StatusInternalServerError
}
