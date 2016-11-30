package service

import (
	"errors"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"

	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
)

// Struct that implements APIProvider
type Service struct {
	// logger?
	// DB to API model converter?
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

func (service *Service) GetTweets() ([]*model.Tweet, *Error) {
	tweets, err := service.db.GetTweets()
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	return tweets, nil
}

// Use GetTweets() with filtering parameters instead, when filtering will be supported
func (service *Service) GetTweetsOfUserWithID(userID int64) ([]*model.Tweet, *Error) {
	tweets, err := service.db.GetTweetsOfUserWithID(userID)
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	return tweets, nil
}

func (service *Service) GetTweet(tweetID int64) (*model.Tweet, *Error) {
	tweet, err := service.db.GetTweet(tweetID)

	if err == database.NoResults {
		return nil, &Error{
			Code: http.StatusNotFound,
			Err:  errors.New("Tweet with given ID was not found."),
		}
	}
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	return tweet, nil
}

func (service *Service) PostTweet(tweet *model.NewTweet) (*model.Tweet, *Error) {
	// TODO: reject if content is empty or when user submitted the same tweet more than once
	newTweet, err := service.db.InsertTweet(tweet)
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	return newTweet, nil
}

func (service *Service) DeleteTweet(userID, tweetID int64) *Error {
	// TODO: Maybe fetch Tweet not TweetWithAuthor
	databaseTweet, err := service.db.GetTweet(tweetID)

	if err == database.NoResults {
		return &Error{
			Code: http.StatusNotFound,
			Err:  errors.New("Tweet with given ID was not found."),
		}
	}
	if err == database.DatabaseError {
		return UnexpectedError
	}

	if databaseTweet.Author.ID != userID {
		return &Error{
			Code: http.StatusForbidden,
			Err:  errors.New("User is not allowed to modify this resource."),
		}
	}

	err = service.db.DeleteTweet(tweetID)
	if err == database.DatabaseError {
		return UnexpectedError
	}

	return nil
}

func (service *Service) GetUsers() ([]*model.PublicUser, *Error) {
	users, databaseError := service.db.GetUsers()

	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	return users, nil
}

func (service *Service) GetUser(userId int64) (*model.PublicUser, *Error) {
	user, databaseError := service.db.GetUserByID(userId)

	if databaseError == database.NoResults {
		return nil, &Error{
			Code: http.StatusNotFound,
			Err:  errors.New("User with given ID was not found."),
		}
	}
	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	return user, nil
}

func (service *Service) RegisterUser(newUserForm *model.NewUserForm) (*model.PublicUser, *Error) {
	newUser, err := service.db.InsertUser(newUserForm)

	if err == database.UserAlreadyExistsError {
		return nil, &Error{
			Code: http.StatusConflict,
			Err:  errors.New("User with given username or email already exists."),
		}
	}
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	return newUser, nil
}

func (service *Service) LoginUser(loginForm *model.LoginForm) (*model.LoginResponse, *Error) {
	email := loginForm.Email
	password := loginForm.Password

	user, databaseError := service.db.GetUserByEmail(email)
	if databaseError == database.NoResults {
		return nil, &Error{
			Code: http.StatusNotFound,
			Err:  errors.New("User with given ID was not found."),
		}
	}
	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	// TODO: hash the password before comparing
	if user.Password != password {
		return nil, &Error{
			Code: http.StatusUnauthorized,
			Err:  errors.New("Invalid email or password."),
		}
	}

	loginTime := time.Now()
	updateError := service.db.UpdateUserLastLoginTime(user.ID, &loginTime)
	if updateError == database.DatabaseError {
		return nil, UnexpectedError
	}

	token, serviceError := service.createTokenForUser(user)
	if serviceError != nil {
		return nil, serviceError
	}

	response := &model.LoginResponse{
		AuthToken: token,
		User:      user,
	}

	return response, nil
}

//TODO: Maybe move out to another package and inject tokenCreator object/closure that can create tokens
func (service *Service) createTokenForUser(user *model.User) (string, *Error) {
	validityDuration := time.Duration(service.configuration.GetTokenValidityPeriod())
	expirationTime := time.Now().Add(validityDuration * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(service.configuration.GetSecretKey())
	if err != nil {
		log.WithError(err).Fatal("Failed to sign the token.")
		return "", UnexpectedError
	}

	return tokenString, nil
}
