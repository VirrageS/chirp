package service

import model "github.com/VirrageS/chirp/backend/model"

// TODO: Maybe split into 2 services: tweet and user service?
type ServiceProvider interface {
	GetTweets(requestingUserID int64) ([]*model.Tweet, error)
	GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error)
	GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error)
	PostTweet(newTweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error)
	DeleteTweet(tweetID, requestingUserID int64) error
	LikeTweet(tweetID, requestingUserID int64) (*model.Tweet, error)
	UnlikeTweet(tweetID, requestingUserID int64) (*model.Tweet, error)

	GetUsers() ([]*model.PublicUser, error)
	GetUser(userId int64) (*model.PublicUser, error)
	RegisterUser(newUserForm *model.NewUserForm) (*model.PublicUser, error)
	LoginUser(loginForm *model.LoginForm) (*model.LoginResponse, error)
	RefreshAuthToken(request *model.RefreshAuthTokenRequest) (*model.RefreshAuthTokenResponse, error)
}
