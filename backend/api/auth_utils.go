package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/model"
	appErrors "github.com/VirrageS/chirp/backend/model/errors"
)

func (api *API) getGoogleUser(loginFormCode string) (*model.UserGoogle, error) {
	token, err := api.googleOAuth2.Exchange(context.Background(), loginFormCode)
	if err != nil {
		log.WithError(err).Error("Error exchanging with Google.")
		return nil, err
	}

	client := api.googleOAuth2.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.WithError(err).Error("Error getting user info from Google.")
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	user := model.UserGoogle{}
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		log.WithError(err).Error("Error unmarshalling user info from Google response.")
		return nil, err
	}

	return &user, nil
}

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
