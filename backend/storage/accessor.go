package storage

import (
	"time"

	"github.com/VirrageS/chirp/backend/model"
)

type TweetDataAccessor interface {
	GetUsersTweets(userID, requestingUserID int64) ([]*model.Tweet, error)
	GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error)
	InsertTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error)
	DeleteTweet(tweetID, requestingUserID int64) error
	LikeTweet(tweetID, userID int64) error
	UnlikeTweet(tweetID, userID int64) error
	GetTweetsUsingQueryString(querystring string, requestingUserID int64) ([]*model.Tweet, error)
}

type UserDataAccessor interface {
	GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error)
	GetUserByEmail(email string) (*model.User, error)
	InsertUser(user *model.NewUserForm) (*model.PublicUser, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error
	FollowUser(followeeID, followerID int64) error
	UnfollowUser(followeeID, followerID int64) error
	GetFollowers(userID, requestingUserID int64) ([]*model.PublicUser, error)
	GetFollowees(userID, requestingUserID int64) ([]*model.PublicUser, error)
	GetUsersUsingQueryString(querystring string, requestingUserID int64) ([]*model.PublicUser, error)
}

type StorageAccessor interface {
	UserDataAccessor
	TweetDataAccessor
}
