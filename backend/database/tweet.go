package database

import (
	"errors"
	"time"
)

type Tweet struct {
	Id           int64
	AuthorId     int64
	LikeCount    int64
	RetweetCount int64
	CreatedAt    time.Time
	Content      string
}

var tweets = []Tweet{
	{
		Id:           1,
		AuthorId:     users[0].Id,
		LikeCount:    0,
		RetweetCount: 0,
		CreatedAt:    time.Unix(0, 0),
		Content:      "siema siema siema",
	},
}

func GetTweets() ([]Tweet, error) {
	return tweets, nil
}

func GetTweet(tweetId int64) (Tweet, error) {
	tweet := getTweetWithId(tweetId)

	/* Emulate DB query fail? */
	if (Tweet{}) == tweet {
		return Tweet{}, errors.New("Tweet with given ID was not found.")
	}

	return tweet, nil
}

func InsertTweet(tweet Tweet) Tweet {
	tweetId := insertTweetToDatabase(tweet)
	tweet.Id = tweetId

	return tweet
}

func getTweetWithId(tweet_id int64) Tweet {
	var found_tweet Tweet

	for _, tweet := range tweets {
		if tweet.Id == tweet_id {
			found_tweet = tweet
		}
	}

	return found_tweet
}

func insertTweetToDatabase(tweet Tweet) int64 {
	tweetId := len(tweets) + 1
	tweet.Id = int64(tweetId)

	tweets = append(tweets, tweet)

	return int64(tweetId)
}
