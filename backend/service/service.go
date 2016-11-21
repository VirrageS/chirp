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
)

// Struct that implements APIProvider
type Service struct {
	// logger?
	// DB to API model converter?
	// API to DB model converter?
	configuration config.ServiceConfigProvider
	db            database.DatabaseAccessor
}

// Constructs a Service that uses given DatabaseAccessor with configuration provided by given ServiceConfigProvider
func NewService(databaseAccessor database.DatabaseAccessor, configuration config.ServiceConfigProvider) ServiceProvider {
	return &Service{
		configuration: configuration,
		db:            databaseAccessor,
	}
}

func (service *Service) GetTweets() ([]*APIModel.Tweet, *Error) {
	databaseTweets, databaseError := service.db.GetTweets()

	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	APITweets, serviceError := service.convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets)

	if serviceError != nil {
		return nil, serviceError
	}

	return APITweets, nil
}

// Use GetTweets() with filtering parameters instead, when filtering will be supported
func (service *Service) GetTweetsOfUserWithID(userID int64) ([]*APIModel.Tweet, *Error) {
	databaseTweets, databaseError := service.db.GetTweetsOfUserWithID(userID)

	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	APITweets, serviceError := service.convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets)

	if serviceError != nil {
		return nil, serviceError
	}

	return APITweets, nil
}

func (service *Service) GetTweet(tweetID int64) (*APIModel.Tweet, *Error) {
	databaseTweet, databaseError := service.db.GetTweet(tweetID)

	if databaseError == database.NoResults {
		return nil, &Error{
			Code: http.StatusNotFound,
			Err:  errors.New("Tweet with given ID was not found."),
		}
	}
	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	APITweet, serviceError := service.convertDatabaseTweetToAPITweet(&databaseTweet)

	if serviceError != nil {
		return nil, serviceError
	}

	return APITweet, nil
}

func (service *Service) PostTweet(newTweet *APIModel.NewTweet) (*APIModel.Tweet, *Error) {
	// TODO: reject if content is empty or when user submitted the same tweet more than once
	databaseTweet := service.convertAPINewTweetToDatabaseTweet(newTweet)

	addedTweet, databaseError := service.db.InsertTweet(*databaseTweet)

	if databaseError == database.DatabaseError {
		return nil, UnexpectedError
	}

	APITweet, serviceError := service.convertDatabaseTweetToAPITweet(&addedTweet)

	if serviceError != nil {
		return nil, serviceError
	}

	return APITweet, nil
}

func (service *Service) DeleteTweet(userID, tweetID int64) *Error {
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

	if databaseTweet.AuthorID != userID {
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

	APIUsers := service.convertArrayOfDatabaseUsersToArrayOfAPIUsers(databaseUsers)

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

	APIUser := service.convertDatabaseUserToAPIUser(databaseUser)

	return APIUser, nil
}

func (service *Service) RegisterUser(newUserForm *APIModel.NewUserForm) (*APIModel.User, *Error) {
	databaseUser := service.covertAPINewUserToDatabaseUser(newUserForm)

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

	apiUser := service.convertDatabaseUserToAPIUser(newUser)

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

	apiUser := service.convertDatabaseUserToAPIUser(databaseUser)
	response := &APIModel.LoginResponse{
		AuthToken: token,
		User:      apiUser,
	}

	return response, nil
}

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

// TODO: Maybe the converter should not access database and databaseModel.Tweet should contain whole user data
func (service *Service) convertDatabaseTweetToAPITweet(tweet *databaseModel.Tweet) (*APIModel.Tweet, *Error) {
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
		return nil, UnexpectedError
	}

	author := service.convertDatabaseUserToAPIUser(authorFullData)

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

func (service *Service) convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets []databaseModel.Tweet) ([]*APIModel.Tweet, *Error) {
	APITweets := make([]*APIModel.Tweet, 0)

	for _, databaseTweet := range databaseTweets {
		APITweet, err := service.convertDatabaseTweetToAPITweet(&databaseTweet)

		if err != nil {
			return nil, UnexpectedError
		}

		APITweets = append(APITweets, APITweet)
	}

	return APITweets, nil
}

func (service *Service) convertArrayOfDatabaseUsersToArrayOfAPIUsers(databaseUsers []*databaseModel.User) []*APIModel.User {
	convertedUsers := make([]*APIModel.User, 0)

	for _, databaseUser := range databaseUsers {
		APIUser := service.convertDatabaseUserToAPIUser(databaseUser)
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
