package main

/*
	TODO:
	  - generate a good secret key
	    (useful: http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key,
	    	     https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/)
*/

import (
	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/server"
)

func main() {
	serverConfig, databaseConfig, redisConfig := config.GetConfig("config")

	db := database.NewConnection(databaseConfig)
	redis := cache.NewRedisCache(redisConfig)

	s := server.New(db, redis, serverConfig)
	s.Run(":8080")
}
