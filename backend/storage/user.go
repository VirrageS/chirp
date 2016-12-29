package storage

// TODO: maybe prepare statements? http://go-database-sql.org/prepared.html

import (
	"database/sql"
	"time"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
	"github.com/lib/pq"
)

// Struct that implements UserDataAcessor using given DAO and cache
type UserStorage struct {
	DAO   database.UserDAO
	cache cache.CacheProvider
}

// Constructs UserDB that uses a DAO and CacheProvider
func NewUserStorage(DAO database.UserDAO, cache cache.CacheProvider) *UserStorage {
	return &UserStorage{
		DAO,
		cache,
	}
}

func (s *UserStorage) GetUsers(requestingUserID int64) ([]*model.PublicUser, error) {
	users := make([]*model.PublicUser, 0)
	if exists, _ := s.cache.GetWithFields(cache.Fields{"users", requestingUserID}, &users); exists {
		return users, nil
	}

	users, err := s.DAO.GetPublicUsers(requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	s.cache.SetWithFields(cache.Fields{"users", requestingUserID}, users)
	return users, nil
}

func (s *UserStorage) GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error) {
	var user *model.PublicUser
	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", "id", userID, requestingUserID}, &user); exists {
		return user, nil
	}

	user, err := s.DAO.GetPublicUserUsingQuery(`
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

	s.cache.SetWithFields(cache.Fields{"user", "id", userID, requestingUserID}, user)
	return user, nil
}

func (s *UserStorage) GetUserByEmail(email string) (*model.User, error) {
	var user *model.User
	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", "email", email}, &user); exists {
		return user, nil
	}

	user, err := s.DAO.GetUserUsingQuery("SELECT * FROM users WHERE email=$1", email)
	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}
	if err != nil {
		return nil, errors.UnexpectedError
	}

	s.cache.SetWithFields(cache.Fields{"user", "email", email}, user)
	return user, nil
}

func (s *UserStorage) InsertUser(newUserForm *model.NewUserForm) (*model.PublicUser, error) {
	userID, err := s.DAO.InsertUserToDatabase(newUserForm)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == database.UniqueConstraintViolationCode {
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

func (s *UserStorage) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	err := s.DAO.UpdateUserLastLoginTime(userID, lastLoginTime)
	if err != nil {
		return errors.UnexpectedError
	}

	// No point updating cache, because that is not a very important data and updating it would need to invalidate
	// whole cache.

	return nil
}

func (s *UserStorage) FollowUser(followeeID, followerID int64) error {
	err := s.DAO.FollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	// TODO: Maybe a smarter way: don't delete, but just update cache with followerCount++ and following=true
	// Just delete from cache for the requesting user, it will be fetched back in next GET query
	s.cache.DeleteWithFields(cache.Fields{"user", followeeID, followerID})

	return nil
}

func (s *UserStorage) UnfollowUser(followeeID, followerID int64) error {
	err := s.DAO.UnfollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	// TODO: Maybe a smarter way: don't delete, but just update cache with followerCount-- and following=false
	// Just delete from cache for the requesting user, it will be fetched back in next GET query
	s.cache.DeleteWithFields(cache.Fields{"user", followeeID, followerID})

	return nil
}

func (s *UserStorage) Followers(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followersIDs, err := s.DAO.Followers(userID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	followers := make([]*model.PublicUser, 0)
	for i, id := range followersIDs {
		var user model.PublicUser

		if exists, _ := s.cache.GetWithFields(cache.Fields{"user", "id", id, requestingUserID}, &user); exists {
			followers = append(followers, &user)

			// remove ID from followingIDs
			followersIDs[i] = followersIDs[len(followersIDs)-1]
			followersIDs = followersIDs[:len(followersIDs)-1]
		}
	}

	if len(followersIDs) > 0 {
		dbFollowers, err := s.DAO.GetPublicUsersFromListOfIDs(requestingUserID, followersIDs)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		followers = append(followers, dbFollowers...)
	}

	return followers, nil
}

func (s *UserStorage) Followees(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followeesIDs, err := s.DAO.Followees(userID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	followees := make([]*model.PublicUser, 0)
	for i, id := range followeesIDs {
		var user model.PublicUser

		if exists, _ := s.cache.GetWithFields(cache.Fields{"user", "id", id, requestingUserID}, &user); exists {
			followees = append(followees, &user)

			// remove ID from followersIDs
			followeesIDs[i] = followeesIDs[len(followeesIDs)-1]
			followeesIDs = followeesIDs[:len(followeesIDs)-1]
		}
	}

	if len(followeesIDs) > 0 {
		dbFollowees, err := s.DAO.GetPublicUsersFromListOfIDs(requestingUserID, followeesIDs)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		followees = append(followees, dbFollowees...)
	}

	return followees, nil
}
