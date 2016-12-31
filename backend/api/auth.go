package api

import (
	"errors"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/config"
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

	newUser, err := api.Service.RegisterUser(&newUserForm)
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

	loginResponse, err := api.Service.LoginUser(&loginForm)
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

	response, err := api.Service.RefreshAuthToken(&request)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, response)
}

func GetoAuth2ConfigGoogle() *oauth2.Config {
	_,_,_, authorizationGoogleConfig := config.GetConfig("config")

	return &oauth2.Config{
		ClientID:     authorizationGoogleConfig.GetClientId(),
		ClientSecret: authorizationGoogleConfig.GetClientSecret(),
		RedirectURL:  "http://127.0.0.1:8080/createOrLoginUserWithGoogle",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: oauth2.Endpoint{
				AuthURL:  authorizationGoogleConfig.GetAuthURL(),
				TokenURL: authorizationGoogleConfig.GetTokenURL(),
		},
	}
}

func (api *API) AddressAuthorizationGoogle(context *gin.Context) {
	conf := GetoAuth2ConfigGoogle()
	token := "TODO"
	context.IndentedJSON(http.StatusOK, conf.AuthCodeURL(token))
}

func (api *API) CreateOrLoginUserWithGoogle(context *gin.Context) {
	if "TODO" != context.Query("state") {
		context.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state"))
		return
	}
	conf := GetoAuth2ConfigGoogle()
	token, err := conf.Exchange(oauth2.NoContext, context.Query("code"))
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	user := model.UserGoogle{}
	json.Unmarshal([]byte(data), &user)

	loginResponse, err := api.Service.CreateOrLoginUserWithGoogle(&user)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}
	context.IndentedJSON(http.StatusOK, loginResponse)
}
