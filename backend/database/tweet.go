package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
	"github.com/lib/pq"
)

// Tweet Data Access Object - provides operations on Tweet database table
type TweetDAO interface {
	GetTweetsIDsByAuthorID(userID int64) ([]int64, error)
	GetTweetsByIDs(tweetsIDs []int64) ([]*model.Tweet, error)
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

func (db *tweetDB) GetTweetsIDsByAuthorID(userID int64) ([]int64, error) {
	rows, err := db.Query(`
		SELECT id
		FROM tweets
		WHERE author_id = $1
		ORDER BY created_at DESC`,
		userID)
	if err != nil {
		log.WithField("userID", userID).WithError(err).Error("GetTweetsIDsByAuthorID query error.")
		return nil, err
	}
	defer rows.Close()

	tweetIDs, err := readMultipleTweetsIDs(rows)
	if err != nil {
		log.WithError(err).Error("GetTweetsIDsByAuthorID rows scan/iteration error.")
	}

	return tweetIDs, nil
}

func (db *tweetDB) GetTweetsByIDs(tweetsIDs []int64) ([]*model.Tweet, error) {
	rows, err := db.Query(`
		SELECT id, created_at, content, author_id
		FROM tweets
		WHERE id = ANY($1)
		ORDER BY created_at DESC`,
		pq.Array(tweetsIDs))
	if err != nil {
		log.WithField("tweetsIDs", tweetsIDs).WithError(err).Error("GetTweetsByIDs query error.")
		return nil, err
	}
	defer rows.Close()

	tweets, err := readMultipleTweets(rows)
	if err != nil {
		log.WithError(err).Error("GetTweetsByIDs rows scan/iteration error.")
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
	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}
	if err != nil {
		log.WithField("tweetID", tweetID).WithError(err).Error("GetTweetByID query error.")
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
