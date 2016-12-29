package database

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/VirrageS/chirp/backend/model"
)

type TweetDAO interface {
	GetTweets() ([]*model.Tweet, error)
	GetTweetsOfUserWithID(userID int64) ([]*model.Tweet, error)
	GetTweetWithID(tweetID int64) (*model.Tweet, error)
	InsertTweet(newTweet *model.NewTweet) (*model.Tweet, error)
	DeleteTweet(tweetID int64) error

	LikeTweet(tweetID, userID int64) error
	UnlikeTweet(tweetID, userID int64) error
	LikeCount(tweetID int64) (int64, error)
	IsLiked(tweetID, userID int64) (bool, error)
}

type tweetDB struct {
	*sql.DB
}

func NewTweetDAO(dbConnection *sql.DB) TweetDAO {
	return &tweetDB{dbConnection}
}

func (db *tweetDB) GetTweets() ([]*model.Tweet, error) {
	rows, err := db.Query(`
		SELECT id, created_at, content, author_id
		FROM tweets
		ORDER BY tweets.created_at DESC`)
	if err != nil {
		log.WithError(err).Error("GetTweets query error.")
		return nil, err
	}
	defer rows.Close()

	tweets, err := readMultipleTweets(rows)
	if err != nil {
		log.WithError(err).Error("GetTweets rows scan/iteration error.")
	}

	return tweets, nil
}

func (db *tweetDB) GetTweetsOfUserWithID(userID int64) ([]*model.Tweet, error) {
	rows, err := db.Query(`
		SELECT id, created_at, content, author_id
		FROM tweets
		WHERE author_id = $1
		ORDER BY created_at DESC`,
		userID)
	if err != nil {
		log.WithField("userID", userID).WithError(err).Error("GetTweetsOfUserWithID query error.")
		return nil, err
	}
	defer rows.Close()

	tweets, err := readMultipleTweets(rows)
	if err != nil {
		log.WithError(err).Error("GetTweetsOfUserWithID rows scan/iteration error.")
	}

	return tweets, nil
}

func (db *tweetDB) GetTweetWithID(tweetID int64) (*model.Tweet, error) {
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
	_, err := db.Exec(`
		DELETE FROM tweets
		WHERE id=$1`,
		tweetID)
	if err != nil {
		log.WithField("tweetID", tweetID).WithError(err).Error("DeleteTweet query error.")
		return err
	}

	return nil
}

func (db *tweetDB) LikeTweet(tweetID, userID int64) error {
	_, err := db.Exec(`
		INSERT INTO likes (tweet_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (tweet_id, user_id) DO NOTHING`,
		tweetID, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"tweetID": tweetID,
			"userID":  userID,
		}).WithError(err).Error("LikeTweet query error.")
		return err
	}

	return nil
}

func (db *tweetDB) UnlikeTweet(tweetID, userID int64) error {
	_, err := db.Exec(`
		DELETE FROM likes
		WHERE tweet_id=$1 AND user_id=$2`,
		tweetID, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"tweetID": tweetID,
			"userID":  userID,
		}).WithError(err).Error("DeleteTweet query error.")
		return err
	}

	return nil
}

func (db *tweetDB) LikeCount(tweetID int64) (int64, error) {
	var likeCount int64

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM likes
		WHERE tweet_id = $1`,
		tweetID).
		Scan(&likeCount)
	if err != nil {
		log.WithField("tweetID", tweetID).WithError(err).Error("LikeCount query error.")
		return 0, err
	}

	return likeCount, nil
}

func (db *tweetDB) IsLiked(tweetID, userID int64) (bool, error) {
	var isLiked bool

	err := db.QueryRow(`
		SELECT exists
			(SELECT true
			FROM likes
			WHERE tweet_id = $1 AND user_id = $2)`,
		tweetID, userID).
		Scan(&isLiked)
	if err != nil {
		log.WithFields(log.Fields{
			"tweetID": tweetID,
			"userID":  userID,
		}).WithError(err).Error("IsLiked query error.")
		return false, err
	}

	return isLiked, nil
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
