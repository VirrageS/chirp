package api

import (
	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/model"
	appErrors "github.com/VirrageS/chirp/backend/model/errors"
	"net/http"
)

func (api *API) createTokens(userID int64, request *http.Request) (string, string, error) {
	authToken, err := api.tokenManager.CreateAuthToken(userID, request)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := api.tokenManager.CreateRefreshToken(userID, request)
	if err != nil {
		return "", "", err
	}

	return authToken, refreshToken, nil
}

func (api *API) refreshAuthToken(requestData *model.RefreshAuthTokenRequest, request *http.Request) (*model.RefreshAuthTokenResponse, error) {
	userID, err := api.tokenManager.ValidateToken(requestData.RefreshToken, request)
	if err != nil {
		log.WithError(err).Error("Error validating token in refreshAuthToken.")
		return nil, err
	}

	// check if authenticating user exists
	_, err = api.service.GetUser(userID, userID)
	if err == appErrors.NoResultsError {
		return nil, appErrors.NotExistingUserAuthenticatingError
	}
	if err != nil {
		return nil, err
	}

	// generate new auth token for the user
	authToken, err := api.tokenManager.CreateAuthToken(userID, request)
	if err != nil {
		return nil, err
	}

	response := &model.RefreshAuthTokenResponse{
		AuthToken: authToken,
	}

	return response, nil
}
