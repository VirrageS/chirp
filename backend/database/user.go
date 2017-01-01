package database

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/model"
)

// User Data Access Object - provides operations on User database table
type UserDAO interface {
	GetPublicUsers() ([]*model.PublicUser, error)
	GetPublicUsersFromListOfIDs(usersToFindIDs []int64) ([]*model.PublicUser, error)
	GetPublicUserWithID(userID int64) (*model.PublicUser, error)
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

func (db *userDB) GetPublicUsersFromListOfIDs(usersToFindIDs []int64) ([]*model.PublicUser, error) {
	// TODO: be careful - this ANY query is said to be super inefficient
	rows, err := db.Query(`SELECT id, username, name, avatar_url FROM users WHERE users.id = ANY($1)`, pq.Array(usersToFindIDs))
	if err != nil {
		log.WithField("usersToFindIDs", usersToFindIDs).WithError(err).Error("GetPublicUsersFromListOfIDs query error.")
		return nil, err
	}
	defer rows.Close()

	users, err := readMultipleUsers(rows)
	if err != nil {
		log.WithError(err).Error("GetPublicUsersFromListOfIDs rows scan/iteration error.")
	}

	return users, nil
}

func (db *userDB) GetPublicUserWithID(userID int64) (*model.PublicUser, error) {
	row := db.QueryRow(`
		SELECT id, username, name, avatar_url
		FROM users
		WHERE id = $1`,
		userID)

	user, err := readPublicUser(row)
	if err != nil && err != sql.ErrNoRows {
		log.WithField("userID", userID).WithError(err).Error("GetPublicUserWithID query error.")
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
	if err != nil && err != sql.ErrNoRows {
		log.WithField("userEmail", userEmail).WithError(err).Error("GetUserWithEmail query error.")
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
