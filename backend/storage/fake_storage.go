package storage

import (
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/storage/cache"
	"github.com/VirrageS/chirp/backend/storage/database"
	"github.com/VirrageS/chirp/backend/storage/fulltextsearch"
)

// FakeStorage exports all fields which can be necessary for tests.
type FakeStorage struct {
	Database *database.Connection
	Cache    cache.Accessor
	Storage  Accessor
}

// NewFakeStorage creates new fake storage necessary for tests.
func NewFakeStorage(postgresConfig config.PostgresConfigProvider) *FakeStorage {
	db := database.NewPostgresDatabase(postgresConfig)
	if db == nil {
		panic("failed to connect to Postgres instance")
	}

	usersDAO := database.NewUserDAO(db)
	tweetsDAO := database.NewTweetDAO(db)
	followsDAO := database.NewFollowsDAO(db)
	likesDAO := database.NewLikesDAO(db)

	cache := cache.NewFakeCache() // TODO this shoud be redis...
	fts := fulltextsearch.NewFakeSearch()

	usersStorage := newUsersStorage(usersDAO, followsDAO, cache, fts)
	tweetsStorage := newTweetsStorage(tweetsDAO, likesDAO, usersStorage, cache, fts)
	return &FakeStorage{
		Database: db,
		Cache:    cache,
		Storage: &storage{
			usersDataAccessor:  usersStorage,
			tweetsDataAccessor: tweetsStorage,
		},
	}
}
