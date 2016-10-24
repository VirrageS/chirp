package main

/*
	TODO: 	- fix error handling
*/

import (
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"

	"github.com/VirrageS/cache"

	"github.com/VirrageS/chirp/backend/api"
)

func main() {
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})
	iris.Use(crs)

	cache := cache.NewCache(time.Minute * 2)
	iris.UseFunc(func(c *iris.Context) {
		c.Set("cache", cache)
		c.Next()
	})

	tweets := iris.Party("/tweets")
	{
		tweets.Get("/", api.GetTweets)
		tweets.Post("/", api.PostTweet)
		tweets.Get("/:id", api.GetTweet)
	}

	users := iris.Party("/users")
	{
		users.Get("/", api.GetUsers)
		users.Post("/", api.PostUser)
		users.Get("/:id", api.GetUser)
	}

	iris.Listen(":8080")
}
