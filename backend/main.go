package main

/*
	TODO: 	- fix error handling
*/

import (
	//"github.com/itsjamie/gin-cors"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api"
)

func main() {
	router := gin.Default()

	//router.Use(cors.Middleware(cors.Config{
	//	Origins:        "*",
	//	Methods:        "GET, PUT, POST, DELETE",
	//	RequestHeaders: "Origin, Authorization, Content-Type",
	//}))

	tweets := router.Group("/tweets")
	{
		tweets.GET("/", api.GetTweets)
		tweets.POST("/", api.PostTweet)
		tweets.GET("/:id", api.GetTweet)
	}

	users := router.Group("/users")
	{
		users.GET("/", api.GetUsers)
		users.POST("/", api.PostUser)
		users.GET("/:id", api.GetUser)
	}

	router.Run(":8080")
}
