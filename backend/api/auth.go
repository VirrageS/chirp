package api

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
)

func (api *API) RegisterUser(context *gin.Context) {
	var newUserForm model.NewUserForm
	if err := context.BindJSON(&newUserForm); err != nil {
		context.AbortWithError(
			http.StatusBadRequest,
			errors.New("Fields: name, username, password and email are required."),
		)
		return
	}

	newUser, err := api.service.RegisterUser(&newUserForm)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", newUser.ID))
	context.IndentedJSON(http.StatusCreated, newUser)
}

func (api *API) LoginUser(context *gin.Context) {
	var loginForm model.LoginForm
	if err := context.BindJSON(&loginForm); err != nil {
		context.AbortWithError(
			http.StatusBadRequest,
			errors.New("Fields: email and password are required."),
		)
		return
	}

	loggedUser, err := api.service.LoginUser(&loginForm)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	authToken, refreshToken, err := api.createTokens(loggedUser.ID, context.Request)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	loginResponse := &model.LoginResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
		User:         loggedUser,
	}

	context.IndentedJSON(http.StatusOK, loginResponse)
}

func (api *API) RefreshAuthToken(context *gin.Context) {
	var requestData model.RefreshAuthTokenRequest
	if err := context.BindJSON(&requestData); err != nil {
		context.AbortWithError(
			http.StatusBadRequest,
			errors.New("Fields: `user_id` and `refresh_token` are required."),
		)
		return
	}

	response, err := api.refreshAuthToken(&requestData, context.Request)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, response)
}

func (api *API) GetGoogleAuthorizationURL(context *gin.Context) {
	token := "TODO" // TODO: this should be generated hash from IP Address and browser name / browser_id
	context.IndentedJSON(http.StatusOK, api.googleOAuth2.AuthCodeURL(token, oauth2.AccessTypeOffline))
}

func (api *API) CreateOrLoginUserWithGoogle(context *gin.Context) {
	var form model.GoogleLoginForm
	if err := context.BindJSON(&form); err != nil {
		context.AbortWithError(
			http.StatusBadRequest,
			errors.New("Fields: `code` and `state` are required."),
		)
		return
	}

	if form.State != "TODO" {
		context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid Google login form."))
		return
	}

	user, err := api.getGoogleUser(form.Code)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Error fetching user from Google."))
	}

	loggedUser, err := api.service.CreateOrLoginUserWithGoogle(user)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	authToken, refreshToken, err := api.createTokens(loggedUser.ID, context.Request)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	loginResponse := &model.LoginResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
		User:         loggedUser,
	}

	context.IndentedJSON(http.StatusOK, loginResponse)
}
