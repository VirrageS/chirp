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
	users := make([]*model.PublicUser, 0)
	if exists, _ := db.cache.GetWithFields(cache.Fields{"users", "requser", requestingUserID}, &users); exists {
		return users, nil
	}

	users, err := db.getPublicUsers(requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	db.cache.SetWithFields(cache.Fields{"users", "requser", requestingUserID}, users)
	return users, nil
}

func (db *UserDB) GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error) {
	var user *model.PublicUser
	if exists, _ := db.cache.GetWithFields(cache.Fields{"user", "id", userID, "requser", requestingUserID}, &user); exists {
		return user, nil
	}

	user, err := db.getPublicUserUsingQuery(`
		SELECT id, username, name, avatar_url,
			COUNT(follows.follower_id) AS follow_count,
			SUM(CASE WHEN follows.follower_id=$2 THEN 1 ELSE 0 END) > 0 AS following
		FROM users
			LEFT JOIN follows
			ON users.id = follows.followee_id
		WHERE users.id = $1
		GROUP BY users.id;`,
		userID, requestingUserID)

	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}

	if err != nil {
		return nil, errors.UnexpectedError
	}

	db.cache.SetWithFields(cache.Fields{"user", "id", userID, "requser", requestingUserID}, user)
	return user, nil
}

func (db *UserDB) GetUserByEmail(email string) (*model.User, error) {
	var user *model.User
	if exists, _ := db.cache.GetWithFields(cache.Fields{"user", "email", email}, &user); exists {
		return user, nil
	}

	user, err := db.getUserUsingQuery("SELECT * FROM users WHERE email=$1", email)
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

	// We don't flush cache on purpose. The data in cache can be not precise for some time.
	// We also don't add the user to cache because this would not make sense, since we compute additional data
	// for each user that depends on requesting user.

	return newPublicUser, nil
}

func (db *UserDB) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	err := db.updateUserLastLoginTime(userID, lastLoginTime)
	if err != nil {
		return errors.UnexpectedError
	}

	// No point updating cache, because that is not a very important data and updating it would need to invalidate
	// whole cache.

	return nil
}

func (db *UserDB) FollowUser(followeeID, followerID int64) error {
	err := db.followUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	// TODO: Maybe a smarter way: don't delete, but just update cache with followerCount++ and following=true
	// Just delete from cache for the requesting user, it will be fetched back in next GET query
	db.cache.DeleteWithFields(cache.Fields{"user", followeeID, "requser", followerID})

	return nil
}

func (db *UserDB) UnfollowUser(followeeID, followerID int64) error {
	err := db.unfollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	// TODO: Maybe a smarter way: don't delete, but just update cache with followerCount-- and following=false
	// Just delete from cache for the requesting user, it will be fetched back in next GET query
	db.cache.DeleteWithFields(cache.Fields{"user", followeeID, "requser", followerID})

	return nil
}

func (db *UserDB) Followers(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followersIDs, err := db.followers(userID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	followers := make([]*model.PublicUser, 0)
	for i, id := range followersIDs {
		var user model.PublicUser

		if exists, _ := db.cache.GetWithFields(cache.Fields{"user", "id", id, "requser", requestingUserID}, &user); exists {
			followers = append(followers, &user)

			// remove ID from followingIDs
			followersIDs[i] = followersIDs[len(followersIDs)-1]
			followersIDs = followersIDs[:len(followersIDs)-1]
		}
	}

	if len(followersIDs) > 0 {
		dbFollowers, err := db.getPublicUsersFromListOfIDs(requestingUserID, followersIDs)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		followers = append(followers, dbFollowers...)
	}

	return followers, nil
}

func (db *UserDB) Followees(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followeesIDs, err := db.followees(userID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	followees := make([]*model.PublicUser, 0)
	for i, id := range followeesIDs {
		var user model.PublicUser

		if exists, _ := db.cache.GetWithFields(cache.Fields{"user", "id", id, "requser", requestingUserID}, &user); exists {
			followees = append(followees, &user)

			// remove ID from followersIDs
			followeesIDs[i] = followeesIDs[len(followeesIDs)-1]
			followeesIDs = followeesIDs[:len(followeesIDs)-1]
		}
	}

	if len(followeesIDs) > 0 {
		dbFollowees, err := db.getPublicUsersFromListOfIDs(requestingUserID, followeesIDs)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		followees = append(followees, dbFollowees...)
	}

	return followees, nil
}

func (db *UserDB) getPublicUsers(requestingUserID int64) ([]*model.PublicUser, error) {
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
func (db *UserDB) getPublicUsersFromListOfIDs(requestingUserID int64, usersToFindIDs []int64) ([]*model.PublicUser, error) {
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
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarUrl, &user.FollowerCount, &user.Following)

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

// TODO: this is a temporary workaround
func (db *UserDB) followers(userID int64) ([]int64, error) {
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
func (db *UserDB) followees(userID int64) ([]int64, error) {
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
