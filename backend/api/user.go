package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/services"
)

func GetUsers(context *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		context.AbortWithError(err.Code, err.Err)
		return
	}

	context.JSON(http.StatusOK, users)
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

	context.JSON(http.StatusOK, responseUser)
}

func PostUser(context *gin.Context) {
	name := context.PostForm("name")
	username := context.PostForm("username")
	email := context.PostForm("email")

	err := validatePostUserParameters(name, username, email)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
	}

	requestUser := model.NewUser{
		Name:     name,
		Username: username,
		Email:    email,
	}

	newUser, err2 := services.PostUser(requestUser)
	if err2 != nil {
		context.AbortWithError(err2.Code, err2.Err)
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", newUser.ID))
	context.JSON(http.StatusCreated, newUser)
}

func validatePostUserParameters(name, username, email string) error {
	var invalidFields []string

	if name == "" {
		invalidFields = append(invalidFields, "name")
	}
	if username == "" {
		invalidFields = append(invalidFields, "username")
	}
	if email == "" {
		invalidFields = append(invalidFields, "email")
	}

	if len(invalidFields) > 0 {
		errorMesssage := "Invalid request. Fields: " + strings.Join(invalidFields, ", ") + " are required."
		return errors.New(errorMesssage)
	}

	return nil
}
