package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/VirrageS/chirp/backend/model"
)

func (api *API) GetTweet(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

	tweetID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid tweet ID. Expected an integer."))
		return
	}

	responseTweet, err := api.service.GetTweet(tweetID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, responseTweet)
}

func (api *API) PostTweet(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG in token_auth middleware
	requestingUserID := (context.MustGet("userID").(int64))
	var newTweet model.NewTweet

	if err := context.BindJSON(&newTweet); err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Field content is required."))
		return
	}

	newTweet.AuthorID = requestingUserID

	responseTweet, err := api.service.PostTweet(&newTweet, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", responseTweet.ID))
	context.IndentedJSON(http.StatusCreated, responseTweet)
}

func (api *API) DeleteTweet(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG
	requestingUserID := (context.MustGet("userID").(int64))
	tweetIDString := context.Param("id")

	tweetID, err := strconv.ParseInt(tweetIDString, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid tweet ID. Expected an integer."))
		return
	}

	err = api.service.DeleteTweet(tweetID, requestingUserID)

	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.Status(http.StatusNoContent)
}

func (api *API) LikeTweet(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

	tweetID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid tweet ID. Expected an integer."))
		return
	}

	tweet, err := api.service.LikeTweet(tweetID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, tweet)
}

func (api *API) UnlikeTweet(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

	tweetID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid tweet ID. Expected an integer."))
		return
	}

	tweet, err := api.service.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, tweet)
}

func (api *API) Feed(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG in token_auth middleware
	requestingUserID := (context.MustGet("userID").(int64))

	tweets, err := api.service.Feed(requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, tweets)
}
