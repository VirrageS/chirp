package database

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/model"
)

type UserDAO interface {
	GetPublicUsers(requestingUserID int64) ([]*model.PublicUser, error)
	GetPublicUsersFromListOfIDs(requestingUserID int64, usersToFindIDs []int64) ([]*model.PublicUser, error)
	GetUserUsingQuery(query string, args ...interface{}) (*model.User, error)
	GetPublicUserUsingQuery(query string, args ...interface{}) (*model.PublicUser, error)
	InsertUserToDatabase(user *model.NewUserForm) (int64, error)
	UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error
	FollowUser(followeeID, followerID int64) error
	UnfollowUser(followeeID, followerID int64) error
	Followers(userID int64) ([]int64, error)
	Followees(userID int64) ([]int64, error)
}

type userDB struct {
	*sql.DB
}

func NewUserDAO(dbConnection *sql.DB) UserDAO {
	return &userDB{dbConnection}
}

func (db *userDB) GetPublicUsers(requestingUserID int64) ([]*model.PublicUser, error) {
	rows, err := db.Query(`
		SELECT id, username, name, avatar_url,
			COUNT(follows.follower_id) as follow_count,
			SUM(CASE WHEN follows.follower_id=$1 THEN 1 ELSE 0 END) > 0 AS following
		FROM users
			LEFT JOIN follows
			ON users.id = follows.followee_id
		GROUP BY users.id;`,
		requestingUserID)
	if err != nil {
		log.WithError(err).Error("getPublicUsers query error.")
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.PublicUser, 0)
	for rows.Next() {
		var user model.PublicUser
		err = rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl, &user.FollowerCount, &user.Following)
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

// TODO this is almost a copy-paste of /\. REFACTOR
func (db *userDB) GetPublicUsersFromListOfIDs(requestingUserID int64, usersToFindIDs []int64) ([]*model.PublicUser, error) {
	// TODO: be careful - this ANY query is said to be super inefficient
	query := `SELECT id, username, name, avatar_url,
			COUNT(follows.follower_id) as follow_count,
			SUM(CASE WHEN follows.follower_id=$1 THEN 1 ELSE 0 END) > 0 AS following
		FROM users
			LEFT JOIN follows
			ON users.id = follows.followee_id
		WHERE users.id = ANY($2)
		GROUP BY users.id;`

	rows, err := db.Query(query, requestingUserID, pq.Array(usersToFindIDs))
	if err != nil {
		log.WithError(err).WithField("query", query).Error("getPublicUsersFromListOfIDs query error.")
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.PublicUser, 0)
	for rows.Next() {
		var user model.PublicUser
		err = rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl, &user.FollowerCount, &user.Following)
		if err != nil {
			log.WithError(err).Error("getPublicUsersFromListOfIDs row scan error.")
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

func (db *userDB) GetUserUsingQuery(query string, args ...interface{}) (*model.User, error) {
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

func (db *userDB) GetPublicUserUsingQuery(query string, args ...interface{}) (*model.PublicUser, error) {
	var user model.PublicUser

	row := db.QueryRow(query, args...)
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl, &user.FollowerCount, &user.Following)

	if err != nil && err != sql.ErrNoRows {
		log.WithField("query", query).WithError(err).Error("getPublicUserUsingQuery database error.")
		return nil, err
	}

	return &user, err
}

func (db *userDB) InsertUserToDatabase(user *model.NewUserForm) (int64, error) {
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

func (db *userDB) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
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

func (db *userDB) FollowUser(followeeID, followerID int64) error {
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

func (db *userDB) UnfollowUser(followeeID, followerID int64) error {
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

// TODO: this is a temporary workaround
func (db *userDB) Followers(userID int64) ([]int64, error) {
	rows, err := db.Query(`
		SELECT follower_id
		FROM users
			INNER JOIN follows
			ON users.id = follows.followee_id
		WHERE users.id = $1;`,
		userID)

	if err != nil {
		log.WithError(err).Error("followers query error")
	}
	defer rows.Close()

	followersIDs := make([]int64, 0)
	for rows.Next() {
		var followerID int64
		err = rows.Scan(&followerID)
		if err != nil {
			log.WithError(err).Error("followers row scan error.")
			return nil, err
		}

		followersIDs = append(followersIDs, followerID)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("followers rows iteration error.")
		return nil, err
	}

	return followersIDs, nil
}

// TODO: this is almost a copy-paste of /\. Refactor.
func (db *userDB) Followees(userID int64) ([]int64, error) {
	rows, err := db.Query(`
		SELECT followee_id
		FROM users
			INNER JOIN follows
			ON users.id = follows.follower_id
		WHERE users.id = $1;`,
		userID)

	if err != nil {
		log.WithError(err).Error("followees query error")
	}
	defer rows.Close()

	followeesIDs := make([]int64, 0)
	for rows.Next() {
		var followeeID int64
		err = rows.Scan(&followeeID)
		if err != nil {
			log.WithError(err).Error("followees row scan error.")
			return nil, err
		}

		followeesIDs = append(followeesIDs, followeeID)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("followees rows iteration error.")
		return nil, err
	}

	return followeesIDs, nil
}
