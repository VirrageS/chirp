package api

import (
	"gopkg.in/gin-gonic/gin.v1"
)

type APIProvider interface {
	RegisterUser(context *gin.Context)
	LoginUser(context *gin.Context)

	GetTweets(context *gin.Context)
	GetTweet(context *gin.Context)
	PostTweet(context *gin.Context)
	DeleteTweet(context *gin.Context)
	HomeFeed(context *gin.Context)

	GetUsers(context *gin.Context)
	GetUser(context *gin.Context)
}