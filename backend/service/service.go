package service

import (
	"time"

	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
	"github.com/VirrageS/chirp/backend/token"
)

// Struct that implements APIProvider
type Service struct {
	// logger?
	config       config.ServiceConfigProvider
	db           database.DatabaseAccessor
	tokenManager token.TokenManagerProvider
}

// Constructs a Service that uses provided objects
func NewService(config config.ServiceConfigProvider, database database.DatabaseAccessor, tokenManager token.TokenManagerProvider) ServiceProvider {
	return &Service{
		config:       config,
		db:           database,
		tokenManager: tokenManager,
	}
}

func (service *Service) GetTweets() ([]*model.Tweet, error) {
	tweets, err := service.db.GetTweets()
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

// Use GetTweets() with filtering parameters instead, when filtering will be supported
func (service *Service) GetTweetsOfUserWithID(userID int64) ([]*model.Tweet, error) {
	tweets, err := service.db.GetTweetsOfUserWithID(userID)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func (service *Service) GetTweet(tweetID int64) (*model.Tweet, error) {
	tweet, err := service.db.GetTweet(tweetID)

	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (service *Service) PostTweet(tweet *model.NewTweet) (*model.Tweet, error) {
	// TODO: reject if content is empty or when user submitted the same tweet more than once
	newTweet, err := service.db.InsertTweet(tweet)
	if err != nil {
		return nil, err
	}

	return newTweet, nil
}

func (service *Service) DeleteTweet(userID, tweetID int64) error {
	// TODO: Maybe fetch Tweet not TweetWithAuthor
	databaseTweet, err := service.db.GetTweet(tweetID)

	if err != nil {
		return err
	}

	if databaseTweet.Author.ID != userID {
		return errors.ForbiddenError
	}

	err = service.db.DeleteTweet(tweetID)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) GetUsers() ([]*model.PublicUser, error) {
	users, err := service.db.GetUsers()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *Service) GetUser(userID int64) (*model.PublicUser, error) {
	user, err := service.db.GetUserByID(userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *Service) RegisterUser(newUserForm *model.NewUserForm) (*model.PublicUser, error) {
	newUser, err := service.db.InsertUser(newUserForm)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// TODO: fix this - maybe service should fetch only 'auth' data and then get fetch user data and return it
func (service *Service) LoginUser(loginForm *model.LoginForm) (*model.LoginResponse, error) {
	email := loginForm.Email
	password := loginForm.Password

	user, databaseError := service.db.GetUserByEmail(email)
	if databaseError == errors.NoResultsError {
		return nil, errors.InvalidCredentialsError // return 401 when user with given email is not found
	} else if databaseError != nil {
		return nil, databaseError
	}

	// TODO: hash the password before comparing
	if user.Password != password {
		return nil, errors.InvalidCredentialsError
	}

	loginTime := time.Now()
	updateError := service.db.UpdateUserLastLoginTime(user.ID, &loginTime)
	if updateError != nil {
		return nil, updateError
	}

	authToken, err := service.createAuthToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := service.createRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	response := &model.LoginResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
		User: &model.PublicUser{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			AvatarUrl: user.AvatarUrl.String,
			Following: user.Following,
		},
	}

	return response, nil
}

func (service *Service) RefreshAuthToken(request *model.RefreshAuthTokenRequest) (*model.RefreshAuthTokenResponse, error) {
	authToken, err := service.createAuthToken(request.UserID)
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
