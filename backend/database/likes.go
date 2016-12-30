package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
)

// Likes Data Access Object - provides operations on Likes database table
type LikesDAO interface {
	LikeTweet(tweetID, userID int64) error
	UnlikeTweet(tweetID, userID int64) error
	LikeCount(tweetID int64) (int64, error)
	IsLiked(tweetID, userID int64) (bool, error)
}

type likesDB struct {
	*sql.DB
}

func NewLikesDAO(dbConnection *sql.DB) LikesDAO {
	return &likesDB{dbConnection}
}

func (db *likesDB) LikeTweet(tweetID, userID int64) error {
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

func (db *likesDB) UnlikeTweet(tweetID, userID int64) error {
	_, err := db.Exec(`DELETE FROM likes WHERE tweet_id=$1 AND user_id=$2`, tweetID, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"tweetID": tweetID,
			"userID":  userID,
		}).WithError(err).Error("DeleteTweet query error.")
		return err
	}

	return nil
}

func (db *likesDB) LikeCount(tweetID int64) (int64, error) {
	var likeCount int64

	err := db.QueryRow(`SELECT COUNT(*) FROM likes WHERE tweet_id = $1`, tweetID).Scan(&likeCount)
	if err != nil {
		log.WithField("tweetID", tweetID).WithError(err).Error("LikeCount query error.")
		return 0, err
	}

	return likeCount, nil
}

func (db *likesDB) IsLiked(tweetID, userID int64) (bool, error) {
	var isLiked bool

	err := db.QueryRow(`SELECT exists (SELECT TRUE FROM likes WHERE tweet_id = $1 AND user_id = $2)`,
		tweetID, userID).Scan(&isLiked)
	if err != nil {
		log.WithFields(log.Fields{
			"tweetID": tweetID,
			"userID":  userID,
		}).WithError(err).Error("IsLiked query error.")
		return false, err
	}

	return isLiked, nil
}
