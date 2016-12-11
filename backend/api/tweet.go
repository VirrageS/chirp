package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
)

func (api *API) GetTweets(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))
	userIDStr := context.Query("userID")
	var tweets []*model.Tweet
	var err error

	// TODO: maybe put the logic of checking request parameters to service when there is more than 1 flag
	if userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			context.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID. Expected an integer."))
			return
		}
		tweets, err = api.Service.GetTweetsOfUserWithID(userID, requestingUserID)
	} else {
		tweets, err = api.Service.GetTweets(requestingUserID)
	}

	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, tweets)
}

func (api *API) GetTweet(context *gin.Context) {
	requestingUserID := (context.MustGet("userID").(int64))
	parameterID := context.Param("id")

	tweetID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid tweet ID. Expected an integer."))
		return
	}

	responseTweet, err := api.Service.GetTweet(tweetID, requestingUserID)
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

	responseTweet, err := api.Service.PostTweet(&newTweet, requestingUserID)
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

	err = api.Service.DeleteTweet(tweetID, requestingUserID)

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

	tweet, err := api.Service.LikeTweet(tweetID, requestingUserID)
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

	tweet, err := api.Service.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, tweet)
}

func (api *API) HomeFeed(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG in token_auth middleware
	requestingUserID := (context.MustGet("userID").(int64))

	tweets, err := api.Service.GetTweetsOfUserWithID(requestingUserID, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, tweets)
}
