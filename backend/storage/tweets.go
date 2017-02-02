package storage

import (
	"github.com/VirrageS/chirp/backend/async"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
	"github.com/VirrageS/chirp/backend/storage/cache"
	"github.com/VirrageS/chirp/backend/storage/database"
	"github.com/VirrageS/chirp/backend/storage/fulltextsearch"
)

// Struct that implements TweetDataAccessor using given DAO, cache and full text search provider
type tweetsStorage struct {
	tweetsDAO    database.TweetsDAO
	likesDAO     database.LikesDAO
	cache        cache.Accessor
	usersStorage usersDataAccessor
	fts          fulltextsearch.TweetsSearcher
}

// newTweetsStorage constructs tweetsStorage that uses given likesDAO, tweetsDAO, usersStorage, cache Accessor and TweetSearcher
func newTweetsStorage(tweetsDAO database.TweetsDAO, likesDAO database.LikesDAO, usersStorage usersDataAccessor, cache cache.Accessor, fts fulltextsearch.TweetsSearcher) tweetsDataAccessor {
	return &tweetsStorage{
		tweetsDAO:    tweetsDAO,
		likesDAO:     likesDAO,
		cache:        cache,
		usersStorage: usersStorage,
		fts:          fts,
	}
}

func (s *tweetsStorage) GetUsersTweets(userID, requestingUserID int64) ([]*model.Tweet, error) {
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

func (s *tweetsStorage) GetTweetsByAuthorIDs(authorsIDs []int64, requestingUserID int64) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0)

	// This should not be parallel because each GetUsersTweets is arleady
	// parallel so we do not want overload database/cache.
	for _, userID := range authorsIDs {
		usersTweets, err := s.GetUsersTweets(userID, requestingUserID)
		if err != nil {
			return nil, err
		}

		tweets = append(tweets, usersTweets...)
	}

	return tweets, nil
}

func (s *tweetsStorage) GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	var (
		tweet *model.Tweet
		err   error
	)

	key := cache.Key{"tweet", tweetID}
	if exists, _ := s.cache.GetSingle(key, tweet); !exists {
		tweet, err = s.tweetsDAO.GetTweetByID(tweetID)
		if err == errors.NoResultsError {
			return nil, errors.NoResultsError
		} else if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, tweet})
	}

	err = s.collectTweetData(tweet, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweet, nil
}

func (s *tweetsStorage) InsertTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error) {
	insertedTweet, err := s.tweetsDAO.InsertTweet(tweet)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	err = s.collectTweetData(insertedTweet, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	s.cache.Set(cache.Entry{cache.Key{"tweet", insertedTweet.ID}, insertedTweet})
	s.cache.SAdd(cache.Key{"tweets.ids", requestingUserID}, insertedTweet.ID)

	return insertedTweet, nil
}

func (s *tweetsStorage) DeleteTweet(tweetID, requestingUserID int64) error {
	err := s.tweetsDAO.DeleteTweet(tweetID)
	if err != nil {
		return errors.UnexpectedError
	}

	s.cache.Delete(cache.Key{"tweet", tweetID})
	s.cache.SRemove(cache.Key{"tweets.ids", requestingUserID}, tweetID)

	return nil
}

func (s *tweetsStorage) LikeTweet(tweetID, requestingUserID int64) error {
	liked, err := s.likesDAO.LikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	if liked {
		s.cache.Incr(cache.Key{"tweet", tweetID, "like.count"})
		s.cache.Set(
			cache.Entry{cache.Key{"tweet", tweetID, "liked.by", requestingUserID}, true},
		)
	}

	return nil
}

func (s *tweetsStorage) UnlikeTweet(tweetID, requestingUserID int64) error {
	unliked, err := s.likesDAO.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		return errors.UnexpectedError
	}

	if unliked {
		s.cache.Decr(cache.Key{"tweet", tweetID, "like.count"})
		s.cache.Set(
			cache.Entry{cache.Key{"tweet", tweetID, "liked.by", requestingUserID}, false},
		)
	}

	return nil
}

func (s *tweetsStorage) GetTweetsUsingQueryString(querystring string, requestingUserID int64) ([]*model.Tweet, error) {
	tweetsIDs := make([]int64, 0)

	key := cache.Key{"tweets", "querystring", querystring}
	if exists, _ := s.cache.GetSingle(key, &tweetsIDs); !exists {
		var err error

		tweetsIDs, err = s.fts.GetTweetsIDs(querystring)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, tweetsIDs})
	}

	return s.getTweetsByIDs(tweetsIDs, requestingUserID)
}

func (s *tweetsStorage) getTweetsIDsByAuthorID(userID int64) ([]int64, error) {
	tweetsIDs := make([]int64, 0)

	key := cache.Key{"tweets.ids", userID}
	if exists, _ := s.cache.GetSingle(key, &tweetsIDs); !exists {
		var err error

		tweetsIDs, err = s.tweetsDAO.GetTweetsIDsByAuthorID(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, tweetsIDs})
	}

	return tweetsIDs, nil
}

func (s *tweetsStorage) collectTweetsData(tweets []*model.Tweet, requestingUserID int64) error {
	pool := async.NewWorkerPool(func(task async.Task) *async.Result {
		err := s.collectTweetData(task.(*model.Tweet), requestingUserID)
		return &async.Result{nil, err}
	})
	defer pool.Close()

	for _, tweet := range tweets {
		pool.PostTask(tweet)
	}

	for range tweets {
		if result := pool.GetResult(); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// Be careful - this is function does SIDE EFFECTS only
func (s *tweetsStorage) collectTweetData(tweet *model.Tweet, requestingUserID int64) error {
	var (
		author    *model.PublicUser
		likeCount int64
		isLiked   bool
	)

	author, err := s.usersStorage.GetUserByID(tweet.Author.ID, requestingUserID)
	if err != nil {
		return err
	}

	key := cache.Key{"tweet", tweet.ID, "like.count"}
	if exists, _ := s.cache.GetSingle(key, &likeCount); !exists {
		likeCount, err = s.likesDAO.GetLikeCount(tweet.ID)
		if err != nil {
			return err
		}

		s.cache.Set(cache.Entry{key, likeCount})
	}

	key = cache.Key{"tweet", tweet.ID, "liked.by", requestingUserID}
	if exists, _ := s.cache.GetSingle(key, &isLiked); !exists {
		isLiked, err = s.likesDAO.IsLiked(tweet.ID, requestingUserID)
		if err != nil {
			return err
		}

		s.cache.Set(cache.Entry{key, isLiked})
	}

	tweet.Author = author
	tweet.LikeCount = likeCount
	tweet.Liked = isLiked

	return nil
}

func (s *tweetsStorage) getTweetsByIDs(tweetsIDs []int64, requestingUserID int64) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0, len(tweetsIDs))

	pool := async.NewWorkerPool(func(task async.Task) *async.Result {
		var tweet *model.Tweet
		id := task.(int64)

		key := cache.Key{"tweet", id}
		if exists, _ := s.cache.GetSingle(key, tweet); !exists {
			var err error

			tweet, err = s.tweetsDAO.GetTweetByID(id)
			if err != nil {
				return &async.Result{nil, err}
			}

			s.cache.Set(cache.Entry{key, tweet})
		}

		return &async.Result{tweet, nil}
	})
	defer pool.Close()

	for _, id := range tweetsIDs {
		pool.PostTask(id)
	}

	for range tweetsIDs {
		result := pool.GetResult()
		if result.Error != nil {
			return nil, errors.UnexpectedError
		}

		tweets = append(tweets, result.Value.(*model.Tweet))
	}

	// fill tweets with missing data
	err := s.collectTweetsData(tweets, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return tweets, nil
}
