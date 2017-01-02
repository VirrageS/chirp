package storage

import (
	_ "github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
)

// Struct that implements StorageAcessor
type Storage struct {
	UserDataAccessor
	TweetDataAccessor
}

// Constructs new Storage that uses given DAOs and cache
func NewStorage(userDAO database.UserDAO, followsDAO database.FollowsDAO, tweetDAO database.TweetDAO, likesDAO database.LikesDAO,
	cache cache.CacheProvider) StorageAccessor {
	userStorage := NewUserStorage(userDAO, followsDAO, cache)
	tweetStorage := NewTweetStorage(tweetDAO, likesDAO, cache, userStorage)

	return &Storage{
		UserDataAccessor:  userStorage,
		TweetDataAccessor: tweetStorage,
	}
}
