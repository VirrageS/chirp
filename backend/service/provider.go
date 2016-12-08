package service

import model "github.com/VirrageS/chirp/backend/model"

// TODO: Maybe split into 2 services: tweet and user service?
type ServiceProvider interface {
	GetTweets() ([]*model.Tweet, error)
	GetTweetsOfUserWithID(userID int64) ([]*model.Tweet, error)
	GetTweet(tweetID int64) (*model.Tweet, error)
	PostTweet(newTweet *model.NewTweet) (*model.Tweet, error)
	DeleteTweet(userID, tweetID int64) error

	GetUsers() ([]*model.PublicUser, error)
	GetUser(userId int64) (*model.PublicUser, error)
	RegisterUser(newUserForm *model.NewUserForm) (*model.PublicUser, error)
	LoginUser(loginForm *model.LoginForm) (*model.LoginResponse, error)
	RefreshAuthToken(request *model.RefreshAuthTokenRequest) (*model.RefreshAuthTokenResponse, error)
}
