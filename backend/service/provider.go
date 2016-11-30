package service

import model "github.com/VirrageS/chirp/backend/model"

// TODO: Maybe split into 2 services: tweet and user service?
type ServiceProvider interface {
	GetTweets() ([]*model.Tweet, *Error)
	GetTweetsOfUserWithID(userID int64) ([]*model.Tweet, *Error)
	GetTweet(tweetID int64) (*model.Tweet, *Error)
	PostTweet(newTweet *model.NewTweet) (*model.Tweet, *Error)
	DeleteTweet(userID, tweetID int64) *Error

	GetUsers() ([]*model.PublicUser, *Error)
	GetUser(userId int64) (*model.PublicUser, *Error)
	RegisterUser(newUserForm *model.NewUserForm) (*model.PublicUser, *Error)
	LoginUser(loginForm *model.LoginForm) (*model.LoginResponse, *Error)
}
