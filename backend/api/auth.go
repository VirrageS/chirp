package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
	"golang.org/x/oauth2"
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

	loginResponse, err := api.service.LoginUser(&loginForm)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, loginResponse)
}

func (api *API) RefreshAuthToken(context *gin.Context) {
	var request model.RefreshAuthTokenRequest
	if err := context.BindJSON(&request); err != nil {
		context.AbortWithError(
			http.StatusBadRequest,
			errors.New("Fields: user_id and refresh_token are required."),
		)
		return
	}

	response, err := api.service.RefreshAuthToken(&request)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, response)
}

func (api *API) GetGoogleAutorizationURL(context *gin.Context) {
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
		context.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state"))
		return
	}

	token, err := api.googleOAuth2.Exchange(oauth2.NoContext, form.Code)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := api.googleOAuth2.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	user := model.UserGoogle{}
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	loginResponse, err := api.service.CreateOrLoginUserWithGoogle(&user)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}
	context.IndentedJSON(http.StatusOK, loginResponse)
}
