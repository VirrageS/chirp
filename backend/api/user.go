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
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
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

	user, err := api.Service.GetUser(userID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}

func (api *API) FollowUser(context *gin.Context) {
	parameterID := context.Param("id")
	requestingUserID := (context.MustGet("userID").(int64))

	userID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID. Expected an integer."))
		return
	}
	if userID == requestingUserID {
		context.AbortWithError(http.StatusBadRequest, errors.New("User can't follow himself."))
		return
	}

	user, err := api.Service.FollowUser(userID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}
