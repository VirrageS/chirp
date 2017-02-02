package storage

import (
	"time"

	"github.com/VirrageS/chirp/backend/model"
)

type tweetsDataAccessor interface {
	GetUsersTweets(userID, requestingUserID int64) ([]*model.Tweet, error)
	GetTweetsByAuthorIDs(authorsIDs []int64, requestingUserID int64) ([]*model.Tweet, error)
	GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error)
	InsertTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error)
	DeleteTweet(tweetID, requestingUserID int64) error
	LikeTweet(tweetID, userID int64) error
	UnlikeTweet(tweetID, userID int64) error
	GetTweetsUsingQueryString(querystring string, requestingUserID int64) ([]*model.Tweet, error)
}

type usersDataAccessor interface {
	GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error)
	GetUserByEmail(email string) (*model.User, error)
	InsertUser(user *model.NewUserForm) (*model.PublicUser, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error
	FollowUser(followeeID, followerID int64) error
	UnfollowUser(followeeID, followerID int64) error
	GetFollowers(userID, requestingUserID int64) ([]*model.PublicUser, error)
	GetFollowees(userID, requestingUserID int64) ([]*model.PublicUser, error)
	GetFolloweesIDs(userID int64) ([]int64, error)
	GetUsersUsingQueryString(querystring string, requestingUserID int64) ([]*model.PublicUser, error)
}

// Accessor is interface which defines all functions used on database/cache/fts
// in the system. Any other packages should use this Accessor instead of using
// eg. database directly.
type Accessor interface {
	usersDataAccessor
	tweetsDataAccessor
}
