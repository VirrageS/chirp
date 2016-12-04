package service

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"

	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
)

// Struct that implements APIProvider
type Service struct {
	// logger?
	configuration config.ServiceConfigProvider
	db            database.DatabaseAccessor
}

// Constructs a Service that uses provided objects
func NewService(databaseAccessor database.DatabaseAccessor, configuration config.ServiceConfigProvider) ServiceProvider {
	return &Service{
		configuration: configuration,
		db:            databaseAccessor,
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

func (service *Service) GetUser(userId int64) (*model.PublicUser, error) {
	user, err := service.db.GetUserByID(userId)

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
	if databaseError != nil {
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

	token, serviceError := service.createTokenForUser(user)
	if serviceError != nil {
		return nil, serviceError
	}

	response := &model.LoginResponse{
		AuthToken: token,
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

//TODO: Maybe move out to another package and inject tokenCreator object/closure that can create tokens
func (service *Service) createTokenForUser(user *model.User) (string, error) {
	validityDuration := time.Duration(service.configuration.GetTokenValidityPeriod())
	expirationTime := time.Now().Add(validityDuration * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(service.configuration.GetSecretKey())
	if err != nil {
		log.WithError(err).Fatal("Failed to sign the token.")
		return "", errors.UnexpectedError
	}

	return tokenString, nil
}
