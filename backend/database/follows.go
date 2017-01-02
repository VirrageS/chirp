package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
)

// Follows Data Access Object - provides operations on Follows database table
type FollowsDAO interface {
	FollowUser(followeeID, followerID int64) error
	UnfollowUser(followeeID, followerID int64) error
	GetFollowersIDs(userID int64) ([]int64, error)
	GetFolloweesIDs(userID int64) ([]int64, error)
	GetFollowerCount(userID int64) (int64, error)
	GetFolloweeCount(userID int64) (int64, error)
	IsFollowing(followerID, followeeID int64) (bool, error)
}

type followsDB struct {
	*sql.DB
}

func NewFollowsDAO(dbConnection *sql.DB) FollowsDAO {
	return &followsDB{dbConnection}
}

func (db *followsDB) FollowUser(followeeID, followerID int64) error {
	_, err := db.Exec(`
		INSERT INTO follows (followee_id, follower_id)
		VALUES ($1, $2)
		ON CONFLICT (followee_id, follower_id) DO NOTHING`,
		followeeID, followerID,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"followeeID": followeeID,
			"followerID": followerID,
		}).WithError(err).Error("FollowUser query error.")
		return err
	}

	return nil
}

func (db *followsDB) UnfollowUser(followeeID, followerID int64) error {
	_, err := db.Exec(`
		DELETE FROM follows
		WHERE followee_id=$1 AND follower_id=$2;
		`,
		followeeID, followerID,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"followeeID": followeeID,
			"followerID": followerID,
		}).WithError(err).Error("UnfollowUser query error.")
		return err
	}

	return nil
}

func (db *followsDB) GetFollowersIDs(userID int64) ([]int64, error) {
	rows, err := db.Query(`SELECT follower_id FROM follows WHERE followee_id = $1`, userID)
	if err != nil {
		log.WithError(err).Error("GetFollowersIDs query error")
	}
	defer rows.Close()

	followersIDs := make([]int64, 0)
	for rows.Next() {
		var followerID int64
		err = rows.Scan(&followerID)
		if err != nil {
			log.WithError(err).Error("GetFollowersIDs row scan error.")
			return nil, err
		}

		followersIDs = append(followersIDs, followerID)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("GetFollowersIDs rows iteration error.")
		return nil, err
	}

	return followersIDs, nil
}

func (db *followsDB) GetFolloweesIDs(userID int64) ([]int64, error) {
	rows, err := db.Query(`SELECT followee_id FROM follows WHERE follower_id = $1`, userID)
	if err != nil {
		log.WithError(err).Error("GetFolloweesIDs query error")
	}
	defer rows.Close()

	followeesIDs := make([]int64, 0)
	for rows.Next() {
		var followeeID int64

		err = rows.Scan(&followeeID)
		if err != nil {
			log.WithError(err).Error("GetFolloweesIDs row scan error.")
			return nil, err
		}

		followeesIDs = append(followeesIDs, followeeID)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("GetFolloweesIDs rows iteration error.")
		return nil, err
	}

	return followeesIDs, nil
}

func (db *followsDB) GetFollowerCount(userID int64) (int64, error) {
	var followerCount int64

	err := db.QueryRow(`SELECT COUNT(*) FROM follows WHERE followee_id = $1`, userID).Scan(&followerCount)
	if err != nil {
		log.WithError(err).Error("GetFollowerCount query error.")
		return 0, err
	}

	return followerCount, nil
}

func (db *followsDB) GetFolloweeCount(userID int64) (int64, error) {
	var followeeCount int64

	err := db.QueryRow(`SELECT COUNT(*) FROM follows WHERE follower_id = $1`, userID).Scan(&followeeCount)
	if err != nil {
		log.WithError(err).Error("GetFolloweeCount query error.")
		return 0, err
	}

	return followeeCount, nil
}

func (db *followsDB) IsFollowing(followerID, followeeID int64) (bool, error) {
	var isFollowing bool

	err := db.QueryRow(`SELECT exists (SELECT TRUE FROM follows WHERE follower_id = $1 AND followee_id = $2)`,
		followerID, followeeID).Scan(&isFollowing)
	if err != nil {
		log.WithError(err).Error("IsFollowing query error.")
		return false, err
	}

	return isFollowing, nil
}
