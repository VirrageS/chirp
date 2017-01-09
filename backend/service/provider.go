package service

import model "github.com/VirrageS/chirp/backend/model"

// TODO: Maybe split into 2 services: tweet and user service?
type ServiceProvider interface {
	GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error)
	GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error)
	PostTweet(newTweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error)
	DeleteTweet(tweetID, requestingUserID int64) error
	LikeTweet(tweetID, requestingUserID int64) (*model.Tweet, error)
	UnlikeTweet(tweetID, requestingUserID int64) (*model.Tweet, error)

	GetUser(userID, requestingUserID int64) (*model.PublicUser, error)
	FollowUser(userID, requestingUserID int64) (*model.PublicUser, error)
	UnfollowUser(userID, requestingUserID int64) (*model.PublicUser, error)
	UserFollowers(userID, requestingUserID int64) ([]*model.PublicUser, error)
	UserFollowees(userID, requestingUserID int64) ([]*model.PublicUser, error)

	FullTextSearch(queryString string, requestingUserID int64) (*model.FullTextSearchResponse, error)

	RegisterUser(newUserForm *model.NewUserForm) (*model.PublicUser, error)
	LoginUser(loginForm *model.LoginForm) (*model.PublicUser, error)
	CreateOrLoginUserWithGoogle(userGoogle *model.UserGoogle) (*model.PublicUser, error)
}
