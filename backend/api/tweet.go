package api

import (
	"fmt"
	"github.com/VirrageS/chirp/backend/apiModel"
	"github.com/VirrageS/chirp/backend/services"
	"github.com/kataras/iris"
	"strconv"
)

func GetTweets(context *iris.Context) {
	// TODO: support filtering
	//expected_user_id := context.Param("author")
	//expected_user_name := context.Param("author")
	// ...

	tweets, error := services.GetTweets()

	if error != nil {
		context.JSON(iris.StatusInternalServerError, iris.Map{
			"error": error,
		})
		return
	}

	context.JSON(iris.StatusOK, tweets)
}

func GetTweet(context *iris.Context) {
	parameterId := context.Param("id")

	tweetId, err := strconv.ParseInt(parameterId, 10, 64)
	if err != nil {
		context.JSON(iris.StatusBadRequest, iris.Map{
			"error": "Invalid tweet ID.",
		})
		return
	}

	responseTweet, error := services.GetTweet(tweetId)

	if error != nil {
		context.JSON(iris.StatusNotFound, iris.Map{
			"error": "Tweet with given ID not found.",
		})
		return
	}

	context.JSON(iris.StatusOK, responseTweet)
}

func PostTweet(context *iris.Context) {
	tweetAuthorIdString := context.PostValue("author_id")
	content := context.PostValue("content")

	tweetAuthorId, error := strconv.ParseInt(tweetAuthorIdString, 10, 64)
	if error != nil {
		context.JSON(iris.StatusBadRequest, iris.Map{
			"error": "Invalid author ID.",
		})
		return
	}

	requestTweet := apiModel.NewTweet{
		AuthorId: tweetAuthorId,
		Content:  content,
	}

	responseTweet, error := services.PostTweet(requestTweet)

	if error != nil {
		context.JSON(iris.StatusInternalServerError, iris.Map{
			"error": error,
		})
		return
	}

	context.SetHeader("Location", fmt.Sprintf("/user/%d", responseTweet.Id))
	context.JSON(iris.StatusCreated, responseTweet)
}
