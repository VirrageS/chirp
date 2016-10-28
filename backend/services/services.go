package services

import (
	"errors"

	APIModel "github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/database"
	databaseModel "github.com/VirrageS/chirp/backend/database/model"
	appErrors "github.com/VirrageS/chirp/backend/errors"

	"net/http"
	"time"
)

func GetTweets() ([]APIModel.Tweet, *appErrors.AppError) {
	databaseTweets, err := database.GetTweets()

	if err != nil {
		return nil, appErrors.UnexpectedError
	}

	APITweets, err := convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets)

	if err != nil {
		return nil, appErrors.UnexpectedError
	}

	return APITweets, nil
}

func GetTweet(tweetID int64) (APIModel.Tweet, *appErrors.AppError) {
	databaseTweet, err := database.GetTweet(tweetID)

	if err != nil {
		// Later on we'll need to add type switch here to check the type of error, because several things
		// can go wrong when fetching data from database: not found, SQL error, db connection error etc
		return APIModel.Tweet{}, &appErrors.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("User with given ID was not found."),
		}
	}

	APITweet, err2 := convertDatabaseTweetToAPITweet(databaseTweet)

	if err2 != nil {
		return APIModel.Tweet{}, err2
	}

	return APITweet, nil
}

func PostTweet(newTweet APIModel.NewTweet) (APIModel.Tweet, *appErrors.AppError) {
	databaseTweet := convertAPINewTweetToDatabaseTweet(newTweet)

	addedTweet, err := database.InsertTweet(databaseTweet)

	if err != nil {
		// for now its an unexpected error, but later on we'll probably need an error type switch here too
		return APIModel.Tweet{}, appErrors.UnexpectedError
	}

	APITweet, err2 := convertDatabaseTweetToAPITweet(addedTweet)

	if err2 != nil {
		return APIModel.Tweet{}, err2
	}

	return APITweet, nil
}

func GetUsers() ([]APIModel.User, *appErrors.AppError) {
	databaseUsers, err := database.GetUsers()

	if err != nil {
		// for now its an unexpected error, but later on we'll probably need an error type switch here too
		return nil, appErrors.UnexpectedError
	}

	APIUsers := convertArrayOfDatabaseUsersToArrayOfAPIUsers(databaseUsers)

	return APIUsers, nil
}

func GetUser(userId int64) (APIModel.User, *appErrors.AppError) {
	databaseUser, err := database.GetUser(userId)

	if err != nil {
		// Maybe later on we'll need to add type switch here to check the type of error, because several things
		// can go wrong when fetching data from database: not found, SQL error, db connection error etc
		return APIModel.User{}, &appErrors.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("User with given ID was not found."),
		}
	}

	APIUser := convertDatabaseUserToAPIUser(databaseUser)

	return APIUser, nil
}

func PostUser(user APIModel.NewUser) (APIModel.User, *appErrors.AppError) {
	ok := validatePostUserParameters(user)

	if !ok {
		return APIModel.User{}, &appErrors.AppError{
			Code: http.StatusBadRequest,
			Err:  errors.New("Invalid request parameters"),
		}
	}

	databaseUser := covertAPINewUserToDatabaseUser(user)

	newUser, err := database.InsertUser(databaseUser)

	if err != nil {
		// again, one error only for now...
		return APIModel.User{}, &appErrors.AppError{
			Code: http.StatusConflict,
			Err:  errors.New("User with given username already exists."),
		}
	}

	APIUser := convertDatabaseUserToAPIUser(newUser)

	return APIUser, nil
}

func validatePostUserParameters(user APIModel.NewUser) bool {
	isOk := true

	if user.Name == "" {
		isOk = false
	}
	if user.Username == "" {
		isOk = false
	}
	if user.Email == "" {
		isOk = false
	}

	return isOk
}

func convertDatabaseTweetToAPITweet(tweet databaseModel.Tweet) (APIModel.Tweet, *appErrors.AppError) {
	id := tweet.ID
	userID := tweet.AuthorID
	likeCount := tweet.LikeCount
	retweetCount := tweet.RetweetCount
	createdAt := tweet.CreatedAt
	content := tweet.Content

	authorFullData, err := database.GetUser(userID)

	if err != nil {
		// log this instead and return an error with proper message
		// errorMessage := fmt.Sprintf("no integrity in database, "+
		//	"user with id = %d was not found (but should have been found)",
		//	userID)
		return APIModel.Tweet{}, appErrors.UnexpectedError
	}

	APIAuthorFullData := convertDatabaseUserToAPIUser(authorFullData)

	APITweet := APIModel.Tweet{
		ID:           id,
		Author:       APIAuthorFullData,
		LikeCount:    likeCount,
		RetweetCount: retweetCount,
		CreatedAt:    createdAt,
		Content:      content,
	}
	return APITweet, nil
}

func convertAPINewTweetToDatabaseTweet(tweet APIModel.NewTweet) databaseModel.Tweet {
	authorId := tweet.AuthorID
	content := tweet.Content

	return databaseModel.Tweet{
		ID:           0,
		AuthorID:     authorId,
		LikeCount:    0,
		RetweetCount: 0,
		CreatedAt:    time.Now(),
		Content:      content,
	}
}

func convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets []databaseModel.Tweet) ([]APIModel.Tweet, *appErrors.AppError) {
	var APITweets []APIModel.Tweet

	for _, databaseTweet := range databaseTweets {
		APITweet, err := convertDatabaseTweetToAPITweet(databaseTweet)

		if err != nil {
			return nil, appErrors.UnexpectedError
		}

		APITweets = append(APITweets, APITweet)
	}

	return APITweets, nil
}

func convertArrayOfDatabaseUsersToArrayOfAPIUsers(databaseUsers []databaseModel.User) []APIModel.User {
	var convertedUsers []APIModel.User

	for _, databaseUser := range databaseUsers {
		APIUser := convertDatabaseUserToAPIUser(databaseUser)
		convertedUsers = append(convertedUsers, APIUser)
	}

	return convertedUsers
}

func convertDatabaseUserToAPIUser(user databaseModel.User) APIModel.User {
	id := user.ID
	name := user.Name
	username := user.Username
	email := user.Email
	createdAt := user.CreatedAt

	return APIModel.User{
		ID:        id,
		Name:      name,
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
	}
}

func covertAPINewUserToDatabaseUser(user APIModel.NewUser) databaseModel.User {
	name := user.Name
	username := user.Username
	email := user.Email

	return databaseModel.User{
		ID:        0,
		Name:      name,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}
}
