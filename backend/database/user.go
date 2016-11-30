package database

// TODO: maybe prepare statements? http://go-database-sql.org/prepared.html

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/model"
	"github.com/lib/pq"
)

// Struct that implements UserDataAccessor using sql (postgres) database
type UserDB struct {
	*sql.DB
}

// Constructs UserDB that uses a given sql.DB connection
func NewUserDB(databaseConnection *sql.DB) *UserDB {
	return &UserDB{databaseConnection}
}

func (db *UserDB) GetUsers() ([]*model.PublicUser, error) {
	users, err := db.getPublicUsers()
	if err != nil {
		return nil, DatabaseError
	}

	return users, nil
}

func (db *UserDB) GetUserByID(userID int64) (*model.PublicUser, error) {
	user, err := db.getPublicUserUsingQuery("SELECT id, username, name, avatar_url from users WHERE id=$1", userID)
	if err == sql.ErrNoRows {
		return nil, NoResults
	}
	if err != nil {
		return nil, DatabaseError
	}

	return user, nil
}

func (db *UserDB) GetUserByEmail(email string) (*model.User, error) {
	user, err := db.getUserUsingQuery("SELECT * from users WHERE email=$1", email)
	if err == sql.ErrNoRows {
		return nil, NoResults
	}
	if err != nil {
		return nil, DatabaseError
	}

	return user, nil
}

func (db *UserDB) InsertUser(newUserForm *model.NewUserForm) (*model.PublicUser, error) {
	// TODO: FIX ME PLEASE
	newUser := &model.User{
		ID:            0,
		TwitterToken:  toSqlNullString(""),
		FacebookToken: toSqlNullString(""),
		GoogleToken:   toSqlNullString(""),
		Username:      newUserForm.Username,
		Password:      newUserForm.Password,
		Email:         newUserForm.Email,
		CreatedAt:     time.Now(),
		LastLogin:     time.Now(),
		Active:        true,
		Name:          newUserForm.Name,
		AvatarUrl:     toSqlNullString(""),
		Following:     false, // TODO: NOT ME HERE PLEASE
	}

	userID, err := db.insertUserToDatabase(newUser)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == UniqueConstraintViolationCode {
			return nil, UserAlreadyExistsError
		}
		return nil, DatabaseError
	}

	// TODO: FIX ME PLEASE TOO!
	newPublicUser := &model.PublicUser{
		ID:        userID,
		Username:  newUser.Username,
		Name:      newUser.Name,
		AvatarUrl: "",
		Following: false,
	}

	return newPublicUser, nil
}

func (db *UserDB) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	err := db.updateUserLastLoginTime(userID, lastLoginTime)
	if err != nil {
		return DatabaseError
	}

	return nil
}

func (db *UserDB) getPublicUsers() ([]*model.PublicUser, error) {
	rows, err := db.Query("SELECT id, username, last_login, name, avatar_url FROM users;")
	if err != nil {
		log.WithError(err).Error("getPublicUsers query error.")
		return nil, err
	}

	var users []*model.PublicUser

	defer rows.Close()
	for rows.Next() {
		var user model.PublicUser
		err = rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl)
		if err != nil {
			log.WithError(err).Error("getPublicUsers row scan error.")
			return nil, err
		}

		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("getPublicUsers rows iteration error.")
		return nil, err
	}

	return users, nil
}

func (db *UserDB) getUserUsingQuery(query string, args ...interface{}) (*model.User, error) {
	var user model.User

	row := db.QueryRow(query, args...)
	err := row.Scan(&user.ID, &user.TwitterToken, &user.FacebookToken, &user.GoogleToken, &user.Username,
		&user.Email, &user.Password, &user.CreatedAt, &user.LastLogin, &user.Active,
		&user.Name, &user.AvatarUrl)

	if err != nil && err != sql.ErrNoRows {
		log.WithField("query", query).WithError(err).Error("getUserUsingQuery database error.")
		return nil, err
	}

	return &user, err
}

func (db *UserDB) getPublicUserUsingQuery(query string, args ...interface{}) (*model.PublicUser, error) {
	var user model.PublicUser

	row := db.QueryRow(query, args...)
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl)

	if err != nil && err != sql.ErrNoRows {
		log.WithField("query", query).WithError(err).Error("getPublicUserUsingQuery database error.")
		return nil, err
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

func toSqlNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
