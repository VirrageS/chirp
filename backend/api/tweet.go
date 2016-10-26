package api

import (
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
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	context.JSON(http.StatusOK, tweets)
}

func GetTweet(context *gin.Context) {
	parameterID := context.Query("id")

	tweetID, err := strconv.ParseInt(parameterID, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tweet ID.",
		})
		return
	}

	responseTweet, err := services.GetTweet(tweetID)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Tweet with given ID not found.",
		})
		return
	}

	context.JSON(http.StatusOK, responseTweet)
}

func PostTweet(context *gin.Context) {
	tweetAuthorIDString := context.PostForm("author_id")
	content := context.PostForm("content")

	tweetAuthorID, err := strconv.ParseInt(tweetAuthorIDString, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Author ID was not a number.",
		})
		return
	}

	requestTweet := model.NewTweet{
		AuthorID: tweetAuthorID,
		Content:  content,
	}

	responseTweet, err := services.PostTweet(requestTweet)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", responseTweet.ID))
	context.JSON(http.StatusCreated, responseTweet)
}
