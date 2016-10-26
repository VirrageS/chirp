package services

import (
	"errors"
	"fmt"

	APIModel "github.com/VirrageS/chirp/backend/api/model"
	"github.com/VirrageS/chirp/backend/database"
	databaseModel "github.com/VirrageS/chirp/backend/database/model"

	"time"
)

func GetTweets() ([]APIModel.Tweet, error) {
	databaseTweets, err := database.GetTweets()

	if err != nil {
		return nil, errors.New("Oooops, something went wrong!")
	}

	APITweets, err := convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets)

	if err != nil {
		return nil, err
	}

	return APITweets, nil
}

func GetTweet(tweetID int64) (APIModel.Tweet, error) {
	databaseTweet, err := database.GetTweet(tweetID)

	if err != nil {
		return APIModel.Tweet{}, err
	}

	APITweet, err := convertDatabaseTweetToAPITweet(databaseTweet)

	if err != nil {
		return APIModel.Tweet{}, errors.New("Oooops, something went wrong!")
	}

	return APITweet, nil
}

func PostTweet(newTweet APIModel.NewTweet) (APIModel.Tweet, error) {
	databaseTweet := convertAPINewTweetToDatabaseTweet(newTweet)

	addedTweet := database.InsertTweet(databaseTweet)

	APITweet, err := convertDatabaseTweetToAPITweet(addedTweet)

	if err != nil {
		return APIModel.Tweet{}, errors.New("Oooops, something went wrong!")
	}

	return APITweet, nil
}

func GetUsers() ([]APIModel.User, error) {
	databaseUsers, err := database.GetUsers()

	if err != nil {
		return nil, errors.New("Oooops, something went wrong!")
	}

	APIUsers := convertArrayOfDatabaseUsersToArrayOfAPIUsers(databaseUsers)

	return APIUsers, nil
}

func GetUser(userId int64) (APIModel.User, error) {
	databaseUser, err := database.GetUser(userId)

	if err != nil {
		return APIModel.User{}, err
	}

	APIUser := convertDatabaseUserToAPIUser(databaseUser)

	return APIUser, nil
}

func PostUser(user APIModel.NewUser) (APIModel.User, error) {
	ok, errs := validatePostUserParameters(user)

	if !ok {
		return APIModel.User{}, errs[0] // TODO: replace with custom error type that can hold array of messages
	}

	databaseUser := covertAPINewUserToDatabaseUser(user)

	newUser, err := database.InsertUser(databaseUser)

	if err != nil {
		return APIModel.User{}, err
	}

	APIUser := convertDatabaseUserToAPIUser(newUser)

	return APIUser, nil
}

func validatePostUserParameters(user APIModel.NewUser) (bool, []error) {
	var errs []error
	isOk := true

	if user.Name == "" {
		isOk = false
		errs = append(errs, errors.New("user name cannot be empty"))
	}
	if user.Username == "" {
		isOk = false
		errs = append(errs, errors.New("user username cannot be empty"))
	}
	if user.Email == "" {
		isOk = false
		errs = append(errs, errors.New("user email cannot be empty"))
	}

	return isOk, errs
}

func convertDatabaseTweetToAPITweet(tweet databaseModel.Tweet) (APIModel.Tweet, error) {
	id := tweet.ID
	userID := tweet.AuthorID
	likeCount := tweet.LikeCount
	retweetCount := tweet.RetweetCount
	createdAt := tweet.CreatedAt
	content := tweet.Content

	authorFullData, err := database.GetUser(userID)

	if err != nil {
		errorMessage := fmt.Sprintf("no integrity in database, "+
			"user with id = %d was not found (but should have been found)",
			userID)
		return APIModel.Tweet{}, errors.New(errorMessage)
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

func convertArrayOfDatabaseTweetsToArrayOfAPITweets(databaseTweets []databaseModel.Tweet) ([]APIModel.Tweet, error) {
	var APITweets []APIModel.Tweet

	for _, databaseTweet := range databaseTweets {
		APITweet, err := convertDatabaseTweetToAPITweet(databaseTweet)

		if err != nil {
			return nil, errors.New("Oooops, something went wrong!")
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
