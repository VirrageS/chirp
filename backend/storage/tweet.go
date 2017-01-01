package storage

import (
	"database/sql"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
)

// Struct that implements TweetDataAccessor using given DAO and cache
type TweetStorage struct {
	tweetDAO    database.TweetDAO
	likesDAO    database.LikesDAO
	cache       cache.CacheProvider
	userStorage UserDataAccessor
}

// Constructs TweetStorage that uses given likesDAO, tweetDAO, CacheProvider and UserStorage
func NewTweetStorage(tweetDAO database.TweetDAO, likesDAO database.LikesDAO,
	cache cache.CacheProvider, userStorage UserDataAccessor) *TweetStorage {

	return &TweetStorage{
		tweetDAO:    tweetDAO,
		likesDAO:    likesDAO,
		cache:       cache,
		userStorage: userStorage,
	}
}

func (s *TweetStorage) GetUsersTweets(userID, requestingUserID int64) ([]*model.Tweet, error) {
	tweetsIDs := make([]int64, 0)

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweetsIDs", userID}, &tweetsIDs); !exists {
		var err error

		tweetsIDs, err = s.tweetDAO.GetTweetsIDsOfUserWithID(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		s.cache.SetWithFields(cache.Fields{"tweetsIDs", userID}, tweetsIDs)
	}

	tweets, err := s.getTweetsByIDs(tweetsIDs, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweets, nil
}

func (s *TweetStorage) GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	var tweet *model.Tweet

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweet", tweetID}, tweet); !exists {
		var err error

		tweet, err = s.tweetDAO.GetTweetByID(tweetID)
		if err == sql.ErrNoRows {
			return nil, errors.NoResultsError
		}
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"tweet", tweetID}, tweet)
	}

	err := s.collectTweetData(tweet, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweet, nil
}

func (s *TweetStorage) InsertTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error) {
	insertedTweet, err := s.tweetDAO.InsertTweet(tweet)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	author, err := s.userStorage.GetUserByID(insertedTweet.Author.ID, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}
	insertedTweet.Author = author
	// we don't need to fetch more data for the tweet, since we know that it has 0 likes, and is not liked

	s.cache.SetWithFields(cache.Fields{"tweet", insertedTweet.ID}, tweet)
	s.cache.SetWithFields(cache.Fields{"tweet", insertedTweet.ID, "liked_by", requestingUserID}, false)
	s.cache.SetWithFields(cache.Fields{"tweet", insertedTweet.ID, "like_count"}, 0)
	s.cache.DeleteWithFields(cache.Fields{"tweets", requestingUserID})

	return insertedTweet, nil
}

func (s *TweetStorage) DeleteTweet(tweetID, requestingUserID int64) error {
	err := s.tweetDAO.DeleteTweet(tweetID)
	if err != nil {
		return errors.UnexpectedError
	}

	s.cache.DeleteWithFields(cache.Fields{"tweet", tweetID})
	// for now delete tweets affects only 'tweets of user with ID' of the author of the tweet
	s.cache.DeleteWithFields(cache.Fields{"tweets", requestingUserID})

	return nil
}

func (s *TweetStorage) LikeTweet(tweetID, requestingUserID int64) error {
	err := s.likesDAO.LikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	s.cache.SetWithFields(cache.Fields{"tweet", tweetID, "is_liked_by", requestingUserID}, true)
	s.cache.IncrementWithFields(cache.Fields{"tweet", tweetID, "like_count"})

	return nil
}

func (s *TweetStorage) UnlikeTweet(tweetID, requestingUserID int64) error {
	err := s.likesDAO.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	s.cache.SetWithFields(cache.Fields{"tweet", tweetID, "is_liked_by", requestingUserID}, false)
	s.cache.DecrementWithFields(cache.Fields{"tweet", tweetID, "like_count"})

	return nil
}

// Be careful - this is function does SIDE EFFECTS only
func (s *TweetStorage) collectTweetData(tweet *model.Tweet, requestingUserID int64) error {
	var author *model.PublicUser
	var likeCount int64
	var isLiked bool

	// here userStorage will take care of everything
	author, err := s.userStorage.GetUserByID(tweet.Author.ID, requestingUserID)
	if err != nil {
		return err
	}

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweet", tweet.ID, "like_count"}, &likeCount); !exists {
		likeCount, err = s.likesDAO.LikeCount(tweet.ID)
		if err != nil {
			return err
		}

		s.cache.SetWithFields(cache.Fields{"tweet", tweet.ID, "like_count"}, likeCount)
	}

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweet", tweet.ID, "is_liked_by", requestingUserID}, &isLiked); !exists {
		isLiked, err = s.likesDAO.IsLiked(tweet.ID, requestingUserID)
		if err != nil {
			return err
		}

		s.cache.SetWithFields(cache.Fields{"tweet", tweet.ID, "is_liked_by", requestingUserID}, isLiked)
	}

	tweet.Author = author
	tweet.LikeCount = likeCount
	tweet.Liked = isLiked

	return nil
}

func (s *TweetStorage) getTweetsByIDs(tweetsIDs []int64, requestingUserID int64) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0)

	// get tweets from cache
	for i, id := range tweetsIDs {
		var tweet model.Tweet

		if exists, _ := s.cache.GetWithFields(cache.Fields{"tweet", id}, &tweet); exists {
			tweets = append(tweets, &tweet)

			// remove ID from tweetsIDs
			tweetsIDs[i] = tweetsIDs[len(tweetsIDs)-1]
			tweetsIDs = tweetsIDs[:len(tweetsIDs)-1]
		}
	}

	// get tweets that are not in cache from database
	if len(tweetsIDs) > 0 {
		dbTweets, err := s.tweetDAO.GetTweetsFromListOfIDs(tweetsIDs)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, dbTweets...)
	}

	// fill tweets with missing data
	for _, tweet := range tweets {
		err := s.collectTweetData(tweet, requestingUserID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
	}

	return tweets, nil
}
