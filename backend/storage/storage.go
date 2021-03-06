package storage

import (
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/storage/cache"
	"github.com/VirrageS/chirp/backend/storage/database"
	"github.com/VirrageS/chirp/backend/storage/fulltextsearch"
)

type storage struct {
	usersDataAccessor
	tweetsDataAccessor
}

// New constructs Accessor that TODO
func New(postgresConfig config.PostgresConfigProvider, redisConfig config.RedisConfigProvider, elasticsearchConfig config.ElasticsearchConfigProvider) Accessor {
	db := database.NewPostgresDatabase(postgresConfig)
	if db == nil {
		panic("failed to connect to Postgres instance")
	}

	usersDAO := database.NewUserDAO(db)
	tweetsDAO := database.NewTweetDAO(db)
	followsDAO := database.NewFollowsDAO(db)
	likesDAO := database.NewLikesDAO(db)

	cache := cache.NewRedisCache(redisConfig)
	if cache == nil {
		panic("failed to connect to Redis instance")
	}

	fts := fulltextsearch.NewElasticsearchSearch(elasticsearchConfig)
	if fts == nil {
		panic("failed to connect to Elasticsearch instance")
	}

	usersStorage := newUsersStorage(usersDAO, followsDAO, cache, fts)
	tweetsStorage := newTweetsStorage(tweetsDAO, likesDAO, usersStorage, cache, fts)
	return &storage{
		usersDataAccessor:  usersStorage,
		tweetsDataAccessor: tweetsStorage,
	}
}
