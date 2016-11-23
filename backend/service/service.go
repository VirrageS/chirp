package service

import (
	"errors"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"

	APIModel "github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	databaseModel "github.com/VirrageS/chirp/backend/database/model"
	"github.com/VirrageS/chirp/backend/service/converters"
)

// Struct that implements APIProvider
type Service struct {
	// logger?
	// DB to API model converter?
	configuration config.ServiceConfigProvider
	db            database.DatabaseAccessor

	userConverter converters.UserModelConverter
}

// Constructs a Service that uses provided objects
func NewService(databaseAccessor database.DatabaseAccessor,
	configuration config.ServiceConfigProvider,
	userConverter converters.UserModelConverter) ServiceProvider {

	return &Service{
		configuration: configuration,
		db:            databaseAccessor,
		userConverter: userConverter,
	}
}

func (service *Service) GetTweets() ([]*APIModel.Tweet, *Error) {
	databaseTweets, err := service.db.GetTweets()
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	APITweets := service.convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets)

	return APITweets, nil
}

// Use GetTweets() with filtering parameters instead, when filtering will be supported
func (service *Service) GetTweetsOfUserWithID(userID int64) ([]*APIModel.Tweet, *Error) {
	databaseTweets, err := service.db.GetTweetsOfUserWithID(userID)
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	APITweets := service.convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets)

	return APITweets, nil
}

func (service *Service) GetTweet(tweetID int64) (*APIModel.Tweet, *Error) {
	databaseTweet, err := service.db.GetTweet(tweetID)

	if err == database.NoResults {
		return nil, &Error{
			Code: http.StatusNotFound,
			Err:  errors.New("Tweet with given ID was not found."),
		}
	}
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	APITweet := service.convertDatabaseTweetToAPITweet(&databaseTweet)

	return APITweet, nil
}

func (service *Service) PostTweet(newTweet *APIModel.NewTweet) (*APIModel.Tweet, *Error) {
	// TODO: reject if content is empty or when user submitted the same tweet more than once
	databaseTweet := service.convertAPINewTweetToDatabaseTweet(newTweet)

	addedTweet, err := service.db.InsertTweet(*databaseTweet)
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	APITweet := service.convertDatabaseTweetToAPITweet(&addedTweet)

	return APITweet, nil
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

func (service *Service) GetUsers() ([]*APIModel.User, *Error) {
	databaseUsers, databaseError := service.db.GetUsers()

	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	APIUsers := service.userConverter.ConvertArrayDatabasePublicUserToAPI(databaseUsers)

	return APIUsers, nil
}

func (service *Service) GetUser(userId int64) (*APIModel.User, *Error) {
	databaseUser, databaseError := service.db.GetUserByID(userId)

	if databaseError == database.NoResults {
		return nil, &Error{
			Code: http.StatusNotFound,
			Err:  errors.New("User with given ID was not found."),
		}
	}
	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	APIUser := service.userConverter.ConvertDatabasePublicUserToAPI(databaseUser)

	return APIUser, nil
}

func (service *Service) RegisterUser(newUserForm *APIModel.NewUserForm) (*APIModel.User, *Error) {
	databaseUser := service.userConverter.ConvertAPItoDatabase(newUserForm)

	newUser, err := service.db.InsertUser(databaseUser)

	if err == database.UserAlreadyExistsError {
		return nil, &Error{
			Code: http.StatusConflict,
			Err:  errors.New("User with given username or email already exists."),
		}
	}
	if err == database.DatabaseError {
		return nil, UnexpectedError
	}

	apiUser := service.userConverter.ConvertDatabaseToAPI(newUser)

	return apiUser, nil
}

func (service *Service) LoginUser(loginForm *APIModel.LoginForm) (*APIModel.LoginResponse, *Error) {
	email := loginForm.Email
	password := loginForm.Password

	databaseUser, databaseError := service.db.GetUserByEmail(&email)
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
	if databaseUser.Password != password {
		return nil, &Error{
			Code: http.StatusUnauthorized,
			Err:  errors.New("Invalid email or password."),
		}
	}

	loginTime := time.Now()
	updateError := service.db.UpdateUserLastLoginTime(databaseUser.ID, &loginTime)
	if updateError == database.DatabaseError {
		return nil, UnexpectedError
	}

	token, serviceError := service.createTokenForUser(databaseUser)
	if serviceError != nil {
		return nil, serviceError
	}

	apiUser := service.userConverter.ConvertDatabaseToAPI(databaseUser)
	response := &APIModel.LoginResponse{
		AuthToken: token,
		User:      apiUser,
	}

	return response, nil
}

//TODO: Maybe move out to another package and inject tokenCreator object/closure that can create tokens
func (service *Service) createTokenForUser(user *databaseModel.User) (*string, *Error) {
	validityDuration := time.Duration(service.configuration.GetTokenValidityPeriod())
	expirationTime := time.Now().Add(validityDuration * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(service.configuration.GetSecretKey())
	if err != nil {
		log.WithError(err).Fatal("Failed to sign the token.")
		return nil, UnexpectedError
	}

	return &tokenString, nil
}

func (service *Service) convertDatabaseTweetToAPITweet(tweet *databaseModel.TweetWithAuthor) *APIModel.Tweet {
	tweetID := tweet.ID
	author := tweet.Author
	likes := tweet.Likes
	retweets := tweet.Retweets
	createdAt := tweet.CreatedAt
	content := tweet.Content

	APIauthor := service.userConverter.ConvertDatabasePublicUserToAPI(author)

	APITweet := &APIModel.Tweet{
		ID:        tweetID,
		Author:    APIauthor,
		Likes:     likes,
		Retweets:  retweets,
		CreatedAt: createdAt,
		Content:   content,
		Liked:     false,
		Retweeted: false,
	}
	return APITweet
}

func (service *Service) convertAPINewTweetToDatabaseTweet(tweet *APIModel.NewTweet) *databaseModel.Tweet {
	authorId := tweet.AuthorID
	content := tweet.Content

	return &databaseModel.Tweet{
		ID:        0,
		AuthorID:  authorId,
		Likes:     0,
		Retweets:  0,
		CreatedAt: time.Now(),
		Content:   content,
	}
}

func (service *Service) convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets []databaseModel.TweetWithAuthor) []*APIModel.Tweet {
	APITweets := make([]*APIModel.Tweet, 0)

	for _, databaseTweet := range databaseTweets {
		APITweet := service.convertDatabaseTweetToAPITweet(&databaseTweet)
		APITweets = append(APITweets, APITweet)
	}

	return APITweets
}
