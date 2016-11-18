package database

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/VirrageS/chirp/backend/database/model"
)

type TweetDataAccessor interface {
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
	tweets, err := db.getTweets()
	if err != nil {
		return nil, DatabaseError
	}

	return tweets, nil
}

func (db *TweetDB) GetTweetsOfUserWithID(userID int64) ([]model.Tweet, error) {
	tweets, err := db.getTweetsOfUserWithID(userID)
	if err != nil {
		return nil, DatabaseError
	}

	return tweets, nil
}

func (db *TweetDB) GetTweet(tweetID int64) (model.Tweet, error) {
	tweet, err := db.getTweetUsingQuery("SELECT * FROM tweets WHERE id=$1;", tweetID)
	if err == sql.ErrNoRows {
		return model.Tweet{}, NoRowsError
	}
	if err != nil {
		return model.Tweet{}, DatabaseError
	}

	return tweet, nil
}

func (db *TweetDB) InsertTweet(tweet model.Tweet) (model.Tweet, error) {
	tweetID, err := db.insertTweetToDatabase(tweet)
	if err != nil {
		return model.Tweet{}, DatabaseError
	}

	tweet.ID = tweetID

	return tweet, nil
}

func (db *TweetDB) DeleteTweet(tweetID int64) error {
	err := db.deleteTweetWithID(tweetID)
	if err != nil {
		return DatabaseError
	}

	return nil
}

// TODO: Maybe it should also fetch tweet's User and embed it inside the returned object
func (db *TweetDB) getTweetUsingQuery(query string, args ...interface{}) (model.Tweet, error) {
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

func (db *TweetDB) insertTweetToDatabase(tweet model.Tweet) (int64, error) {
	query, err := db.Prepare("INSERT INTO tweets (author_id, created_at, content) " +
		"VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		log.WithField("query", query).WithError(err).Error("insertTweetToDatabase query prepare error.")
		return 0, err
	}
	defer query.Close()

	var newID int64
	// for Postgres we need to use query with RETURNING id to get the ID of the inserted tweet
	err = query.QueryRow(tweet.AuthorID, tweet.CreatedAt, tweet.Content).Scan(&newID)
	if err != nil {
		log.WithError(err).Error("insertTweetToDatabase query execute error.")
		return 0, err
	}

	return newID, nil
}

func (db *TweetDB) deleteTweetWithID(tweetID int64) error {
	statement, err := db.Prepare("DELETE FROM tweets WHERE id=$1")
	if err != nil {
		log.WithField("query", statement).WithError(err).Error("deleteTweetWithID query prepare error.")
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(tweetID)
	if err != nil {
		log.WithError(err).Error("deleteTweetWithID query execute error.")
		return err
	}

	return nil
}

func (db *TweetDB) getTweets() ([]model.Tweet, error) {
	rows, err := db.Query("SELECT * FROM tweets;")
	if err != nil {
		log.WithError(err).Error("GetTweets query error.")
	}

	var tweets []model.Tweet

	defer rows.Close()
	for rows.Next() {
		var tweet model.Tweet

		err := rows.Scan(&tweet.ID, &tweet.AuthorID, &tweet.CreatedAt, &tweet.Content)
		if err != nil {
			log.WithError(err).Error("getTweets row scan error.")
			return nil, err
		}

		tweets = append(tweets, tweet)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("getTweets rows iteration error.")
		return nil, err
	}

	return tweets, nil
}

// TODO: almost the same as getTweets()...
func (db *TweetDB) getTweetsOfUserWithID(userID int64) ([]model.Tweet, error) {
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
			return nil, err
		}

		tweets = append(tweets, tweet)
	}

	return tweets, nil
}
