package converters

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	APIModel "github.com/VirrageS/chirp/backend/api/model"
	databaseModel "github.com/VirrageS/chirp/backend/database/model"
)

var expectedAPIUSer = &APIModel.User{
	ID:        1,
	Username:  "username",
	Name:      "name",
	AvatarUrl: "url",
	Following: false,
}

// mock user converter
type TestUserConverter struct{}

func (c *TestUserConverter) ConvertAPIToDatabase(user *APIModel.NewUserForm) *databaseModel.User {
	return &databaseModel.User{}
}

func (c *TestUserConverter) ConvertArrayDatabaseToAPI(databaseUsers []*databaseModel.User) []*APIModel.User {
	return make([]*APIModel.User, 0)
}

func (c *TestUserConverter) ConvertArrayDatabasePublicUserToAPI(databaseUsers []*databaseModel.PublicUser) []*APIModel.User {
	return make([]*APIModel.User, 0)
}

func (c *TestUserConverter) ConvertDatabaseToAPI(user *databaseModel.User) *APIModel.User {
	return &APIModel.User{}
}

func (c *TestUserConverter) ConvertDatabasePublicUserToAPI(user *databaseModel.PublicUser) *APIModel.User {
	return expectedAPIUSer
}

func TestConvertDatabaseTweetToAPITweet(t *testing.T) {
	// subject
	var converter = NewTweetConverter(&TestUserConverter{})

	tweetCreationTime := time.Now()

	testCases := []struct {
		DBTweet  *databaseModel.TweetWithAuthor
		APITweet *APIModel.Tweet
	}{
		{ // positive scenario (the only one?)
			DBTweet: &databaseModel.TweetWithAuthor{
				ID: 1,
				Author: &databaseModel.PublicUser{
					ID:        1,
					Username:  "username",
					Name:      "name",
					AvatarUrl: sql.NullString{String: "url", Valid: true},
				},
				Likes:     2,
				Retweets:  3,
				CreatedAt: tweetCreationTime,
				Content:   "tweet",
			},
			APITweet: &APIModel.Tweet{
				ID:        1,
				Author:    expectedAPIUSer,
				Likes:     2,
				Retweets:  3,
				CreatedAt: tweetCreationTime,
				Content:   "tweet",
				Liked:     false,
				Retweeted: false,
			},
		},
	}

	for _, testCase := range testCases {
		actualAPITweet := converter.ConvertDatabaseTweetToAPITweet(testCase.DBTweet)
		expectedAPITweet := testCase.APITweet

		if !reflect.DeepEqual(actualAPITweet, expectedAPITweet) {
			t.Errorf("Got: %v, but expected: %v", actualAPITweet, expectedAPITweet)
		}
	}
}

func TestConvertArrayOfDatabaseTweetsToArrayOfAPITweets(t *testing.T) {
	// subject
	var converter = NewTweetConverter(&TestUserConverter{})

	tweetCreationTime := time.Now()

	testCases := []struct {
		DBTweets  []*databaseModel.TweetWithAuthor
		APITweets []*APIModel.Tweet
	}{
		{ // positive case
			DBTweets: []*databaseModel.TweetWithAuthor{
				{
					ID: 1,
					Author: &databaseModel.PublicUser{
						ID:        1,
						Username:  "username",
						Name:      "name",
						AvatarUrl: sql.NullString{String: "url", Valid: true},
					},
					Likes:     2,
					Retweets:  3,
					CreatedAt: tweetCreationTime,
					Content:   "tweet",
				},
			},
			APITweets: []*APIModel.Tweet{
				{
					ID:        1,
					Author:    expectedAPIUSer,
					Likes:     2,
					Retweets:  3,
					CreatedAt: tweetCreationTime,
					Content:   "tweet",
					Liked:     false,
					Retweeted: false,
				},
			},
		},
		{ // nil case
			DBTweets:  nil,
			APITweets: make([]*APIModel.Tweet, 0),
		},
	}

	for _, testCase := range testCases {
		actualAPITweetSlice := converter.ConvertArrayOfDatabaseTweetsToArrayOfAPITweets(testCase.DBTweets)
		expectedAPITweetSlice := testCase.APITweets

		if !reflect.DeepEqual(actualAPITweetSlice, expectedAPITweetSlice) {
			t.Errorf("Got: %v, but expected: %v", actualAPITweetSlice, expectedAPITweetSlice)
		}
	}
}
