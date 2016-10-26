package api

import (
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/services"
)

func GetUsers(context *gin.Context) {
	users, err := services.GetUsers()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	context.JSON(http.StatusOK, users)
}

func GetUser(context *gin.Context) {
	parameterID := context.Query("id")
	userID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID.",
		})
		return
	}

	responseUser, err := services.GetUser(userID)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User with given ID not found.",
		})
		return
	}

	context.JSON(http.StatusOK, responseUser)
}

func PostUser(context *gin.Context) {
	name := context.PostForm("name")
	username := context.PostForm("username")
	email := context.PostForm("email")

	requestUser := model.NewUser{
		Name:     name,
		Username: username,
		Email:    email,
	}

	newUser, err := services.PostUser(requestUser)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", newUser.ID))
	context.JSON(http.StatusCreated, newUser)
}
