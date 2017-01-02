package database

import (
	"database/sql"

	"github.com/VirrageS/chirp/backend/model"
)

// Helper that wraps rows and row so they can be used in the same function
type scannable interface {
	Scan(dest ...interface{}) error
}

func readMultipleUsers(rows *sql.Rows) ([]*model.PublicUser, error) {
	users := make([]*model.PublicUser, 0)

	for rows.Next() {
		user, err := readPublicUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func readPublicUser(row scannable) (*model.PublicUser, error) {
	var user model.PublicUser

	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func readUser(row scannable) (*model.User, error) {
	var user model.User

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Name,
		&user.TwitterToken, &user.FacebookToken, &user.GoogleToken,
		&user.CreatedAt, &user.LastLogin, &user.Active, &user.AvatarUrl)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func readMultipleTweetsIDs(rows *sql.Rows) ([]int64, error) {
	tweetsIDs := make([]int64, 0)
	for rows.Next() {
		var tweetID int64

		err := rows.Scan(&tweetID)
		if err != nil {
			return nil, err
		}

		tweetsIDs = append(tweetsIDs, tweetID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tweetsIDs, nil
}

func readMultipleTweets(rows *sql.Rows) ([]*model.Tweet, error) {
	tweets := make([]*model.Tweet, 0)
	for rows.Next() {
		tweet, err := readTweet(rows)
		if err != nil {
			return nil, err
		}

		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tweets, nil
}

func readTweet(row scannable) (*model.Tweet, error) {
	var tweet model.Tweet
	var authorID int64

	err := row.Scan(&tweet.ID, &tweet.CreatedAt, &tweet.Content, &authorID)
	if err != nil {
		return nil, err
	}

	tweet.Author = &model.PublicUser{ID: authorID}

	return &tweet, nil
}
