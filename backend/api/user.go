package api

import (
	"errors"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/services"
)

func GetUsers(context *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		context.AbortWithError(err.Code, err.Err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func GetUser(context *gin.Context) {
	parameterID := context.Param("id")

	userID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID. Expected an integer."))
		return
	}

	responseUser, err2 := services.GetUser(userID)
	if err2 != nil {
		context.AbortWithError(err2.Code, err2.Err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"user": responseUser,
	})
}
