package api

import "gopkg.in/gin-gonic/gin.v1"

type APIProvider interface {
	RegisterUser(context *gin.Context)
	LoginUser(context *gin.Context)
	RefreshAuthToken(context *gin.Context)
	GetGoogleAuthorizationURL(context *gin.Context)
	CreateOrLoginUserWithGoogle(context *gin.Context)

	GetTweet(context *gin.Context)
	PostTweet(context *gin.Context)
	DeleteTweet(context *gin.Context)
	LikeTweet(context *gin.Context)
	UnlikeTweet(context *gin.Context)
	Feed(context *gin.Context)

	GetUser(context *gin.Context)
	FollowUser(context *gin.Context)
	UnfollowUser(context *gin.Context)
	UserFollowers(context *gin.Context)
	UserFollowees(context *gin.Context)
	UserTweets(context *gin.Context)

	Search(context *gin.Context)
}
