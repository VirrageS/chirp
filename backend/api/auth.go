package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/services"
	"gopkg.in/gin-gonic/gin.v1"
)

func RegisterUser(context *gin.Context) {
	name := context.PostForm("name")
	username := context.PostForm("username")
	email := context.PostForm("email")
	password := context.PostForm("password")

	validationError := validateRegisterUserParameters(name, username, email, password)
	if validationError != nil {
		context.AbortWithError(http.StatusBadRequest, validationError)
		return
	}

	requestUser := model.NewUser{
		Username: username,
		Password: password,
		Email:    email,
		Name:     name,
	}

	newUser, serviceError := services.RegisterUser(requestUser)
	if serviceError != nil {
		context.AbortWithError(serviceError.Code, serviceError.Err)
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", newUser.ID))
	context.JSON(http.StatusCreated, gin.H{
		"user": newUser,
	})
}

func LoginUser(context *gin.Context) {
	email := context.PostForm("email")
	password := context.PostForm("password")

	err := validateLoginUserParameters(email, password)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, serviceError := services.LoginUser(email, password)
	if serviceError != nil {
		context.AbortWithError(serviceError.Code, serviceError.Err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"auth_token": token,
	})
}

func validateRegisterUserParameters(name, username, email, password string) error {
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
	if password == "" {
		invalidFields = append(invalidFields, "password")
	}

	if len(invalidFields) > 0 {
		errorMessage := "Invalid request. Fields: " + strings.Join(invalidFields, ", ") + " are required."
		return errors.New(errorMessage)
	}

	return nil
}

func validateLoginUserParameters(email, password string) error {
	var invalidFields []string

	if email == "" {
		invalidFields = append(invalidFields, "email")
	}
	if password == "" {
		invalidFields = append(invalidFields, "password")
	}

	if len(invalidFields) > 0 {
		errorMessage := "Invalid request. Fields " + strings.Join(invalidFields, ", ") + " are required."
		return errors.New(errorMessage)
	}

	return nil
}
