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
	if bindError := context.BindJSON(&newUserForm); bindError != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Fields: name, username, password and email are required."))
		return
	}

	newUser, serviceError := api.Service.RegisterUser(&newUserForm)
	if serviceError != nil {
		context.AbortWithError(serviceError.Code, serviceError.Err)
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", newUser.ID))
	context.IndentedJSON(http.StatusCreated, newUser)
}

func (api *API) LoginUser(context *gin.Context) {
	var loginForm model.LoginForm
	if bindError := context.BindJSON(&loginForm); bindError != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Fields: email and password are required."))
		return
	}

	loginResponse, serviceError := api.Service.LoginUser(&loginForm)
	if serviceError != nil {
		context.AbortWithError(serviceError.Code, serviceError.Err)
		return
	}

	context.IndentedJSON(http.StatusOK, loginResponse)
}
