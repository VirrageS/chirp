package storage

import (
	"database/sql"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
)

// Struct that implements TweetDataAccessor using sql (postgres) database
type TweetStorage struct {
	DAO   database.TweetDAO
	cache cache.CacheProvider
}

// Constructs TweetDB that uses a given sql.DB connection and CacheProvider
func NewTweetStorage(DAO database.TweetDAO, cache cache.CacheProvider) *TweetStorage {
	return &TweetStorage{
		DAO,
		cache,
	}
}

func (db *TweetStorage) GetTweets(requestingUserID int64) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0)
	if exists, _ := db.cache.GetWithFields(cache.Fields{"tweets", requestingUserID}, &tweets); exists {
		return tweets, nil
	}

	tweets, err := db.DAO.GetTweets(requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	db.cache.SetWithFields(cache.Fields{"tweets", requestingUserID}, tweets)
	return tweets, nil
}

func (db *TweetStorage) GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0)
	if exists, _ := db.cache.GetWithFields(cache.Fields{"tweets", userID, requestingUserID}, &tweets); exists {
		return tweets, nil
	}

	tweets, err := db.DAO.GetTweetsOfUserWithID(userID, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	db.cache.SetWithFields(cache.Fields{"tweets", userID, requestingUserID}, tweets)
	return tweets, nil
}

func (db *TweetStorage) GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	var tweet *model.Tweet
	if exists, _ := db.cache.GetWithFields(cache.Fields{"tweet", tweetID, requestingUserID}, tweet); exists {
		return tweet, nil
	}

	tweet, err := db.DAO.GetTweetUsingQuery(`
		SELECT tweets.id, tweets.created_at, tweets.content,
		 	users.id, users.username, users.name, users.avatar_url,
		 	COUNT(likes.tweet_id) as likes,
		 	SUM(case when likes.user_id=$1 then 1 else 0 end) > 0 as liked
		FROM tweets
			JOIN users ON tweets.author_id = users.id AND tweets.id=$2
			LEFT JOIN likes ON tweets.id = likes.tweet_id
		GROUP BY tweets.id, users.id
		ORDER BY tweets.created_at DESC;`,
		requestingUserID, tweetID)

	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}

	if err != nil {
		return nil, errors.UnexpectedError
	}

	db.cache.SetWithFields(cache.Fields{"tweet", tweetID, requestingUserID}, tweet)
	return tweet, nil
}

func (db *TweetStorage) InsertTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error) {
	tweetID, err := db.DAO.InsertTweetToDatabase(tweet)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	// TODO: this is probably super ugly. Maybe fetch user only?
	// Probably could just fetch user from cache
	newTweet, err := db.DAO.GetTweetUsingQuery(`
		SELECT tweets.id, tweets.created_at, tweets.content,
		 	users.id, users.username, users.name, users.avatar_url,
		 	COUNT(likes.tweet_id) as likes,
		 	SUM(case when likes.user_id=$1 then 1 else 0 end) > 0 as liked
		FROM tweets
			JOIN users ON tweets.author_id = users.id AND tweets.id=$2
			LEFT JOIN likes ON tweets.id = likes.tweet_id
		GROUP BY tweets.id, users.id
		ORDER BY tweets.created_at DESC;`,
		requestingUserID, tweetID)

	if err != nil {
		return nil, errors.UnexpectedError
	}

	// We don't flush cache on purpose. The data in cache can be not precise for some time.
	db.cache.SetWithFields(cache.Fields{"tweet", tweetID, requestingUserID}, tweet)

	return newTweet, nil
}

func (db *TweetStorage) DeleteTweet(tweetID int64) error {
	err := db.DAO.DeleteTweetWithID(tweetID)
	if err != nil {
		return errors.UnexpectedError
	}

	// Its better to just flush the cache here, because almost everything changes.
	db.cache.Flush()

	return nil
}

func (db *TweetStorage) LikeTweet(tweetID, requestingUserID int64) error {
	err := db.DAO.LikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	// TODO: Maybe a smarter way: don't delete, but just update cache with likeCount++ and liked=true,
	// Just delete from cache for the requesting user, it will be fetched back in next GET query
	db.cache.DeleteWithFields(cache.Fields{"tweet", tweetID, requestingUserID})

	return nil
}

func (db *TweetStorage) UnlikeTweet(tweetID, requestingUserID int64) error {
	err := db.DAO.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	// TODO: Maybe a smarter way: don't delete, but just update cache with likeCount-- and liked=false
	// Just delete from cache for the requesting user, it will be fetched back in next GET query
	db.cache.DeleteWithFields(cache.Fields{"tweet", tweetID, requestingUserID})

	return nil
}
