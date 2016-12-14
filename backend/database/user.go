package database

// TODO: maybe prepare statements? http://go-database-sql.org/prepared.html

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
	"github.com/lib/pq"
)

// Struct that implements UserDataAccessor using sql (postgres) database
type UserDB struct {
	*sql.DB
	cache cache.CacheProvider
}

// Constructs UserDB that uses a given sql.DB connection and CacheProvider
func NewUserDB(databaseConnection *sql.DB, cache cache.CacheProvider) *UserDB {
	return &UserDB{
		databaseConnection,
		cache,
	}
}

func (db *UserDB) GetUsers(requestingUserID int64) ([]*model.PublicUser, error) {
	if users, exists := db.cache.Get("users"); exists {
		return users.([]*model.PublicUser), nil
	}

	users, err := db.getPublicUsers(requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	db.cache.Set("users", users)
	return users, nil
}

func (db *UserDB) GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error) {
	if user, exists := db.cache.GetWithFields(cache.Fields{"user", "id", userID}); exists {
		return user.(*model.PublicUser), nil
	}

	user, err := db.getPublicUserUsingQuery(`
		SELECT id, username, name, avatar_url, SUM(case when follows.follower_id=$2 then 1 else 0 end) > 0 as following
		FROM users
			LEFT JOIN follows on users.id = follows.followee_id
		WHERE users.id = $1
		GROUP BY users.id;`,
		userID, requestingUserID)

	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}

	if err != nil {
		return nil, errors.UnexpectedError
	}

	db.cache.SetWithFields(cache.Fields{"user", "id", userID}, user)
	return user, nil
}

func (db *UserDB) GetUserByEmail(email string) (*model.User, error) {
	if user, exists := db.cache.GetWithFields(cache.Fields{"user", "email", email}); exists {
		return user.(*model.User), nil
	}

	user, err := db.getUserUsingQuery("SELECT * from users WHERE email=$1", email)
	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}
	if err != nil {
		return nil, errors.UnexpectedError
	}

	db.cache.SetWithFields(cache.Fields{"user", "email", email}, user)
	return user, nil
}

func (db *UserDB) InsertUser(newUserForm *model.NewUserForm) (*model.PublicUser, error) {
	userID, err := db.insertUserToDatabase(newUserForm)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == UniqueConstraintViolationCode {
			return nil, errors.UserAlreadyExistsError
		}
		return nil, errors.UnexpectedError
	}

	// TODO: how bad is this? This is ugly, but saves a database query
	newPublicUser := &model.PublicUser{
		ID:        userID,
		Username:  newUserForm.Username,
		Name:      newUserForm.Name,
		AvatarUrl: "",
		Following: false,
	}

	return newPublicUser, nil
}

func (db *UserDB) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	err := db.updateUserLastLoginTime(userID, lastLoginTime)
	if err != nil {
		return errors.UnexpectedError
	}

	return nil
}

func (db *UserDB) FollowUser(followeeID, followerID int64) error {
	err := db.followUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	return nil
}

func (db *UserDB) UnfollowUser(followeeID, followerID int64) error {
	err := db.unfollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	return nil
}

func (db *UserDB) getPublicUsers(requestingUserID int64) ([]*model.PublicUser, error) {
	rows, err := db.Query(`
		SELECT id, username, name, avatar_url, SUM(case when follows.follower_id=$1 then 1 else 0 end) > 0 as following
		FROM users
			LEFT JOIN follows on users.id = follows.followee_id
		GROUP BY users.id;`,
		requestingUserID)
	if err != nil {
		log.WithError(err).Error("getPublicUsers query error.")
		return nil, err
	}

	var users []*model.PublicUser

	defer rows.Close()
	for rows.Next() {
		var user model.PublicUser
		err = rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl, &user.Following)
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
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl, &user.Following)

	if err != nil && err != sql.ErrNoRows {
		log.WithField("query", query).WithError(err).Error("getPublicUserUsingQuery database error.")
		return nil, err
	}

	return &user, err
}

func (db *UserDB) insertUserToDatabase(user *model.NewUserForm) (int64, error) {
	query, err := db.Prepare("INSERT INTO users (username, email, password, name)" +
		"VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		log.WithError(err).Error("insertUserToDatabase query prepare error.")
		return 0, err
	}
	defer query.Close()

	var newID int64
	// for Postgres we need to use query with RETURNING id to get the ID of the inserted user
	err = query.QueryRow(user.Username, user.Email, user.Password, user.Name).Scan(&newID)

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

func (db *UserDB) followUser(followeeID, followerID int64) error {
	query, err := db.Prepare(`
		INSERT INTO follows (followee_id, follower_id)
		VALUES ($1, $2)
		ON CONFLICT (followee_id, follower_id) DO NOTHING;
		`)

	if err != nil {
		log.WithError(err).Error("followUser query prepare error")
		return err
	}
	defer query.Close()

	_, err = query.Exec(followeeID, followerID)
	if err != nil {
		log.WithFields(log.Fields{
			"followeeID": followeeID,
			"followerID": followerID,
		}).WithError(err).Error("followUser query execute error.")
		return err
	}

	return nil
}

func (db *UserDB) unfollowUser(followeeID, followerID int64) error {
	query, err := db.Prepare(`
		DELETE FROM follows
		WHERE followee_id=$1 AND follower_id=$2;
		`)

	if err != nil {
		log.WithError(err).Error("unfollowUser query prepare error")
		return err
	}
	defer query.Close()

	_, err = query.Exec(followeeID, followerID)
	if err != nil {
		log.WithFields(log.Fields{
			"followeeID": followeeID,
			"followerID": followerID,
		}).WithError(err).Error("unfollowUser query execute error.")
		return err
	}

	return nil
}
