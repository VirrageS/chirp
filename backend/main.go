package main

/*
	TODO:
	  - use c.bind() feature of Gin
	  - add logging (probably using https://github.com/golang/glog)
	  - generate a good secret key
	    (useful: http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key,
	    	     https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/)
*/

import (
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api"
	"github.com/VirrageS/chirp/backend/middleware"
)

func main() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	tweets := router.Group("/tweets")
	{
		tweets.GET("/", api.GetTweets)
		tweets.POST("/", middleware.TokenAuthenticator, api.PostTweet)
		tweets.GET("/:id", api.GetTweet)
		tweets.DELETE("/:id", middleware.TokenAuthenticator, api.DeleteTweet)
	}

	homeFeed := router.Group("/home_feed")
	{
		homeFeed.GET("/", middleware.TokenAuthenticator, api.HomeFeed)
	}

	users := router.Group("/users")
	{
		users.GET("/", api.GetUsers)
		users.GET("/:id", api.GetUser)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/register", api.RegisterUser)
		auth.POST("/login", api.LoginUser)
	}

	router.Run(":8080")
}
