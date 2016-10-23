package api

import (
	"strconv"
	"github.com/kataras/iris"
	"github.com/VirrageS/chirp/backend/services"
	"github.com/VirrageS/chirp/backend/apimodel"
	"fmt"
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
	parameter_id := context.Param("id")

	tweet_id, err := strconv.ParseInt(parameter_id, 10, 64)
	if err != nil {
		context.JSON(iris.StatusBadRequest, iris.Map{
			"error": "Invalid tweet ID.",
		})
		return
	}

	response_tweet, error := services.GetTweet(tweet_id)

	if error != nil {
		context.JSON(iris.StatusNotFound, iris.Map{
			"error": "Tweet with given ID not found.",
		})
		return
	}

	context.JSON(iris.StatusOK, response_tweet)
}

func PostTweet(context *iris.Context) {
	tweet_author_id_string := context.PostValue("author_id")
	content := context.PostValue("content")

	tweet_author_id, error := strconv.ParseInt(tweet_author_id_string, 10, 64)
	if error != nil {
		context.JSON(iris.StatusBadRequest, iris.Map{
			"error": "Invalid author ID.",
		})
		return
	}

	request_tweet := apimodel.NewTweet{
		AuthorId: tweet_author_id,
		Content: content,
	}

	response_tweet, error := services.PostTweet(request_tweet)

	if error != nil {
		context.JSON(iris.StatusInternalServerError, iris.Map{
			"error": error,
		})
		return
	}

	context.SetHeader("Location", fmt.Sprintf("/user/%d", response_tweet.Id))
	context.JSON(iris.StatusCreated, response_tweet)
}