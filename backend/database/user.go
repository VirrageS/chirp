package database

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
)

// User Data Access Object - provides operations on User database table
type UserDAO interface {
	GetPublicUsers() ([]*model.PublicUser, error)
	GetPublicUsersByIDs(usersIDs []int64) ([]*model.PublicUser, error)
	GetPublicUserByID(userID int64) (*model.PublicUser, error)
	GetUserByEmail(userEmail string) (*model.User, error)
	InsertUser(user *model.NewUserForm) (*model.PublicUser, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error
}

type userDB struct {
	*sql.DB
}

func NewUserDAO(dbConnection *sql.DB) UserDAO {
	return &userDB{dbConnection}
}

func (db *userDB) GetPublicUsers() ([]*model.PublicUser, error) {
	rows, err := db.Query(`SELECT id, username, name, avatarl_url FROM users ORDER BY id DESC`)
	if err != nil {
		log.WithError(err).Error("GetPublicUsers query error.")
		return nil, err
	}
	defer rows.Close()

	users, err := readMultipleUsers(rows)
	if err != nil {
		log.WithError(err).Error("GetPublicUsers rows scan/iteration error.")
	}

	return users, nil
}

func (db *userDB) GetPublicUsersByIDs(usersIDs []int64) ([]*model.PublicUser, error) {
	// TODO: be careful - this ANY query is said to be super inefficient
	rows, err := db.Query(`SELECT id, username, name, avatar_url FROM users WHERE users.id = ANY($1)`, pq.Array(usersIDs))
	if err != nil {
		log.WithField("usersIDs", usersIDs).WithError(err).Error("GetPublicUsersByIDs query error.")
		return nil, err
	}
	defer rows.Close()

	users, err := readMultipleUsers(rows)
	if err != nil {
		log.WithError(err).Error("GetPublicUsersByIDs rows scan/iteration error.")
	}

	return users, nil
}

func (db *userDB) GetPublicUserByID(userID int64) (*model.PublicUser, error) {
	row := db.QueryRow(`
		SELECT id, username, name, avatar_url
		FROM users
		WHERE id = $1`,
		userID)

	user, err := readPublicUser(row)
	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}
	if err != nil {
		log.WithField("userID", userID).WithError(err).Error("GetPublicUserByID query error.")
		return nil, err
	}

	return user, err
}

func (db *userDB) GetUserByEmail(userEmail string) (*model.User, error) {
	row := db.QueryRow(`
		SELECT id, username, password, email, name,
			twitter_token, facebook_token, google_token,
			created_at, last_login, active, avatar_url
		FROM users
		WHERE email = $1`,
		userEmail)

	user, err := readUser(row)
	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}
	if err != nil {
		log.WithField("userEmail", userEmail).WithError(err).Error("GetUserByEmail query error.")
		return nil, err
	}

	return user, err
}

func (db *userDB) InsertUser(newUser *model.NewUserForm) (*model.PublicUser, error) {
	// for Postgres we need to use query with RETURNING id to get the ID of the inserted user
	row := db.QueryRow(`
		INSERT INTO users (username, email, password, name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, name, avatar_url`,
		newUser.Username, newUser.Email, newUser.Password, newUser.Name)

	insertedUser, err := readPublicUser(row)
	if err != nil {
		log.WithField("user", *newUser).WithError(err).Error("InsertUser query error.")
		return nil, err
	}

	return insertedUser, nil
}

func (db *userDB) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	_, err := db.Exec(`
		UPDATE users
		SET last_login = $1
		WHERE id = $2`,
		lastLoginTime, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"userID":        userID,
			"lastLoginTime": lastLoginTime,
		}).WithError(err).Error("UpdateUserLastLoginTime query error.")

		return err
	}

	return nil
}
