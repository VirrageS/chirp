package database

import (
	"errors"

	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/VirrageS/chirp/backend/database/model"
)

type TweetDataAcessor interface {
	GetTweets() ([]model.Tweet, error)
	GetTweetsOfUserWithID(userID int64) ([]model.Tweet, error)
	GetTweet(tweetID int64) (model.Tweet, error)
	InsertTweet(tweet model.Tweet) (model.Tweet, error)
	DeleteTweet(tweetID int64) error
}

type TweetDB struct {
	*sql.DB
}

func NewTweetDB(databaseConnection *sql.DB) *TweetDB {
	return &TweetDB{databaseConnection}
}

func (db *TweetDB) GetTweets() ([]model.Tweet, error) {
	return getTweets(db)
}

func (db *TweetDB) GetTweetsOfUserWithID(userID int64) ([]model.Tweet, error) {
	return getTweetsOfUserWithID(db, userID)
}

func (db *TweetDB) GetTweet(tweetID int64) (model.Tweet, error) {
	tweet, err := getTweetUsingQuery(db, "SELECT * FROM tweets WHERE id=$1;", tweetID)
	if err == sql.ErrNoRows {
		return model.Tweet{}, errors.New("") // no users found error
	}
	if err != nil {
		return model.Tweet{}, errors.New("") // db error
	}

	return tweet, nil
}

func (db *TweetDB) InsertTweet(tweet model.Tweet) (model.Tweet, error) {
	tweetID, err := insertTweetToDatabase(db, tweet)
	if err != nil {
		return model.Tweet{}, errors.New("") // db error
	}

	tweet.ID = tweetID

	return tweet, nil
}

func (db *TweetDB) DeleteTweet(tweetID int64) error {
	err := deleteTweetWithID(db, tweetID)
	if err != nil {
		return errors.New("")
	}

	return nil
}

func getTweetUsingQuery(db *TweetDB, query string, args ...interface{}) (model.Tweet, error) {
	var tweet model.Tweet

	row := db.QueryRow(query, args...)
	err := row.Scan(&tweet.ID, &tweet.AuthorID, &tweet.CreatedAt, &tweet.Content)

	if err != nil && err != sql.ErrNoRows {
		log.WithFields(log.Fields{
			"error": err,
			"query": query,
			"args":  args,
		}).Error("GetTweetUsingQuery database error.")
	}

	return tweet, err
}

func insertTweetToDatabase(db *TweetDB, tweet model.Tweet) (int64, error) {
	query, err := db.Prepare("INSERT INTO tweets (author_id, created_at, content) " +
		"VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		log.WithField("query", query).WithError(err).Error("insertTweetToDatabase query prepare error.")
		return 0, errors.New("")
	}
	defer query.Close()

	var newID int64
	// for Postgres we need to use query with RETURNING id to get the ID of the inserted tweet
	err = query.QueryRow(tweet.AuthorID, tweet.CreatedAt, tweet.Content).Scan(&newID)
	if err != nil {
		log.WithError(err).Error("insertTweetToDatabase query execute error.")
		return 0, errors.New("")
	}

	return newID, nil
}

func deleteTweetWithID(db *TweetDB, tweetID int64) error {
	statement, err := db.Prepare("DELETE FROM tweets WHERE id=$1")
	if err != nil {
		log.WithField("query", statement).WithError(err).Error("deleteTweetWithID query prepare error.")
		errors.New("")
	}
	defer statement.Close()

	_, err = statement.Exec(tweetID)
	if err != nil {
		log.WithError(err).Error("deleteTweetWithID query execute error.")
		return errors.New("")
	}

	return nil
}

func getTweets(db *TweetDB) ([]model.Tweet, error) {
	rows, err := db.Query("SELECT * FROM tweets;")
	if err != nil {
		log.WithError(err).Error("GetTweets query error.")
	}

	var tweets []model.Tweet

	defer rows.Close()
	for rows.Next() {
		var tweet model.Tweet
		err := rows.Scan(&tweet.ID, &tweet.AuthorID, &tweet.CreatedAt, &tweet.Content)
		// TODO: move error outside of the loop
		if err != nil {
			log.WithError(err).Error("GetTweets row scan error.")
			return nil, errors.New("")
		}

		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

// TODO: almost the same as getTweets()...
func getTweetsOfUserWithID(db *TweetDB, userID int64) ([]model.Tweet, error) {
	rows, err := db.Query("SELECT * FROM tweets WHERE id=$1;", userID)
	if err != nil {
		log.WithError(err).Error("GetTweets query error.")
	}

	var tweets []model.Tweet

	defer rows.Close()
	for rows.Next() {
		var tweet model.Tweet
		err := rows.Scan(&tweet.ID, &tweet.AuthorID, &tweet.CreatedAt, &tweet.Content)
		// TODO: move error outside of the loop
		if err != nil {
			log.WithError(err).Error("GetTweets row scan error.")
			return nil, errors.New("")
		}

		tweets = append(tweets, tweet)
	}

	return tweets, nil
}
