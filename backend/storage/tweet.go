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
	tweetDAO    database.TweetDAO
	likesDAO    database.LikesDAO
	cache       cache.CacheProvider
	userStorage UserDataAccessor
}

// Constructs TweetDB that uses a given sql.DB connection and CacheProvider
func NewTweetStorage(tweetDAO database.TweetDAO, likesDAO database.LikesDAO,
	cache cache.CacheProvider, userStorage UserDataAccessor) *TweetStorage {

	return &TweetStorage{
		tweetDAO:    tweetDAO,
		likesDAO:    likesDAO,
		cache:       cache,
		userStorage: userStorage,
	}
}

func (db *TweetStorage) GetTweets(requestingUserID int64) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0)
	if exists, _ := db.cache.GetWithFields(cache.Fields{"tweets", requestingUserID}, &tweets); exists {
		return tweets, nil
	}

	tweets, err := db.tweetDAO.GetTweets()
	if err != nil {
		return nil, errors.UnexpectedError
	}

	for _, tweet := range tweets {
		author, err := db.userStorage.GetUserByID(tweet.Author.ID, requestingUserID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		likeCount, err := db.likesDAO.LikeCount(tweet.ID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		isLiked, err := db.likesDAO.IsLiked(tweet.ID, requestingUserID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		tweet.Author = author
		tweet.LikeCount = likeCount
		tweet.Liked = isLiked
	}

	db.cache.SetWithFields(cache.Fields{"tweets", requestingUserID}, tweets)
	return tweets, nil
}

func (db *TweetStorage) GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0)
	if exists, _ := db.cache.GetWithFields(cache.Fields{"tweets", userID, requestingUserID}, &tweets); exists {
		return tweets, nil
	}

	tweets, err := db.tweetDAO.GetTweetsOfUserWithID(userID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	for _, tweet := range tweets {
		author, err := db.userStorage.GetUserByID(tweet.Author.ID, requestingUserID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		likeCount, err := db.likesDAO.LikeCount(tweet.ID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		isLiked, err := db.likesDAO.IsLiked(tweet.ID, requestingUserID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		tweet.Author = author
		tweet.LikeCount = likeCount
		tweet.Liked = isLiked
	}

	db.cache.SetWithFields(cache.Fields{"tweets", userID, requestingUserID}, tweets)
	return tweets, nil
}

func (db *TweetStorage) GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	var tweet *model.Tweet
	if exists, _ := db.cache.GetWithFields(cache.Fields{"tweet", tweetID, requestingUserID}, tweet); exists {
		return tweet, nil
	}

	tweet, err := db.tweetDAO.GetTweetWithID(tweetID)
	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}

	if err != nil {
		return nil, errors.UnexpectedError
	}

	author, err := db.userStorage.GetUserByID(tweet.Author.ID, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}
	likeCount, err := db.likesDAO.LikeCount(tweet.ID)
	if err != nil {
		return nil, errors.UnexpectedError
	}
	isLiked, err := db.likesDAO.IsLiked(tweet.ID, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	tweet.Author = author
	tweet.LikeCount = likeCount
	tweet.Liked = isLiked

	db.cache.SetWithFields(cache.Fields{"tweet", tweetID, requestingUserID}, tweet)
	return tweet, nil
}

func (db *TweetStorage) InsertTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error) {
	insertedTweet, err := db.tweetDAO.InsertTweet(tweet)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	author, err := db.userStorage.GetUserByID(insertedTweet.Author.ID, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}
	likeCount, err := db.likesDAO.LikeCount(insertedTweet.ID)
	if err != nil {
		return nil, errors.UnexpectedError
	}
	isLiked, err := db.likesDAO.IsLiked(insertedTweet.ID, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	insertedTweet.Author = author
	insertedTweet.LikeCount = likeCount
	insertedTweet.Liked = isLiked

	// We don't flush cache on purpose. The data in cache can be not precise for some time.
	db.cache.SetWithFields(cache.Fields{"tweet", insertedTweet.ID, requestingUserID}, tweet)

	return insertedTweet, nil
}

func (db *TweetStorage) DeleteTweet(tweetID int64) error {
	err := db.tweetDAO.DeleteTweet(tweetID)
	if err != nil {
		return errors.UnexpectedError
	}

	// Its better to just flush the cache here, because almost everything changes.
	db.cache.Flush()

	return nil
}

func (db *TweetStorage) LikeTweet(tweetID, requestingUserID int64) error {
	err := db.likesDAO.LikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	// TODO: Maybe a smarter way: don't delete, but just update cache with likeCount++ and liked=true,
	// Just delete from cache for the requesting user, it will be fetched back in next GET query
	db.cache.DeleteWithFields(cache.Fields{"tweet", tweetID, requestingUserID})

	return nil
}

func (db *TweetStorage) UnlikeTweet(tweetID, requestingUserID int64) error {
	err := db.likesDAO.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	// TODO: Maybe a smarter way: don't delete, but just update cache with likeCount-- and liked=false
	// Just delete from cache for the requesting user, it will be fetched back in next GET query
	db.cache.DeleteWithFields(cache.Fields{"tweet", tweetID, requestingUserID})

	return nil
}
