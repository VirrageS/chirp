package database

// TODO: maybe prepare statements? http://go-database-sql.org/prepared.html

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/database/model"
	"github.com/lib/pq"
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

func (db *UserDB) GetUsers() ([]*model.User, error) {
	users, err := db.getUsers()
	if err != nil {
		return nil, DatabaseError
	}

	return users, nil
}

func (db *UserDB) GetUserByID(userID int64) (*model.User, error) {
	user, err := db.getUserUsingQuery("SELECT * from users WHERE id=$1", userID)
	if err == sql.ErrNoRows {
		return nil, NoRowsError
	}
	if err != nil {
		return nil, DatabaseError
	}

	return user, nil
}

func (db *UserDB) GetUserByEmail(email *string) (*model.User, error) {
	user, err := db.getUserUsingQuery("SELECT * from users WHERE email=$1", email)
	if err == sql.ErrNoRows {
		return nil, NoRowsError
	}
	if err != nil {
		return nil, DatabaseError
	}

	return user, nil
}

func (db *UserDB) InsertUser(user *model.User) (*model.User, error) {
	userID, err := db.insertUserToDatabase(user)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == UniqueConstraintViolationCode {
			return nil, UserAlreadyExistsError
		}
		return nil, DatabaseError
	}

	user.ID = userID

	return user, nil
}

func (db *UserDB) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	err := db.updateUserLastLoginTime(userID, lastLoginTime)
	if err != nil {
		return DatabaseError
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
		log.WithField("query", query).WithError(err).Error("GetUserUsingQuery database error.")
	}

	return &user, err
}

func (db *UserDB) insertUserToDatabase(user *model.User) (int64, error) {
	query, err := db.Prepare("INSERT INTO users (username, email, password, created_at, last_login, name)" +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		log.WithError(err).Error("insertUserToDatabase query prepare error.")
		return 0, err
	}
	defer query.Close()

	var newID int64
	// for Postgres we need to use query with RETURNING id to get the ID of the inserted user
	err = query.QueryRow(user.Username, user.Email, user.Password, user.CreatedAt, user.LastLogin, user.Name).Scan(&newID)

	if err != nil {
		log.WithError(err).Error("insertUserToDatabase query execute error.")
		return 0, err
	}

	return newID, nil
}

// TODO: add filtering parameters
func (db *UserDB) getUsers() ([]*model.User, error) {
	rows, err := db.Query("SELECT * FROM users;")
	if err != nil {
		log.WithError(err).Error("GetUsers query error.")
		return nil, err
	}

	var users []*model.User

	defer rows.Close()
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.TwitterToken, &user.FacebookToken, &user.GoogleToken, &user.Username,
			&user.Email, &user.Password, &user.CreatedAt, &user.LastLogin, &user.Active,
			&user.Name, &user.AvatarUrl)
		if err != nil {
			log.WithError(err).Error("getUsers row scan error.")
			return nil, err
		}

		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("getUsers rows iteration error.")
		return nil, err
	}

	return users, nil
}

func (db *UserDB) updateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	query, err := db.Prepare("UPDATE users SET last_login=$1 WHERE id=$2;")
	if err != nil {
		log.WithError(err).Error("updateUserLastLoginTime query prepare error.")
		return err
	}
	defer query.Close()

	_, err = query.Exec(lastLoginTime, userID)
	if err != nil {
		log.WithError(err).Error("updateUserLastLoginTime query execute error.")
		return err
	}

	return nil
}
