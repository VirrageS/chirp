package services

import (
	"errors"
	"fmt"
	"github.com/VirrageS/chirp/backend/apiModel"
	"github.com/VirrageS/chirp/backend/database"
	"time"
)

func GetTweets() ([]apiModel.Tweet, error) {
	databaseTweets, err := database.GetTweets()

	if err != nil {
		return nil, errors.New("Oooops, something went wrong!")
	}

	apiTweets, err := convertArrayOfDatabaseTweetsToArrayOfApiTweets(databaseTweets)

	if err != nil {
		return nil, err
	}

	return apiTweets, nil
}

func GetTweet(tweetId int64) (apiModel.Tweet, error) {
	databaseTweet, err := database.GetTweet(tweetId)

	if err != nil {
		return apiModel.Tweet{}, err
	}

	apiTweet, err := convertDatabaseTweetToApiTweet(databaseTweet)

	if err != nil {
		return apiModel.Tweet{}, errors.New("Oooops, something went wrong!")
	}

	return apiTweet, nil
}

func PostTweet(newTweet apiModel.NewTweet) (apiModel.Tweet, error) {
	databaseTweet := convertApiNewTweetToDatabaseTweet(newTweet)

	addedTweet := database.InsertTweet(databaseTweet)

	apiTweet, err := convertDatabaseTweetToApiTweet(addedTweet)

	if err != nil {
		return apiModel.Tweet{}, errors.New("Oooops, something went wrong!")
	}

	return apiTweet, nil
}

func GetUsers() ([]apiModel.User, error) {
	databaseUsers, err := database.GetUsers()

	if err != nil {
		return nil, errors.New("Oooops, something went wrong!")
	}

	apiUsers := convertArrayOfDatabaseUsersToArrayOfApiUsers(databaseUsers)

	return apiUsers, nil
}

func GetUser(userId int64) (apiModel.User, error) {
	databaseUser, err := database.GetUser(userId)

	if err != nil {
		return apiModel.User{}, err
	}

	apiUser := convertDatabaseUserToApiUser(databaseUser)

	return apiUser, nil
}

func PostUser(user apiModel.NewUser) (apiModel.User, error) {
	if user.Name == "" {
		return apiModel.User{}, errors.New("No name was provided.")
	}
	if user.Username == "" {
		return apiModel.User{}, errors.New("No username was provided.")
	}
	if user.Email == "" {
		return apiModel.User{}, errors.New("No email was provided.")
	}

	databaseUser := covertApiNewUserToDatabaseUser(user)

	newUser, err := database.InsertUser(databaseUser)

	if err != nil {
		return apiModel.User{}, err
	}

	apiUser := convertDatabaseUserToApiUser(newUser)

	return apiUser, nil
}

func convertDatabaseTweetToApiTweet(tweet database.Tweet) (apiModel.Tweet, error) {
	id := tweet.Id
	userId := tweet.AuthorId
	likeCount := tweet.LikeCount
	retweetCount := tweet.RetweetCount
	createdAt := tweet.CreatedAt
	content := tweet.Content

	authorFullData, err := database.GetUser(userId)

	if err != nil {
		errorMessage := fmt.Sprintf("No integrity in database, "+
			"user with id = %d was not found (but should have been found)",
			userId)
		return apiModel.Tweet{}, errors.New(errorMessage)
	}

	apiAuthorFullData := convertDatabaseUserToApiUser(authorFullData)

	apiTweet := apiModel.Tweet{
		Id:           id,
		Author:       apiAuthorFullData,
		LikeCount:    likeCount,
		RetweetCount: retweetCount,
		CreatedAt:    createdAt,
		Content:      content,
	}
	return apiTweet, nil
}

func convertApiNewTweetToDatabaseTweet(tweet apiModel.NewTweet) database.Tweet {
	authorId := tweet.AuthorId
	content := tweet.Content

	return database.Tweet{
		Id:           0,
		AuthorId:     authorId,
		LikeCount:    0,
		RetweetCount: 0,
		CreatedAt:    time.Now(),
		Content:      content,
	}
}

func convertArrayOfDatabaseTweetsToArrayOfApiTweets(databaseTweets []database.Tweet) ([]apiModel.Tweet, error) {
	var apiTweets []apiModel.Tweet

	for _, databaseTweet := range databaseTweets {
		apiTweet, err := convertDatabaseTweetToApiTweet(databaseTweet)

		if err != nil {
			return nil, errors.New("Oooops, something went wrong!")
		}

		apiTweets = append(apiTweets, apiTweet)
	}

	return apiTweets, nil
}

func convertArrayOfDatabaseUsersToArrayOfApiUsers(databaseUsers []database.User) []apiModel.User {
	var convertedUsers []apiModel.User

	for _, databaseUser := range databaseUsers {
		api_user := convertDatabaseUserToApiUser(databaseUser)
		convertedUsers = append(convertedUsers, api_user)
	}

	return convertedUsers
}

func convertDatabaseUserToApiUser(user database.User) apiModel.User {
	id := user.Id
	name := user.Name
	username := user.Username
	email := user.Email
	createdAt := user.CreatedAt

	return apiModel.User{
		Id:        id,
		Name:      name,
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
	}
}

func covertApiNewUserToDatabaseUser(user apiModel.NewUser) database.User {
	name := user.Name
	username := user.Username
	email := user.Email

	return database.User{
		Id:        0,
		Name:      name,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}
}
