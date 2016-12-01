package api

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
)

func (api *API) RegisterUser(context *gin.Context) {
	var newUserForm model.NewUserForm
	if err := context.BindJSON(&newUserForm); err != nil {
		context.AbortWithError(http.StatusBadRequest,
			errors.New("Fields: name, username, password and email are required."))

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
		context.AbortWithError(http.StatusBadRequest,
			errors.New("Fields: email and password are required."))

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
