package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/services"
)

func GetTweets(context *gin.Context) {
	// TODO: support filtering
	//expected_user_id := context.Query("author")
	//expected_user_name := context.Query("author")
	// ...

	tweets, err := services.GetTweets()
	if err != nil {
		context.AbortWithError(err.Code, err.Err)
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"tweets": tweets,
	})
}

func GetTweet(context *gin.Context) {
	parameterID := context.Param("id")

	tweetID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid tweet ID. Expected an integer."))
		return
	}

	responseTweet, err2 := services.GetTweet(tweetID)
	if err2 != nil {
		context.AbortWithError(err2.Code, err2.Err)
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"tweet": responseTweet,
	})
}

func PostTweet(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG in token_auth middleware
	tweetAuthorID := (context.MustGet("userID").(int64))
	var newTweetContent model.NewTweetContent

	if bindError := context.BindJSON(&newTweetContent); bindError != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Field content is required."))
	}

	requestTweet := model.NewTweet{
		AuthorID: tweetAuthorID,
		Content:  newTweetContent.Content,
	}

	responseTweet, err2 := services.PostTweet(requestTweet)
	if err2 != nil {
		context.AbortWithError(err2.Code, err2.Err)
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", responseTweet.ID))
	context.IndentedJSON(http.StatusCreated, gin.H{
		"tweet": responseTweet,
	})
}

func DeleteTweet(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG
	authenticatingUserID := (context.MustGet("userID").(int64))
	tweetIDString := context.Param("id")

	tweetID, parseError := strconv.ParseInt(tweetIDString, 10, 64)
	if parseError != nil {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid tweet ID. Expected an integer."))
		return
	}

	serviceError := services.DeleteTweet(authenticatingUserID, tweetID)

	if serviceError != nil {
		context.AbortWithError(serviceError.Code, serviceError.Err)
		return
	}

	context.Status(http.StatusNoContent)
}

func HomeFeed(context *gin.Context) {
	// for now lets panic when userID is not set, or when its not an int because that would mean a BUG in token_auth middleware
	tweetAuthorID := (context.MustGet("userID").(int64))

	tweets, err := services.GetTweetsOfUserWithID(tweetAuthorID)
	if err != nil {
		context.AbortWithError(err.Code, err.Err)
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"tweets": tweets,
	})
}
