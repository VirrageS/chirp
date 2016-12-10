package database

import (
	"time"

	"github.com/VirrageS/chirp/backend/model"
)

type TweetDataAccessor interface {
	GetTweets() ([]*model.Tweet, error)
	GetTweetsOfUserWithID(userID int64) ([]*model.Tweet, error)
	GetTweet(tweetID int64) (*model.Tweet, error)
	InsertTweet(tweet *model.NewTweet) (*model.Tweet, error)
	DeleteTweet(tweetID int64) error
	LikeTweet(tweetID, userID int64) error
}

type UserDataAccessor interface {
	GetUsers() ([]*model.PublicUser, error)
	GetUserByID(userID int64) (*model.PublicUser, error)
	GetUserByEmail(email string) (*model.User, error)
	InsertUser(user *model.NewUserForm) (*model.PublicUser, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error
}

type DatabaseAccessor interface {
	UserDataAccessor
	TweetDataAccessor
}
