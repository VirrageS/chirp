package database

import (
	"time"

	"github.com/VirrageS/chirp/backend/database/model"
)

type TweetDataAccessor interface {
	GetTweets() ([]*model.TweetWithAuthor, error)
	GetTweetsOfUserWithID(userID int64) ([]*model.TweetWithAuthor, error)
	GetTweet(tweetID int64) (*model.TweetWithAuthor, error)
	InsertTweet(tweet model.Tweet) (*model.TweetWithAuthor, error)
	DeleteTweet(tweetID int64) error
}

type UserDataAccessor interface {
	GetUsers() ([]*model.PublicUser, error)
	GetUserByID(userID int64) (*model.PublicUser, error)
	GetUserByEmail(email string) (*model.User, error)
	InsertUser(user *model.User) (*model.User, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error
}

type DatabaseAccessor interface {
	UserDataAccessor
	TweetDataAccessor
}
