package database

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/model"
)

type UserDAO interface {
	GetPublicUsers() ([]*model.PublicUser, error)
	GetPublicUsersFromListOfIDs(usersToFindIDs []int64) ([]*model.PublicUser, error)
	GetPublicUserWithID(userID int64) (*model.PublicUser, error)
	GetUserWithEmail(userEmail string) (*model.User, error)
	InsertUser(user *model.NewUserForm) (int64, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error

	FollowUser(followeeID, followerID int64) error
	UnfollowUser(followeeID, followerID int64) error
	IDsOfFollowers(userID int64) ([]int64, error)
	IDsOfFollowees(userID int64) ([]int64, error)
	FollowerCount(userID int64) (int64, error)
	FolloweeCount(userID int64) (int64, error)
	IsFollowing(followerID, followeeID int64) (bool, error)
}

type userDB struct {
	*sql.DB
}

func NewUserDAO(dbConnection *sql.DB) UserDAO {
	return &userDB{dbConnection}
}

func (db *userDB) GetPublicUsers() ([]*model.PublicUser, error) {
	rows, err := db.Query(`
		SELECT id, username, name, avatarl_url
		FROM users
		ORDER BY id DESC`)
	if err != nil {
		log.WithError(err).Error("GetPublicUsers query error.")
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.PublicUser, 0)
	for rows.Next() {
		var user model.PublicUser
		err = rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl)
		if err != nil {
			log.WithError(err).Error("GetPublicUsers row scan error.")
			return nil, err
		}

		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("GetPublicUsers rows iteration error.")
		return nil, err
	}

	return users, nil
}

func (db *userDB) GetPublicUsersFromListOfIDs(usersToFindIDs []int64) ([]*model.PublicUser, error) {
	// TODO: be careful - this ANY query is said to be super inefficient
	rows, err := db.Query(`SELECT id, username, name, avatar_url
		FROM users
		WHERE users.id = ANY($1)
		GROUP BY users.id`,
		pq.Array(usersToFindIDs))
	if err != nil {
		log.WithError(err).Error("GetPublicUsersFromListOfIDs query error.")
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.PublicUser, 0)
	for rows.Next() {
		var user model.PublicUser
		err = rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl)
		if err != nil {
			log.WithError(err).Error("GetPublicUsersFromListOfIDs row scan error.")
			return nil, err
		}

		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("GetPublicUsersFromListOfIDs rows iteration error.")
		return nil, err
	}

	return users, nil
}

func (db *userDB) GetPublicUserWithID(userID int64) (*model.PublicUser, error) {
	var user model.PublicUser

	err := db.QueryRow(`
		SELECT id, username, name, avatar_url
		FROM users
		WHERE id = $1`,
		userID).
		Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("GetPublicUserWithID database error.")
		return nil, err
	}

	return &user, err
}

func (db *userDB) GetUserWithEmail(userEmail string) (*model.User, error) {
	var user model.User

	err := db.QueryRow(`
		SELECT id, username, password, email, name,
			twitter_token, facebook_token, google_token,
			created_at, last_login, active, avatar_url
		FROM users
		WHERE email = $1`,
		userEmail).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Name,
			&user.TwitterToken, &user.FacebookToken, &user.GoogleToken,
			&user.CreatedAt, &user.LastLogin, &user.Active, &user.AvatarUrl)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("GetUserWithEmail database error.")
		return nil, err
	}

	return &user, err
}

func (db *userDB) InsertUser(user *model.NewUserForm) (int64, error) {
	var newID int64

	// for Postgres we need to use query with RETURNING id to get the ID of the inserted user
	err := db.QueryRow(`
		INSERT INTO users (username, email, password, name)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		user.Username, user.Email, user.Password, user.Name).
		Scan(&newID)
	if err != nil {
		log.WithError(err).Error("InsertUser query execute error.")
		return 0, err
	}

	return newID, nil
}

func (db *userDB) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	_, err := db.Exec(`
		UPDATE users
		SET last_login = $1
		WHERE id = $2`,
		lastLoginTime, userID,
	)
	if err != nil {
		log.WithError(err).Error("UpdateUserLastLoginTime query execute error.")
		return err
	}

	return nil
}

func (db *userDB) FollowUser(followeeID, followerID int64) error {
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
		}).WithError(err).Error("FollowUser query execute error.")
		return err
	}

	return nil
}

func (db *userDB) UnfollowUser(followeeID, followerID int64) error {
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
		}).WithError(err).Error("UnfollowUser query execute error.")
		return err
	}

	return nil
}

func (db *userDB) IDsOfFollowers(userID int64) ([]int64, error) {
	rows, err := db.Query(`
		SELECT follower_id
		FROM follows
		WHERE followee_id = $1`,
		userID)
	if err != nil {
		log.WithError(err).Error("IDsOfFollowers query error")
	}
	defer rows.Close()

	followersIDs := make([]int64, 0)
	for rows.Next() {
		var followerID int64
		err = rows.Scan(&followerID)
		if err != nil {
			log.WithError(err).Error("IDsOfFollowers row scan error.")
			return nil, err
		}

		followersIDs = append(followersIDs, followerID)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("IDsOfFollowers rows iteration error.")
		return nil, err
	}

	return followersIDs, nil
}

// TODO: this is almost a copy-paste of /\. Refactor.
func (db *userDB) IDsOfFollowees(userID int64) ([]int64, error) {
	rows, err := db.Query(`
		SELECT followee_id
		FROM follows
		WHERE follower_id = $1`,
		userID)
	if err != nil {
		log.WithError(err).Error("IDsOfFollowees query error")
	}
	defer rows.Close()

	followeesIDs := make([]int64, 0)
	for rows.Next() {
		var followeeID int64
		err = rows.Scan(&followeeID)
		if err != nil {
			log.WithError(err).Error("IDsOfFollowees row scan error.")
			return nil, err
		}

		followeesIDs = append(followeesIDs, followeeID)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("IDsOfFollowees rows iteration error.")
		return nil, err
	}

	return followeesIDs, nil
}

func (db *userDB) FollowerCount(userID int64) (int64, error) {
	var followerCount int64

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM follows
		WHERE followee_id = $1`,
		userID).
		Scan(&followerCount)
	if err != nil {
		log.WithError(err).Error("FollowerCount query error.")
		return 0, err
	}

	return followerCount, nil
}

// TODO: this is a copy paste of /\
func (db *userDB) FolloweeCount(userID int64) (int64, error) {
	var followeeCount int64

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM follows
		WHERE follower_id = $1`,
		userID).
		Scan(&followeeCount)
	if err != nil {
		log.WithError(err).Error("FolloweeCount query error.")
		return 0, err
	}

	return followeeCount, nil
}

func (db *userDB) IsFollowing(followerID, followeeID int64) (bool, error) {
	var isFollowing bool

	err := db.QueryRow(`
		SELECT exists
			(SELECT true
			FROM follows
			WHERE follower_id = $1 AND followee_id = $2)
		`,
		followerID, followeeID).
		Scan(&isFollowing)
	if err != nil {
		log.WithError(err).Error("IsFollowing query error.")
		return false, err
	}

	return isFollowing, nil
}
