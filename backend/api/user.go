package api

import (
	"errors"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

func (api *API) GetUsers(context *gin.Context) {
	users, err := api.Service.GetUsers()
	if err != nil {
		context.AbortWithError(err.Code, err.Err)
		return
	}

	context.IndentedJSON(http.StatusOK, users)
}

func (api *API) GetUser(context *gin.Context) {
	parameterID := context.Param("id")

	userID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID. Expected an integer."))
		return
	}

	responseUser, err2 := api.Service.GetUser(userID)
	if err2 != nil {
		context.AbortWithError(err2.Code, err2.Err)
		return
	}

	context.IndentedJSON(http.StatusOK, responseUser)
}
