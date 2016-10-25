package api

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/VirrageS/chirp/backend/apiModel"
	"github.com/VirrageS/chirp/backend/services"
	"gopkg.in/gin-gonic/gin.v1"
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
	parameterId := context.Query("id")
	userId, err := strconv.ParseInt(parameterId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID.",
		})
		return
	}

	responseUser, err := services.GetUser(userId)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User with given ID not found.",
		})
		return
	}

	context.JSON(http.StatusOK, responseUser)
}

// TODO: now returns 404 if user already exists
func PostUser(context *gin.Context) {
	name := context.PostForm("name")
	username := context.PostForm("username")
	email := context.PostForm("email")

	requestUser := apiModel.NewUser{
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

	context.Header("Location", fmt.Sprintf("/user/%d", newUser.Id))
	context.JSON(http.StatusCreated, newUser)
}
