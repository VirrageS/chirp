package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/services"
	"gopkg.in/gin-gonic/gin.v1"
)

func RegisterUser(context *gin.Context) {
	var newUserForm model.NewUserForm
	if bindError := context.BindJSON(&newUserForm); bindError != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Fields: name, username, password and email are required."))
		return
	}

	newUser, serviceError := services.RegisterUser(newUserForm)
	if serviceError != nil {
		context.AbortWithError(serviceError.Code, serviceError.Err)
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", newUser.ID))
	context.IndentedJSON(http.StatusCreated, gin.H{
		"user": newUser,
	})
}

func LoginUser(context *gin.Context) {
	var loginForm model.LoginForm
	if bindError := context.BindJSON(&loginForm); bindError != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Fields: email and password are required."))
	}

	token, serviceError := services.LoginUser(loginForm)
	if serviceError != nil {
		context.AbortWithError(serviceError.Code, serviceError.Err)
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"auth_token": token,
	})
}
