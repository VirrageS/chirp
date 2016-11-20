package service

import APIModel "github.com/VirrageS/chirp/backend/api/model"

// TODO: Maybe split into 2 services: tweet and user service?
type ServiceProvider interface {
	GetTweets() ([]*APIModel.Tweet, *Error)
	GetTweetsOfUserWithID(userID int64) ([]*APIModel.Tweet, *Error)
	GetTweet(tweetID int64) (*APIModel.Tweet, *Error)
	PostTweet(newTweet *APIModel.NewTweet) (*APIModel.Tweet, *Error)
	DeleteTweet(userID, tweetID int64) *Error

	GetUsers() ([]*APIModel.User, *Error)
	GetUser(userId int64) (*APIModel.User, *Error)
	RegisterUser(newUserForm *APIModel.NewUserForm) (*APIModel.User, *Error)
	LoginUser(loginForm *APIModel.LoginForm) (*APIModel.LoginResponse, *Error)
}
