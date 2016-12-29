package storage

import (
	_ "github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
)

// Struct that implements DatabaseAccessor
type Storage struct {
	UserDataAccessor
	TweetDataAccessor
}

// Constructs new Database that uses given sql.DB connection
func NewStorage(userDAO database.UserDAO, tweetDAO database.TweetDAO, cache cache.CacheProvider) DatabaseAccessor {
	return &Storage{
		NewUserStorage(userDAO, cache),
		NewTweetStorage(tweetDAO, cache),
	}
}
