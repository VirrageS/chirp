package main

/*
	TODO:
	  - use c.bind() feature of Gin
	  - add logging (probably using https://github.com/golang/glog)
	  - generate a good secret key and implement reading it from config file
	    (useful: http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key,
	    	     https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/)
*/

import (
	//"github.com/itsjamie/gin-cors"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api"
	"github.com/VirrageS/chirp/backend/middleware"
)

func main() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	//router.Use(cors.Middleware(cors.Config{
	//	Origins:        "*",
	//	Methods:        "GET, PUT, POST, DELETE",
	//	RequestHeaders: "Origin, Authorization, Content-Type",
	//}))

	tweets := router.Group("/tweets")
	{
		tweets.GET("/", api.GetTweets)
		tweets.POST("/", middleware.TokenAuthenticator(), api.PostTweet)
		tweets.GET("/:id", api.GetTweet)
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
