package main

/*
	TODO:
	  - generate a good secret key
	    (useful: http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key,
	    	     https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/)
*/

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api"
	"github.com/VirrageS/chirp/backend/middleware"
)

func init() {
	log.SetOutput(os.Stderr) // setup logrus logging library
}

func main() {
	router := gin.Default()
	router.Use(cors.New(configureCORS()))
	router.Use(middleware.ErrorHandler())

	authorizedRoutes := router.Group("/", middleware.TokenAuthenticator)
	{
		tweets := authorizedRoutes.Group("tweets")
		tweets.GET("/", api.GetTweets)
		tweets.POST("/", api.PostTweet)
		tweets.GET("/:id", api.GetTweet)
		tweets.DELETE("/:id", api.DeleteTweet)

		homeFeed := authorizedRoutes.Group("home_feed")
		homeFeed.GET("/", api.HomeFeed)

		users := authorizedRoutes.Group("users")
		users.GET("/", api.GetUsers)
		users.GET("/:id", api.GetUser)
	}

	accounts := router.Group("")
	{
		accounts.POST("/signup", api.RegisterUser)
		accounts.POST("/login", api.LoginUser)
	}

	router.Run(":8080")
}

func configureCORS() (config cors.Config) {
	config = cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	return
}
