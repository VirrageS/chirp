package main

import (
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"

	"github.com/VirrageS/cache"
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

	iris.Listen(":8000")
}
