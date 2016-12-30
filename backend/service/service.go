package service

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
	"github.com/VirrageS/chirp/backend/storage"
	"github.com/VirrageS/chirp/backend/token"
)

// Struct that implements APIProvider
type Service struct {
	// logger?
	config       config.ServiceConfigProvider
	db           storage.DatabaseAccessor
	tokenManager token.TokenManagerProvider
}

// Constructs a Service that uses provided objects
func NewService(config config.ServiceConfigProvider, database storage.DatabaseAccessor, tokenManager token.TokenManagerProvider) ServiceProvider {
	return &Service{
		config:       config,
		db:           database,
		tokenManager: tokenManager,
	}
}

func (service *Service) GetTweets(requestingUserID int64) ([]*model.Tweet, error) {
	tweets, err := service.db.GetTweets(requestingUserID)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

// Use GetTweets() with filtering parameters instead, when filtering will be supported
func (service *Service) GetTweetsOfUserWithID(userID, requestingUserID int64) ([]*model.Tweet, error) {
	tweets, err := service.db.GetTweetsOfUserWithID(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func (service *Service) GetTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	tweet, err := service.db.GetTweet(tweetID, requestingUserID)

	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (service *Service) PostTweet(tweet *model.NewTweet, requestingUserID int64) (*model.Tweet, error) {
	// TODO: reject if content is empty or when user submitted the same tweet more than once
	newTweet, err := service.db.InsertTweet(tweet, requestingUserID)
	if err != nil {
		return nil, err
	}

	return newTweet, nil
}

func (service *Service) DeleteTweet(tweetID, requestingUserID int64) error {
	// TODO: Maybe fetch Tweet not TweetWithAuthor
	databaseTweet, err := service.db.GetTweet(tweetID, requestingUserID)

	if err != nil {
		return err
	}

	if databaseTweet.Author.ID != requestingUserID {
		return errors.ForbiddenError
	}

	err = service.db.DeleteTweet(tweetID, requestingUserID)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) LikeTweet(tweetID, requestingUserID int64) (*model.Tweet, error) {
	err := service.db.LikeTweet(tweetID, requestingUserID)
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
	err := service.db.UnlikeTweet(tweetID, requestingUserID)
	if err != nil {
		return nil, err
	}

	tweet, err := service.GetTweet(tweetID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (service *Service) GetUsers(requestingUserID int64) ([]*model.PublicUser, error) {
	users, err := service.db.GetUsers(requestingUserID)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *Service) GetUser(userID, requestingUserID int64) (*model.PublicUser, error) {
	user, err := service.db.GetUserByID(userID, requestingUserID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *Service) FollowUser(userID, requestingUserID int64) (*model.PublicUser, error) {
	err := service.db.FollowUser(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	user, err := service.db.GetUserByID(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *Service) UnfollowUser(userID, requestingUserID int64) (*model.PublicUser, error) {
	err := service.db.UnfollowUser(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	user, err := service.db.GetUserByID(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *Service) UserFollowers(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followers, err := service.db.Followers(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (service *Service) UserFollowees(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followers, err := service.db.Followees(userID, requestingUserID)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (service *Service) RegisterUser(newUserForm *model.NewUserForm) (*model.PublicUser, error) {
	newUser, err := service.db.InsertUser(newUserForm)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *Service) LoginUser(loginForm *model.LoginForm) (*model.LoginResponse, error) {
	email := loginForm.Email
	password := loginForm.Password

	userAuthData, databaseError := service.db.GetAuthDataOfUserWithEmail(email)
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
	updateError := service.db.UpdateUserLastLoginTime(userAuthData.ID, &loginTime)
	if updateError != nil {
		return nil, updateError
	}

	authToken, err := service.createAuthToken(userAuthData.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := service.createRefreshToken(userAuthData.ID)
	if err != nil {
		return nil, err
	}

	user, err := service.db.GetUserByID(userAuthData.ID, userAuthData.ID)
	if err != nil {
		return nil, err
	}

	response := &model.LoginResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	return response, nil
}

func (service *Service) RefreshAuthToken(request *model.RefreshAuthTokenRequest) (*model.RefreshAuthTokenResponse, error) {
	userID, err := service.tokenManager.ValidateToken(request.RefreshToken)
	if err != nil {
		log.WithError(err).Error("Error validating token in RefreshAuthToken.")
		return nil, err
	}

	// check if authenticating user exists
	_, err = service.db.GetUserByID(userID, userID)
	if err == errors.NoResultsError {
		return nil, errors.NotExistingUserAuthenticatingError
	}
	if err != nil {
		return nil, err
	}

	// generate new auth token for the user
	authToken, err := service.createAuthToken(userID)
	if err != nil {
		return nil, err
	}

	response := &model.RefreshAuthTokenResponse{
		AuthToken: authToken,
	}

	return response, nil
}

func (service *Service) createAuthToken(userID int64) (string, error) {
	return service.tokenManager.CreateToken(
		userID,
		service.config.GetAuthTokenValidityPeriod(),
	)
}

func (service *Service) createRefreshToken(userID int64) (string, error) {
	return service.tokenManager.CreateToken(
		userID,
		service.config.GetRefreshTokenValidityPeriod(),
	)
}
