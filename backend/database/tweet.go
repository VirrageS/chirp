package database

import (
	"errors"
	"time"

	"github.com/VirrageS/chirp/backend/database/model"
)

var tweets = []model.Tweet{
	{
		ID:           1,
		AuthorID:     users[0].ID,
		LikeCount:    0,
		RetweetCount: 0,
		CreatedAt:    time.Unix(0, 0),
		Content:      "siema siema siema",
	},
}

func GetTweets() ([]model.Tweet, error) {
	return tweets, nil
}

func GetTweet(tweetID int64) (model.Tweet, error) {
	tweet, err := getTweetWithId(tweetID)

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

func getTweetWithId(tweetID int64) (model.Tweet, error) {
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
