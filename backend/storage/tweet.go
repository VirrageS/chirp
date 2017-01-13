package storage

import (
	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/fulltextsearch"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
)

// Struct that implements TweetDataAccessor using given DAO, cache and full text search provider
type TweetStorage struct {
	tweetDAO    database.TweetDAO
	likesDAO    database.LikesDAO
	cache       cache.CacheProvider
	userStorage UserDataAccessor
	fts         fulltextsearch.TweetSearcher
}

// Constructs TweetStorage that uses given likesDAO, tweetDAO, CacheProvider, UserStorage and TweetSearcher
func NewTweetStorage(tweetDAO database.TweetDAO, likesDAO database.LikesDAO,
	cache cache.CacheProvider, userStorage UserDataAccessor, fts fulltextsearch.TweetSearcher) *TweetStorage {

	return &TweetStorage{
		tweetDAO:    tweetDAO,
		likesDAO:    likesDAO,
		cache:       cache,
		userStorage: userStorage,
		fts:         fts,
	}
}

func (s *TweetStorage) GetUsersTweets(userID, requestingUserID int64) ([]*model.Tweet, error) {
	tweetsIDs, err := s.getTweetsIDsByAuthorID(userID)
	if err != nil {
		return nil, err
	}

	tweets, err := s.getTweetsByIDs(tweetsIDs, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweets, nil
}

func (s *TweetStorage) GetTweetsByAuthorIDs(authorsIDs []int64, requestingUserID int64) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0)

	for _, userID := range authorsIDs {
		usersTweets, err := s.GetUsersTweets(userID, requestingUserID)
		if err != nil {
			return nil, err
		}

		tweets = append(tweets, usersTweets...)
	}

	return tweets, nil
}

func (s *TweetStorage) GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	var tweet *model.Tweet

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweet", tweetID}, &tweet); !exists {
		var err error

		tweet, err = s.tweetDAO.GetTweetByID(tweetID)
		if err == errors.NoResultsError {
			return nil, err
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

	s.cache.SetWithFields(cache.Fields{"tweet", insertedTweet.ID}, insertedTweet)
	s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"tweet", insertedTweet.ID, "likedBy", requestingUserID}, false)
	s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"tweet", insertedTweet.ID, "likeCount"}, 0)
	s.cache.DeleteWithFields(cache.Fields{"tweetsIDs", requestingUserID})

	return insertedTweet, nil
}

func (s *TweetStorage) DeleteTweet(tweetID, requestingUserID int64) error {
	err := s.tweetDAO.DeleteTweet(tweetID)
	if err != nil {
		return errors.UnexpectedError
	}

	s.cache.DeleteWithFields(cache.Fields{"tweet", tweetID})
	// for now delete tweets affects only 'tweets of user with ID' of the author of the tweet
	s.cache.DeleteWithFields(cache.Fields{"tweetsIDs", requestingUserID})

	return nil
}

func (s *TweetStorage) LikeTweet(tweetID, requestingUserID int64) error {
	liked, err := s.likesDAO.LikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	if liked {
		s.cache.IncrementWithFields(cache.Fields{"tweet", tweetID, "likeCount"})
	}
	s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"tweet", tweetID, "likedBy", requestingUserID}, true)

	return nil
}

func (s *TweetStorage) UnlikeTweet(tweetID, requestingUserID int64) error {
	unliked, err := s.likesDAO.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	if unliked {
		s.cache.DecrementWithFields(cache.Fields{"tweet", tweetID, "likeCount"})
	}
	s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"tweet", tweetID, "likedBy", requestingUserID}, false)

	return nil
}

func (s *TweetStorage) GetTweetsUsingQueryString(querystring string, requestingUserID int64) ([]*model.Tweet, error) {
	tweetsIDs := make([]int64, 0)

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweets", "querystring", querystring}, &tweetsIDs); !exists {
		var err error

		tweetsIDs, err = s.fts.GetTweetsIDs(querystring)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"tweets", "querystring", querystring}, tweetsIDs)
	}

	return s.getTweetsByIDs(tweetsIDs, requestingUserID)
}

func (s *TweetStorage) getTweetsIDsByAuthorID(userID int64) ([]int64, error) {
	tweetsIDs := make([]int64, 0)

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweetsIDs", userID}, &tweetsIDs); !exists {
		var err error

		tweetsIDs, err = s.tweetDAO.GetTweetsIDsByAuthorID(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		s.cache.SetWithFields(cache.Fields{"tweetsIDs", userID}, tweetsIDs)
	}

	return tweetsIDs, nil
}

// on deadlocks blame VirrageS
func (s *TweetStorage) collectTweetsData(tweets []*model.Tweet, requestingUserID int64) error {
	errChan := make(chan error, len(tweets))

	for _, tweet := range tweets {
		go func(tweet *model.Tweet) {
			err := s.collectTweetData(tweet, requestingUserID)
			errChan <- err
		}(tweet)
	}

	for range tweets {
		if err := <-errChan; err != nil {
			return err
		}
	}

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

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweet", tweet.ID, "likeCount"}, &likeCount); !exists {
		likeCount, err = s.likesDAO.GetLikeCount(tweet.ID)
		if err != nil {
			return err
		}

		s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"tweet", tweet.ID, "likeCount"}, likeCount)
	}

	if exists, _ := s.cache.GetWithFields(cache.Fields{"tweet", tweet.ID, "likedBy", requestingUserID}, &isLiked); !exists {
		isLiked, err = s.likesDAO.IsLiked(tweet.ID, requestingUserID)
		if err != nil {
			return err
		}

		s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"tweet", tweet.ID, "likedBy", requestingUserID}, isLiked)
	}

	tweet.Author = author
	tweet.LikeCount = likeCount
	tweet.Liked = isLiked

	return nil
}

func (s *TweetStorage) getTweetsByIDs(tweetsIDs []int64, requestingUserID int64) ([]*model.Tweet, error) {
	expectedTweetCount := len(tweetsIDs)

	tweets := make([]*model.Tweet, 0, expectedTweetCount)
	notInCacheCount := 0

	// get tweets from cache
	for _, id := range tweetsIDs {
		var tweet model.Tweet

		if exists, _ := s.cache.GetWithFields(cache.Fields{"tweet", id}, &tweet); exists {
			tweets = append(tweets, &tweet)
		} else {
			tweetsIDs[notInCacheCount] = id
			notInCacheCount++
		}
	}

	// get tweets that are not in cache from database
	if notInCacheCount > 0 {
		dbTweets, err := s.tweetDAO.GetTweetsByIDs(tweetsIDs[:notInCacheCount])
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, dbTweets...)
	}

	if len(tweets) != expectedTweetCount {
		log.WithFields(log.Fields{
			"number of tweets found":   len(tweets),
			"expected tweets of users": expectedTweetCount,
		}).Error("Found less tweets than expected in getTweetsByIDs.")
	}

	// fill tweets with missing data
	err := s.collectTweetsData(tweets, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweets, nil
}
