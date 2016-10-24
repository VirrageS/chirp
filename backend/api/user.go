package api

import (
	"fmt"
	"github.com/VirrageS/chirp/backend/apiModel"
	"github.com/VirrageS/chirp/backend/services"
	"github.com/kataras/iris"
	"strconv"
)

func GetUsers(context *iris.Context) {
	users, err := services.GetUsers()

	if err != nil {
		context.JSON(iris.StatusInternalServerError, iris.Map{
			"error": err,
		})
		return
	}

	context.JSON(iris.StatusOK, users)
}

func GetUser(context *iris.Context) {
	parameterId := context.Param("id")
	userId, err := strconv.ParseInt(parameterId, 10, 64)
	if err != nil {
		context.JSON(iris.StatusBadRequest, iris.Map{
			"error": "Invalid user ID.",
		})
		return
	}

	responseUser, err := services.GetUser(userId)

	if err != nil {
		context.JSON(iris.StatusNotFound, iris.Map{
			"error": "User with given ID not found.",
		})
		return
	}

	context.JSON(iris.StatusOK, responseUser)
}

// TODO: now returns 404 if user already exists
func PostUser(context *iris.Context) {
	name := context.PostValue("name")
	username := context.PostValue("username")
	email := context.PostValue("email")

	requestUser := apiModel.NewUser{
		Name:     name,
		Username: username,
		Email:    email,
	}

	newUser, err := services.PostUser(requestUser)

	if err != nil {
		context.JSON(iris.StatusNotFound, iris.Map{
			"error": err,
		})
		return
	}

	context.SetHeader("Location", fmt.Sprintf("/user/%d", newUser.Id))
	context.JSON(iris.StatusCreated, newUser)
}
