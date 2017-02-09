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
	"github.com/VirrageS/chirp/backend/password"
	"github.com/VirrageS/chirp/backend/server"
	"github.com/VirrageS/chirp/backend/token"
)

func main() {
	conf := config.New()
	if conf == nil {
		panic("Failed to get config.")
	}

	db := database.NewConnection(conf.Database)
	redis := cache.NewRedisCache(conf.Redis)
	elasticsearch := fulltextsearch.NewElasticsearch(conf.Elasticsearch)
	tokenManager := token.NewManager(conf.Token)
	passwordManager := password.NewBcryptManager(conf.Password)

	s := server.New(db, redis, elasticsearch, tokenManager, passwordManager, conf.AuthorizationGoogle)
	s.Run(":8080")
}
