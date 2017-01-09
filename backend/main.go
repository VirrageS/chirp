package main

/*
	TODO:
	  - generate a good secret key
	    (useful: http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key,
	    	     https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/)

	  - server should not start before index is created in elasticsearch,
	    see: https://github.com/VirrageS/chirp/issues/190
*/

import (
	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/fulltextsearch"
	"github.com/VirrageS/chirp/backend/server"
	"github.com/VirrageS/chirp/backend/token"
)

func main() {
	serverConfig, databaseConfig, redisConfig, authorizationGoogleConfig, elasticsearchConfig := config.GetConfig("config")

	db := database.NewConnection(databaseConfig)
	redis := cache.NewRedisCache(redisConfig)
	elasticsearch := fulltextsearch.NewElasticsearch(elasticsearchConfig)
	tokenManager := token.NewTokenManager(serverConfig)

	s := server.New(db, redis, elasticsearch, tokenManager, authorizationGoogleConfig)
	s.Run(":8080")
}
