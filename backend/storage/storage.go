package storage

import (
	_ "github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/fulltextsearch"
)

// Struct that implements StorageAccessor
type Storage struct {
	UserDataAccessor
	TweetDataAccessor
}

// Constructs new Storage that uses given DAOs, cache and full text search provider
func NewStorage(userDAO database.UserDAO, followsDAO database.FollowsDAO, tweetDAO database.TweetDAO, likesDAO database.LikesDAO,
	cache cache.CacheProvider, fts fulltextsearch.Searcher) StorageAccessor {
	userStorage := NewUserStorage(userDAO, followsDAO, cache, fts)
	tweetStorage := NewTweetStorage(tweetDAO, likesDAO, cache, userStorage, fts)

	return &Storage{
		UserDataAccessor:  userStorage,
		TweetDataAccessor: tweetStorage,
	}
}
