package api

import (
	"errors"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

func (api *API) GetUsers(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))

	users, err := api.Service.GetUsers(requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, users)
}

func (api *API) GetUser(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

	userID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID. Expected an integer."))
		return
	}

	user, err := api.Service.GetUser(userID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}

func (api *API) FollowUser(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

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

func (api *API) UnfollowUser(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

	userID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID. Expected an integer."))
		return
	}

	user, err := api.Service.UnfollowUser(userID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}

func (api *API) UserFollowers(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

	userID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID. Expected an integer."))
		return
	}

	user, err := api.Service.UserFollowers(userID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}

func (api *API) UserFollowees(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

	userID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID. Expected an integer."))
		return
	}

	user, err := api.Service.UserFollowees(userID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}
