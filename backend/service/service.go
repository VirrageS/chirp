package service

import (
	"time"

	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
	"github.com/VirrageS/chirp/backend/storage"
)

// Struct that implements APIProvider
type Service struct {
	storage storage.StorageAccessor
}

// Constructs a Service that uses provided objects
func NewService(
	storage storage.StorageAccessor,
) ServiceProvider {
	return &Service{
		storage: storage,
	}
}

func (service *Service) GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error) {
	tweets, err := service.storage.GetUsersTweets(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func (service *Service) GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	tweet, err := service.storage.GetTweet(tweetID, requestingUserID)

	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (service *Service) PostTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error) {
	// TODO: reject if content is empty or when user submitted the same tweet more than once
	newTweet, err := service.storage.InsertTweet(tweet, requestingUserID)
	if err != nil {
		return nil, err
	}

	return newTweet, nil
}

func (service *Service) DeleteTweet(tweetID, requestingUserID int64) error {
	// TODO: Maybe fetch Tweet not TweetWithAuthor
	databaseTweet, err := service.storage.GetTweet(tweetID, requestingUserID)

	if err != nil {
		return err
	}

	if databaseTweet.Author.ID != requestingUserID {
		return errors.ForbiddenError
	}

	err = service.storage.DeleteTweet(tweetID, requestingUserID)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) LikeTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	err := service.storage.LikeTweet(tweetID, requestingUserID)
	if err != nil {
		return nil, err
	}

	tweet, err := service.GetTweet(tweetID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (service *Service) UnlikeTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	err := service.storage.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		return nil, err
	}

	tweet, err := service.GetTweet(tweetID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (service *Service) GetUser(userID, requestingUserID int64) (*model.PublicUser, error) {
	user, err := service.storage.GetUserByID(userID, requestingUserID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *Service) FollowUser(userID, requestingUserID int64) (*model.PublicUser, error) {
	err := service.storage.FollowUser(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	user, err := service.storage.GetUserByID(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *Service) UnfollowUser(userID, requestingUserID int64) (*model.PublicUser, error) {
	err := service.storage.UnfollowUser(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	user, err := service.storage.GetUserByID(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *Service) UserFollowers(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followers, err := service.storage.GetFollowers(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (service *Service) UserFollowees(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followers, err := service.storage.GetFollowees(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (service *Service) FullTextSearch(queryString string, requestingUserID int64) (*model.FullTextSearchResponse, error) {
	tweets, err := service.storage.GetTweetsUsingQueryString(queryString, requestingUserID)
	if err != nil {
		return nil, err
	}

	users, err := service.storage.GetUsersUsingQueryString(queryString, requestingUserID)
	if err != nil {
		return nil, err
	}

	result := &model.FullTextSearchResponse{
		Users:  users,
		Tweets: tweets,
	}

	return result, nil
}

func (service *Service) RegisterUser(newUserForm *model.NewUserForm) (*model.PublicUser, error) {
	newUser, err := service.storage.InsertUser(newUserForm)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *Service) LoginUser(loginForm *model.LoginForm) (*model.PublicUser, error) {
	email := loginForm.Email
	password := loginForm.Password

	userAuthData, databaseError := service.storage.GetUserByEmail(email)
	if databaseError == errors.NoResultsError {
		return nil, errors.InvalidCredentialsError // return 401 when user with given email is not found
	} else if databaseError != nil {
		return nil, databaseError
	}

	// TODO: hash the password before comparing
	if userAuthData.Password != password {
		return nil, errors.InvalidCredentialsError
	}

	loginTime := time.Now()
	updateError := service.storage.UpdateUserLastLoginTime(userAuthData.ID, &loginTime)
	if updateError != nil {
		return nil, updateError
	}

	return service.storage.GetUserByID(userAuthData.ID, userAuthData.ID)
}

func (service *Service) CreateOrLoginUserWithGoogle(newUserGoogle *model.UserGoogle) (*model.PublicUser, error) {
	user, err := service.storage.GetUserByEmail(newUserGoogle.Email)
	if err == errors.NoResultsError {
		// TODO: add picture field and mark that this user is connected to google
		newUserForm := &model.NewUserForm{
			Username: newUserGoogle.GivenName,
			Password: "superwierdhash", // TODO: change this (this should be super unhackable some kind of 128 char hash)
			Email:    newUserGoogle.Email,
			Name:     newUserGoogle.Name,
		}

		_, err := service.storage.InsertUser(newUserForm)
		if err != nil {
			return nil, err
		}

		// get user after creating it (now we should be able to do it...)
		user, err = service.storage.GetUserByEmail(newUserGoogle.Email)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	loginForm := &model.LoginForm{
		Email:    user.Email,
		Password: user.Password,
	}

	return service.LoginUser(loginForm)
}
