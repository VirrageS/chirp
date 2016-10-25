package api

import (
	"fmt"
	"net/http"
	"github.com/VirrageS/chirp/backend/apiModel"
	"github.com/VirrageS/chirp/backend/services"
	"gopkg.in/gin-gonic/gin.v1"
	"strconv"
)

func GetTweets(context *gin.Context) {
	// TODO: support filtering
	//expected_user_id := context.Query("author")
	//expected_user_name := context.Query("author")
	// ...

	tweets, error := services.GetTweets()

	if error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": error,
		})
		return
	}

	context.JSON(http.StatusOK, tweets)
}

func GetTweet(context *gin.Context) {
	parameterId := context.Query("id")

	tweetId, err := strconv.ParseInt(parameterId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tweet ID.",
		})
		return
	}

	responseTweet, error := services.GetTweet(tweetId)

	if error != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Tweet with given ID not found.",
		})
		return
	}

	context.JSON(http.StatusOK, responseTweet)
}

func PostTweet(context *gin.Context) {
	tweetAuthorIdString := context.PostForm("author_id")
	content := context.PostForm("content")

	tweetAuthorId, error := strconv.ParseInt(tweetAuthorIdString, 10, 64)
	if error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Author ID was not a number.",
		})
		return
	}

	requestTweet := apiModel.NewTweet{
		AuthorId: tweetAuthorId,
		Content:  content,
	}

	responseTweet, error := services.PostTweet(requestTweet)

	if error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": error,
		})
		return
	}

	context.Header("Location", fmt.Sprintf("/user/%d", responseTweet.Id))
	context.JSON(http.StatusCreated, responseTweet)
}
