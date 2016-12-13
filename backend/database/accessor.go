package database

import (
	"time"

	"github.com/VirrageS/chirp/backend/model"
)

type TweetDataAccessor interface {
	GetTweets(requestingUserID int64) ([]*model.Tweet, error)
	GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error)
	GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error)
	InsertTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error)
	DeleteTweet(tweetID int64) error
	LikeTweet(tweetID, userID int64) error
	UnlikeTweet(tweetID, userID int64) error
}

type UserDataAccessor interface {
	GetUsers(requestingUserID int64) ([]*model.PublicUser, error)
	GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error)
	GetUserByEmail(email string) (*model.User, error)
	InsertUser(user *model.NewUserForm) (*model.PublicUser, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error
	FollowUser(followeeID, followerID int64) error
}

type DatabaseAccessor interface {
	UserDataAccessor
	TweetDataAccessor
}
