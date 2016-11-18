package database

// TODO: maybe prepare statements? http://go-database-sql.org/prepared.html

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/database/model"
)

type UserDataAccessor interface {
	GetUsers() ([]*model.User, error)
	GetUserByID(userID int64) (*model.User, error)
	GetUserByEmail(email *string) (*model.User, error)
	InsertUser(user *model.User) (*model.User, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error
}

type UserDB struct {
	*sql.DB
}

func NewUserDB(databaseConnection *sql.DB) *UserDB {
	return &UserDB{databaseConnection}
}

// Interface implementations

func (db *UserDB) GetUsers() ([]*model.User, error) {
	return db.getUsers()
}

func (db *UserDB) GetUserByID(userID int64) (*model.User, error) {
	user, err := db.getUserUsingQuery("SELECT * from users WHERE id=$1", userID)
	if err == sql.ErrNoRows {
		return nil, errors.New("") // no users found
	}
	if err != nil {
		return nil, errors.New("") // db error
	}

	return user, nil
}

func (db *UserDB) GetUserByEmail(email *string) (*model.User, error) {
	user, err := db.getUserUsingQuery("SELECT * from users WHERE email=$1", email)
	if err != nil {
		return nil, errors.New("") // db error
	}

	return user, nil
}

func (db *UserDB) InsertUser(user *model.User) (*model.User, error) {
	exists, err := db.checkIfUserAlreadyExists(user)
	if err != nil {
		return nil, errors.New("") // db error
	}

	if exists {
		// TODO: return a message that informs the user which one of username/email is already taken
		return nil, errors.New("") // user already exists
	}

	userID, err := db.insertUserToDatabase(user)
	if err != nil {
		return nil, errors.New("") // db error
	}

	user.ID = userID

	return user, nil
}

func (db *UserDB) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	err := db.updateUserLastLoginTime(userID, lastLoginTime)
	if err != nil {
		return errors.New("") // unexpected error, internal server error
	}

	return nil
}

func (db *UserDB) getUserUsingQuery(query string, args ...interface{}) (*model.User, error) {
	var user model.User

	row := db.QueryRow(query, args...)
	err := row.Scan(&user.ID, &user.TwitterToken, &user.FacebookToken, &user.GoogleToken, &user.Username,
		&user.Email, &user.Password, &user.CreatedAt, &user.LastLogin, &user.Active,
		&user.Name, &user.AvatarUrl)

	if err != nil && err != sql.ErrNoRows {
		log.WithFields(log.Fields{
			"error": err,
			"query": query,
			"args":  args,
		}).Error("GetUserUsingQuery database error.")
	}

	return &user, err
}

func (db *UserDB) insertUserToDatabase(user *model.User) (int64, error) {
	query, err := db.Prepare("INSERT INTO users (username, email, password, created_at, last_login, name)" +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		log.WithField("query", query).WithError(err).Error("insertUserToDatabase query prepare error.")
		return 0, errors.New("")
	}
	defer query.Close()

	var newID int64
	// for Postgres we need to use query with RETURNING id to get the ID of the inserted user
	err = query.QueryRow(user.Username, user.Email, user.Password, user.CreatedAt, user.LastLogin, user.Name).Scan(&newID)
	if err != nil {
		log.WithError(err).Error("insertUserToDatabase query execute error.")
		return 0, errors.New("")
	}

	return newID, nil
}

// TODO: find a better name and design for this function
func (db *UserDB) checkIfUserAlreadyExists(userToCheck *model.User) (bool, error) {
	_, err := db.getUserUsingQuery(
		// can be done with an 'exists' query,
		// but we will need to return info about which field is already taken back to the user
		// TODO: return a message indicating which field is already taken
		"SELECT * from users WHERE email=$1 OR username=$2",
		userToCheck.Email,
		userToCheck.Username)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("userAlreadyExists database error.")
		return false, errors.New("")
	}

	return true, nil
}

// TODO: add filtering parameters
func (db *UserDB) getUsers() ([]*model.User, error) {
	rows, err := db.Query("SELECT * FROM users;")
	if err != nil {
		log.WithError(err).Error("GetUsers query error.")
		return nil, errors.New("")
	}

	var users []*model.User

	defer rows.Close()
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.TwitterToken, &user.FacebookToken, &user.GoogleToken, &user.Username,
			&user.Email, &user.Password, &user.CreatedAt, &user.LastLogin, &user.Active,
			&user.Name, &user.AvatarUrl)
		// TODO: move error outside of the loop
		if err != nil {
			log.WithError(err).Error("GetUsers row scan error.")
			return nil, errors.New("")
		}

		users = append(users, &user)
	}

	return users, nil
}

func (db *UserDB) updateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	query, err := db.Prepare("UPDATE users SET last_login=$1 WHERE id=$2;")
	if err != nil {
		log.WithField("query", query).WithError(err).Error("updateUserLastLoginTime query prepare error.")
		return errors.New("")
	}
	defer query.Close()

	_, err = query.Exec(lastLoginTime, userID)
	if err != nil {
		log.WithError(err).Error("updateUserLastLoginTime query execute error.")
		return errors.New("")
	}

	return nil
}
