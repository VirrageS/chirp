package database

import (
	"errors"
	"time"

	"github.com/VirrageS/chirp/backend/database/model"
)

var tweets = []model.Tweet{
	{
		ID:        1,
		AuthorID:  users[0].ID,
		Likes:     0,
		Retweets:  0,
		CreatedAt: time.Unix(0, 0),
		Content:   "siema siema siema",
	},
}

func GetTweets() ([]model.Tweet, error) {
	return tweets, nil
}

func GetTweetsOfUserWithID(userID int64) ([]model.Tweet, error) {
	var usersTweets []model.Tweet

	for _, tweet := range tweets {
		if tweet.AuthorID == userID {
			usersTweets = append(usersTweets, tweet)
		}
	}

	return usersTweets, nil
}

func GetTweet(tweetID int64) (model.Tweet, error) {
	tweet, err := getTweetWithID(tweetID)
	if err != nil {
		return model.Tweet{}, errors.New("")
	}

	return tweet, nil
}

func InsertTweet(tweet model.Tweet) (model.Tweet, error) {
	tweetID := insertTweetToDatabase(tweet)
	tweet.ID = tweetID

	return tweet, nil
}

func DeleteTweet(tweetID int64) error {
	err := deleteTweetWithID(tweetID)
	if err != nil {
		return errors.New("")
	}

	return nil
}

/* Functions that mock databse queries */

func getTweetWithID(tweetID int64) (model.Tweet, error) {
	for _, tweet := range tweets {
		if tweet.ID == tweetID {
			return tweet, nil
		}
	}

	return model.Tweet{}, errors.New("Tweet with given ID was not found.")
}

func insertTweetToDatabase(tweet model.Tweet) int64 {
	tweetID := len(tweets) + 1
	tweet.ID = int64(tweetID)

	tweets = append(tweets, tweet)

	return int64(tweetID)
}

func deleteTweetWithID(tweetID int64) error {
	for i, tweet := range tweets {
		if tweet.ID == tweetID {
			// remove tweet from the slice
			tweets[i] = tweets[len(tweets)-1] // Replace with the last one
			tweets = tweets[:len(tweets)-1]   // Chop off the last one

			return nil
		}
	}

	return errors.New("Tweet with given ID was not found.")
}
