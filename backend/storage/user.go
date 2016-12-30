package storage

// TODO: maybe prepare statements? http://go-database-sql.org/prepared.html

import (
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
)

// Struct that implements UserDataAcessor using given DAO and cache
type UserStorage struct {
	userDAO    database.UserDAO
	followsDAO database.FollowsDAO
	cache      cache.CacheProvider
}

// Constructs UserDB that uses a DAO and CacheProvider
func NewUserStorage(userDAO database.UserDAO, followsDAO database.FollowsDAO, cache cache.CacheProvider) *UserStorage {
	return &UserStorage{
		userDAO:    userDAO,
		followsDAO: followsDAO,
		cache:      cache,
	}
}

func (s *UserStorage) GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error) {
	var user *model.PublicUser

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", userID}, &user); !exists {
		var err error

		user, err = s.userDAO.GetPublicUserWithID(userID)
		if err == sql.ErrNoRows {
			return nil, errors.NoResultsError
		}
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"user", userID}, user)
	}

	err := s.collectAllPublicUserData(user, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return user, nil
}

func (s *UserStorage) GetAuthDataOfUserWithEmail(email string) (*model.User, error) {
	var user *model.User

	// Here we don't need any additional data other than what we fetched from database.
	// Dont use cache here, since this function will be used only for authentication users and we want
	// 100% real data for this.

	user, err := s.userDAO.GetUserWithEmail(email)
	if err == sql.ErrNoRows {
		return nil, errors.NoResultsError
	}
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return user, nil
}

func (s *UserStorage) InsertUser(newUserForm *model.NewUserForm) (*model.PublicUser, error) {
	insertedUser, err := s.userDAO.InsertUser(newUserForm)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == database.UniqueConstraintViolationCode {
			return nil, errors.UserAlreadyExistsError
		}
		return nil, errors.UnexpectedError
	}

	s.cache.SetWithFields(cache.Fields{"user", insertedUser.ID}, insertedUser.ID)
	s.cache.SetIntWithFields(cache.Fields{"user", insertedUser.ID, "followerCount"}, 0)

	return insertedUser, nil
}

func (s *UserStorage) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	err := s.userDAO.UpdateUserLastLoginTime(userID, lastLoginTime)
	if err != nil {
		return errors.UnexpectedError
	}

	// Maybe we could Get and old one, update and set new one? Might not make a very big difference, since
	// login doesnt happen very often. (and also this function is called only in one place
	// (login function in service) which calls 'getUser' that will fetch the data back to cache anyway.
	s.cache.DeleteWithFields(cache.Fields{"user", userID})

	// data about this user stays in some cached lists of user, but it doesn't make a big difference,
	// because this is not a very important information

	return nil
}

func (s *UserStorage) FollowUser(followeeID, followerID int64) error {
	err := s.followsDAO.FollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	s.cache.SetWithFields(cache.Fields{"user", followeeID, "isFollowedBy", followerID}, true)
	s.cache.IncrementWithFields(cache.Fields{"user", followeeID, "followerCount"})

	return nil
}

func (s *UserStorage) UnfollowUser(followeeID, followerID int64) error {
	err := s.followsDAO.UnfollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	s.cache.SetWithFields(cache.Fields{"user", followeeID, "isFollowedBy", followerID}, false)
	s.cache.DecrementWithFields(cache.Fields{"user", followeeID, "followerCount"})

	return nil
}

// TODO: this all could be done nicely in paralell
func (s *UserStorage) Followers(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	var followersIDs []int64
	followers := make([]*model.PublicUser, 0)

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", userID, "followers"}, &followers); !exists {
		var err error

		followersIDs, err = s.followsDAO.IDsOfFollowers(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		followers, err = s.collectUsersFromListOfIDs(followersIDs)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"user", userID, "followers"}, &followers)
	}

	for _, user := range followers {
		err := s.collectAllPublicUserData(user, requestingUserID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
	}

	return followers, nil
}

// TODO: those to functions /\ and \/ are almost identical. It might be able to refactor them

// TODO: this all could be done nicely in paralell
func (s *UserStorage) Followees(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	var followeesIDs []int64
	followees := make([]*model.PublicUser, 0)

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", userID, "followees"}, &followees); !exists {
		var err error

		followeesIDs, err = s.followsDAO.IDsOfFollowees(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		followees, err = s.collectUsersFromListOfIDs(followeesIDs)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"user", userID, "followees"}, &followees)
	}

	for _, user := range followees {
		err := s.collectAllPublicUserData(user, requestingUserID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
	}

	return followees, nil
}

// This is useful only for testing/debugging, skip cache here
func (s *UserStorage) GetUsers(requestingUserID int64) ([]*model.PublicUser, error) {
	users := make([]*model.PublicUser, 0)

	users, err := s.userDAO.GetPublicUsers()
	if err != nil {
		return nil, errors.UnexpectedError
	}

	for _, user := range users {
		followerCount, err := s.followsDAO.FollowerCount(user.ID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		following, err := s.followsDAO.IsFollowing(requestingUserID, user.ID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		user.FollowerCount = followerCount
		user.Following = following
	}

	return users, nil
}

// Be careful - this is function does SIDE EFFECTS only
func (s *UserStorage) collectAllPublicUserData(user *model.PublicUser, requestingUserID int64) error {
	var followerCount int64
	var following bool

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", user.ID, "followerCount"}, &followerCount); !exists {
		var err error

		followerCount, err = s.followsDAO.FollowerCount(user.ID)
		if err != nil {
			return errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"user", user.ID, "followerCount"}, followerCount)
	}

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", user.ID, "isFollowedBy", requestingUserID}, &following); !exists {
		var err error

		following, err = s.followsDAO.IsFollowing(requestingUserID, user.ID)
		if err != nil {
			return errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"user", user.ID, "isFollowedBy", requestingUserID}, following)
	}

	user.FollowerCount = followerCount
	user.Following = following

	return nil
}

func (s *UserStorage) collectUsersFromListOfIDs(usersIDs []int64) ([]*model.PublicUser, error) {
	users := make([]*model.PublicUser, 0)

	// get users from cache
	for i, id := range usersIDs {
		var user model.PublicUser

		if exists, _ := s.cache.GetWithFields(cache.Fields{"user", id}, &user); exists {
			users = append(users, &user)

			// remove ID from usersIDs
			usersIDs[i] = usersIDs[len(usersIDs)-1]
			usersIDs = usersIDs[:len(usersIDs)-1]
		}
	}

	// get users that are not in cache from database
	if len(usersIDs) > 0 {
		dbFollowers, err := s.userDAO.GetPublicUsersFromListOfIDs(usersIDs)
		if err != nil {
			return nil, err
		}
		users = append(users, dbFollowers...)
	}

	return users, nil
}
