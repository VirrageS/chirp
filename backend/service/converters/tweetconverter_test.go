package converters

import (
	"database/sql"
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

		if actualAPITweet.ID != expectedAPITweet.ID ||
			*actualAPITweet.Author != *expectedAPITweet.Author ||
			actualAPITweet.Likes != expectedAPITweet.Likes ||
			actualAPITweet.Retweets != expectedAPITweet.Retweets ||
			actualAPITweet.CreatedAt != expectedAPITweet.CreatedAt ||
			actualAPITweet.Content != expectedAPITweet.Content ||
			actualAPITweet.Liked != expectedAPITweet.Liked ||
			actualAPITweet.Retweeted != expectedAPITweet.Retweeted {
			t.Errorf("Got: %v, but expected: %v", actualAPITweet, expectedAPITweet)
		}
	}
}
