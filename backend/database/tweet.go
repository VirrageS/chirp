package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/model"
	"github.com/lib/pq"
)

// Tweet Data Access Object - provides operations on Tweet database table
type TweetDAO interface {
	GetTweetsIDsOfUserWithID(userID int64) ([]int64, error)
	GetTweetsFromListOfIDs(tweetsToFindIDs []int64) ([]*model.Tweet, error)
	GetTweetByID(tweetID int64) (*model.Tweet, error)
	InsertTweet(newTweet *model.NewTweet) (*model.Tweet, error)
	DeleteTweet(tweetID int64) error
}

type tweetDB struct {
	*sql.DB
}

func NewTweetDAO(dbConnection *sql.DB) TweetDAO {
	return &tweetDB{dbConnection}
}

func (db *tweetDB) GetTweetsIDsOfUserWithID(userID int64) ([]int64, error) {
	rows, err := db.Query(`
		SELECT id
		FROM tweets
		WHERE author_id = $1
		ORDER BY created_at DESC`,
		userID)
	if err != nil {
		log.WithField("userID", userID).WithError(err).Error("GetTweetsIDsOfUserWithID query error.")
		return nil, err
	}
	defer rows.Close()

	tweetIDs, err := readMultipleTweetsIDs(rows)
	if err != nil {
		log.WithError(err).Error("GetTweetsIDsOfUserWithID rows scan/iteration error.")
	}

	return tweetIDs, nil
}

func (db *tweetDB) GetTweetsFromListOfIDs(tweetsToFindIDs []int64) ([]*model.Tweet, error) {
	rows, err := db.Query(`
		SELECT id, created_at, content, author_id
		FROM tweets
		WHERE id = ANY($1)
		ORDER BY created_at DESC`,
		pq.Array(tweetsToFindIDs))
	if err != nil {
		log.WithField("tweetsToFindIDs", tweetsToFindIDs).WithError(err).Error("GetTweetsFromListOfIDs query error.")
		return nil, err
	}
	defer rows.Close()

	tweets, err := readMultipleTweets(rows)
	if err != nil {
		log.WithError(err).Error("GetTweetsFromListOfIDs rows scan/iteration error.")
	}

	return tweets, nil
}

func (db *tweetDB) GetTweetByID(tweetID int64) (*model.Tweet, error) {
	row := db.QueryRow(`
		SELECT id, created_at, content, author_id
		FROM tweets
		WHERE id = $1
		ORDER BY created_at DESC`,
		tweetID)

	tweet, err := readTweet(row)
	if err != nil && err != sql.ErrNoRows {
		log.WithField("tweetID", tweetID).WithError(err).Error("GetTweetWithID query error.")
		return nil, err
	}

	return tweet, err
}

func (db *tweetDB) InsertTweet(newTweet *model.NewTweet) (*model.Tweet, error) {
	row := db.QueryRow(`
		INSERT INTO tweets (author_id, content)
		VALUES ($1, $2)
		RETURNING id, created_at, content, author_id`,
		newTweet.AuthorID, newTweet.Content)

	insertedTweet, err := readTweet(row)
	if err != nil {
		log.WithField("newTweet", *newTweet).WithError(err).Error("InsertTweet query error.")
		return nil, err
	}

	return insertedTweet, nil
}

func (db *tweetDB) DeleteTweet(tweetID int64) error {
	_, err := db.Exec(`DELETE FROM tweets WHERE id=$1`, tweetID)
	if err != nil {
		log.WithField("tweetID", tweetID).WithError(err).Error("DeleteTweet query error.")
		return err
	}

	return nil
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
