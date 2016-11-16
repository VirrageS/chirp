package service

import (
	"errors"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"

	"database/sql"
	APIModel "github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	databaseModel "github.com/VirrageS/chirp/backend/database/model"
	appErrors "github.com/VirrageS/chirp/backend/errors"
)

// TODO: Maybe split into 2 services: tweet and user service?
type ServiceProvider interface {
	GetTweets() ([]*APIModel.Tweet, *appErrors.AppError)
	GetTweetsOfUserWithID(userID int64) ([]*APIModel.Tweet, *appErrors.AppError)
	GetTweet(tweetID int64) (*APIModel.Tweet, *appErrors.AppError)
	PostTweet(newTweet *APIModel.NewTweet) (*APIModel.Tweet, *appErrors.AppError)
	DeleteTweet(userID, tweetID int64) *appErrors.AppError

	GetUsers() ([]*APIModel.User, *appErrors.AppError)
	GetUser(userId int64) (*APIModel.User, *appErrors.AppError)
	RegisterUser(newUserForm *APIModel.NewUserForm) (*APIModel.User, *appErrors.AppError)
	LoginUser(loginForm *APIModel.LoginForm) (*APIModel.LoginResponse, *appErrors.AppError)
}

type Service struct {
	// logger?
	// DB to API model converter?
	// API to DB model converter?
	configuration *config.ServiceConfig
	db            database.DatabaseAccessor
}

func NewService(databaseAccessor database.DatabaseAccessor, configuration *config.ServiceConfig) ServiceProvider {
	return &Service{
		configuration: configuration,
		db:            databaseAccessor,
	}
}

func (service *Service) GetTweets() ([]*APIModel.Tweet, *appErrors.AppError) {
	databaseTweets, databaseError := service.db.GetTweets()

	if databaseError != nil {
		return nil, appErrors.UnexpectedError
	}

	APITweets, serviceError := service.convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets)

	if serviceError != nil {
		return nil, serviceError
	}

	return APITweets, nil
}

// Use GetTweets() with filtering parameters instead, when filtering will be supported
func (service *Service) GetTweetsOfUserWithID(userID int64) ([]*APIModel.Tweet, *appErrors.AppError) {
	databaseTweets, databaseError := service.db.GetTweetsOfUserWithID(userID)

	if databaseError != nil {
		return nil, appErrors.UnexpectedError
	}

	APITweets, serviceError := service.convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets)

	if serviceError != nil {
		return nil, serviceError
	}

	return APITweets, nil
}

func (service *Service) GetTweet(tweetID int64) (*APIModel.Tweet, *appErrors.AppError) {
	databaseTweet, databaseError := service.db.GetTweet(tweetID)

	if databaseError != nil {
		// Later on we'll need to add type switch here to check the type of error, because several things
		// can go wrong when fetching data from database: not found, SQL error, db connection error etc
		return nil, &appErrors.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("User with given ID was not found."),
		}
	}

	APITweet, serviceError := service.convertDatabaseTweetToAPITweet(&databaseTweet)

	if serviceError != nil {
		return nil, serviceError
	}

	return APITweet, nil
}

func (service *Service) PostTweet(newTweet *APIModel.NewTweet) (*APIModel.Tweet, *appErrors.AppError) {
	// TODO: reject if content is empty or when user submitted the same tweet more than once
	databaseTweet := service.convertAPINewTweetToDatabaseTweet(newTweet)

	addedTweet, databaseError := service.db.InsertTweet(*databaseTweet)

	if databaseError != nil {
		// for now its an unexpected error, but later on we'll probably need an error type switch here too
		return nil, appErrors.UnexpectedError
	}

	APITweet, serviceError := service.convertDatabaseTweetToAPITweet(&addedTweet)

	if serviceError != nil {
		return nil, serviceError
	}

	return APITweet, nil
}

func (service *Service) DeleteTweet(userID, tweetID int64) *appErrors.AppError {
	databaseTweet, err := service.db.GetTweet(tweetID)

	if err != nil {
		return &appErrors.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("Tweet with given ID was not found."),
		}
	}
	if databaseTweet.AuthorID != userID {
		return &appErrors.AppError{
			Code: http.StatusForbidden,
			Err:  errors.New("User is not allowed to modify this resource."),
		}
	}

	err = service.db.DeleteTweet(tweetID)
	if err != nil {
		return appErrors.UnexpectedError
	}

	return nil
}

func (service *Service) GetUsers() ([]*APIModel.User, *appErrors.AppError) {
	databaseUsers, databaseError := service.db.GetUsers()

	if databaseError != nil {
		// for now its an unexpected error, but later on we'll probably need an error type switch here too
		return nil, appErrors.UnexpectedError
	}

	APIUsers := service.convertArrayOfDatabaseUsersToArrayOfAPIUsers(databaseUsers)

	return APIUsers, nil
}

func (service *Service) GetUser(userId int64) (*APIModel.User, *appErrors.AppError) {
	databaseUser, databaseError := service.db.GetUserByID(userId)

	if databaseError != nil {
		// Maybe later on we'll need to add type switch here to check the type of error, because several things
		// can go wrong when fetching data from database: not found, SQL error, db connection error etc
		return nil, &appErrors.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("User with given ID was not found."),
		}
	}

	APIUser := service.convertDatabaseUserToAPIUser(&databaseUser)

	return APIUser, nil
}

func (service *Service) RegisterUser(newUserForm *APIModel.NewUserForm) (*APIModel.User, *appErrors.AppError) {
	databaseUser := service.covertAPINewUserToDatabaseUser(newUserForm)

	newUser, err := service.db.InsertUser(*databaseUser)

	if err != nil {
		// again, one error only for now...
		return nil, &appErrors.AppError{
			Code: http.StatusConflict,
			Err:  errors.New("User with given username or email already exists."),
		}
	}

	apiUser := service.convertDatabaseUserToAPIUser(&newUser)

	return apiUser, nil
}

func (service *Service) LoginUser(loginForm *APIModel.LoginForm) (*APIModel.LoginResponse, *appErrors.AppError) {
	email := loginForm.Email
	password := loginForm.Password

	databaseUser, databaseError := service.db.GetUserByEmail(email)
	// TODO: hash the password before comparing
	if databaseError != nil || databaseUser.Password != password {
		return nil, &appErrors.AppError{
			Code: http.StatusUnauthorized,
			Err:  errors.New("Invalid email or password."),
		}
	}
	// TODO: update users last login time

	token, serviceError := service.createTokenForUser(&databaseUser)
	if serviceError != nil {
		return nil, serviceError
	}

	apiUser := service.convertDatabaseUserToAPIUser(&databaseUser)
	response := &APIModel.LoginResponse{
		AuthToken: token,
		User:      apiUser,
	}

	return response, nil
}

func (service *Service) createTokenForUser(user *databaseModel.User) (string, *appErrors.AppError) {
	validityDuration := time.Duration(service.configuration.TokenValidityPeriod)
	expirationTime := time.Now().Add(validityDuration * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(service.configuration.SecretKey)
	if err != nil {
		log.WithError(err).Error("Failed to sign the token.")
		// we should probably panic here, because the server will not be able to run if it can't auth users
		return "", appErrors.UnexpectedError
	}

	return tokenString, nil
}

func (service *Service) convertDatabaseTweetToAPITweet(tweet *databaseModel.Tweet) (*APIModel.Tweet, *appErrors.AppError) {
	tweetID := tweet.ID
	userID := tweet.AuthorID
	likes := tweet.Likes
	retweets := tweet.Retweets
	createdAt := tweet.CreatedAt
	content := tweet.Content

	authorFullData, err := service.db.GetUserByID(userID)

	if err != nil {
		// TODO: here we will also need to check the error type and have different handling for different erros
		log.WithFields(log.Fields{
			"tweetID": tweetID,
			"userID":  userID,
		}).Error("Failed to convert database tweet to API tweet. User was not found in database.")
		return nil, appErrors.UnexpectedError
	}

	author := service.convertDatabaseUserToAPIUser(&authorFullData)

	APITweet := &APIModel.Tweet{
		ID:        tweetID,
		Author:    author,
		Likes:     likes,
		Retweets:  retweets,
		CreatedAt: createdAt,
		Content:   content,
		Liked:     false,
		Retweeted: false,
	}
	return APITweet, nil
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

func (service *Service) convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets []databaseModel.Tweet) ([]*APIModel.Tweet, *appErrors.AppError) {
	APITweets := make([]*APIModel.Tweet, 0)

	for _, databaseTweet := range databaseTweets {
		APITweet, err := service.convertDatabaseTweetToAPITweet(&databaseTweet)

		if err != nil {
			return nil, appErrors.UnexpectedError
		}

		APITweets = append(APITweets, APITweet)
	}

	return APITweets, nil
}

func (service *Service) convertArrayOfDatabaseUsersToArrayOfAPIUsers(databaseUsers []databaseModel.User) []*APIModel.User {
	convertedUsers := make([]*APIModel.User, 0)

	for _, databaseUser := range databaseUsers {
		APIUser := service.convertDatabaseUserToAPIUser(&databaseUser)
		convertedUsers = append(convertedUsers, APIUser)
	}

	return convertedUsers
}

func (service *Service) convertDatabaseUserToAPIUser(user *databaseModel.User) *APIModel.User {
	id := user.ID
	username := user.Username
	email := user.Email
	createdAt := user.CreatedAt
	lastLogin := user.LastLogin
	name := user.Name
	avatarUrl := user.AvatarUrl.String

	return &APIModel.User{
		ID:        id,
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
		LastLogin: lastLogin,
		Name:      name,
		AvatarUrl: avatarUrl,
		Following: false,
	}
}

func (service *Service) covertAPINewUserToDatabaseUser(user *APIModel.NewUserForm) *databaseModel.User {
	username := user.Username
	password := user.Password
	email := user.Email
	name := user.Name
	creationTime := time.Now()

	return &databaseModel.User{
		ID:            0,
		TwitterToken:  toSqlNullString(""),
		FacebookToken: toSqlNullString(""),
		GoogleToken:   toSqlNullString(""),
		Username:      username,
		Password:      password,
		Email:         email,
		CreatedAt:     creationTime,
		LastLogin:     creationTime,
		Active:        true,
		Name:          name,
		AvatarUrl:     toSqlNullString(""),
	}
}

// converts string to database NullString
// TODO: Maybe move to a new 'util' package
func toSqlNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
