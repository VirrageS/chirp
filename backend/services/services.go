package services

import (
	"fmt"
	"errors"
	"time"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/apimodel"
)

func GetTweets() ([]apimodel.Tweet, error) {
	database_tweets := database.GetTweets()
	var api_tweets []apimodel.Tweet

	for _, database_tweet := range database_tweets {
		api_tweet, error := convertDatabaseTweetToApiTweet(database_tweet)

		if error != nil {
			return nil, errors.New("Oooops, something went wrong!")
		}

		api_tweets = append(api_tweets, api_tweet)
	}

	return api_tweets, nil
}

func GetTweet(tweet_id int64) (apimodel.Tweet, error) {
	databse_tweet, error := database.GetTweet(tweet_id)

	if error == nil {
		return apimodel.Tweet{}, error
	}

	api_tweet, error := convertDatabaseTweetToApiTweet(databse_tweet)

	if error != nil {
		return apimodel.Tweet{}, errors.New("Oooops, something went wrong!")
	}

	return api_tweet, nil
}

func PostTweet(new_tweet apimodel.NewTweet) (apimodel.Tweet, error) {
	database_tweet := convertApiNewTweetToDatabaseTweet(new_tweet)

	added_tweet := database.InsertTweet(database_tweet)

	api_tweet, error := convertDatabaseTweetToApiTweet(added_tweet)

	if error != nil {
		return apimodel.Tweet{}, errors.New("Oooops, something went wrong!")
	}

	return api_tweet, nil
}

func convertDatabaseTweetToApiTweet(tweet database.Tweet) (apimodel.Tweet, error) {
	id := tweet.Id
	user_id := tweet.AuthorId
	like_count := tweet.LikeCount
	retweet_count := tweet.RetweetCount
	created_at := tweet.CreatedAt
	content := tweet.Content

	author_full_data, error := database.GetUser(user_id)

	if error != nil {
		error_message := fmt.Sprintf("No integrity in database, " +
			"user with id = %d was not found (but should have been found)",
			user_id)
		return apimodel.Tweet{}, errors.New(error_message)
	}

	api_author_full_data := convertDatabaseUserToApiUser(author_full_data)

	api_tweet := apimodel.Tweet{
		Id: id,
		Author: api_author_full_data,
		LikeCount: like_count,
		RetweetCount: retweet_count,
		CreatedAt: created_at,
		Content: content,
	}
	return api_tweet, nil
}

func convertApiNewTweetToDatabaseTweet(tweet apimodel.NewTweet) database.Tweet {
	author_id := tweet.AuthorId
	content := tweet.Content

	return database.Tweet{
		Id: 0,
		AuthorId: author_id,
		LikeCount: 0,
		RetweetCount: 0,
		CreatedAt: time.Now(),
		Content: content,
	}
}

func convertArrayOfDatabaseUsersToArrayOfApiUsers(database_users []database.User) []apimodel.User {
	var converted_users []apimodel.User

	for _, database_user := range database_users {
		api_user := convertDatabaseUserToApiUser(database_user)
		converted_users = append(converted_users, api_user)
	}

	return converted_users
}

func convertDatabaseUserToApiUser(user database.User) apimodel.User {
	id := user.Id
	name := user.Name
	username := user.Username
	email := user.Email
	created_at := user.CreatedAt

	return apimodel.User{
		Id: id,
		Name: name,
		Username: username,
		Email: email,
		CreatedAt: created_at,
	}
}

func covertApiUserToDatabaseUser(user apimodel.User) database.User {
	id := user.Id
	name := user.Name
	username := user.Username
	email := user.Email
	created_at := user.CreatedAt

	return database.User{
		Id: id,
		Name: name,
		Username: username,
		Email: email,
		CreatedAt: created_at,
	}
}