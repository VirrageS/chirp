package storage

// TODO: maybe prepare statements? http://go-database-sql.org/prepared.html

import (
	"time"

	"github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/fulltextsearch"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
)

// Struct that implements UserDataAcessor using given DAO, cache and full text search provider
type UserStorage struct {
	userDAO    database.UserDAO
	followsDAO database.FollowsDAO
	cache      cache.CacheProvider
	fts        fulltextsearch.UserSearcher
}

// Constructs UserStorage that uses given userDAO, followsDAO, CacheProvider and UserSearcher
func NewUserStorage(userDAO database.UserDAO, followsDAO database.FollowsDAO, cache cache.CacheProvider,
	fts fulltextsearch.UserSearcher) *UserStorage {
	return &UserStorage{
		userDAO:    userDAO,
		followsDAO: followsDAO,
		cache:      cache,
		fts:        fts,
	}
}

func (s *UserStorage) GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error) {
	var user *model.PublicUser

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", userID}, &user); !exists {
		var err error

		user, err = s.userDAO.GetPublicUserByID(userID)
		if err == errors.NoResultsError {
			return nil, err
		}
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"user", userID}, user)
	}

	err := s.collectPublicUserData(user, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return user, nil
}

func (s *UserStorage) GetUserByEmail(email string) (*model.User, error) {
	var user *model.User

	// Here we don't need any additional data other than what we fetched from database.
	// Dont use cache here, since this function will be used only for authentication users and we want
	// 100% real data for this.

	user, err := s.userDAO.GetUserByEmail(email)
	if err == errors.NoResultsError {
		return nil, err
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
	s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"user", insertedUser.ID, "followerCount"}, 0)
	s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"user", insertedUser.ID, "followeeCount"}, 0)

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
	followed, err := s.followsDAO.FollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	if followed {
		s.cache.IncrementWithFields(cache.Fields{"user", followeeID, "followerCount"})
		s.cache.IncrementWithFields(cache.Fields{"user", followerID, "followeeCount"})
	}
	s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"user", followeeID, "isFollowedBy", followerID}, true)

	return nil
}

func (s *UserStorage) UnfollowUser(followeeID, followerID int64) error {
	unfollowed, err := s.followsDAO.UnfollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	if unfollowed {
		s.cache.DecrementWithFields(cache.Fields{"user", followeeID, "followerCount"})
		s.cache.DecrementWithFields(cache.Fields{"user", followerID, "followeeCount"})
	}
	s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"user", followeeID, "isFollowedBy", followerID}, false)

	return nil
}

// TODO: this all could be done nicely in paralell
func (s *UserStorage) GetFollowers(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followersIDs := make([]int64, 0)

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", userID, "followersIDs"}, &followersIDs); !exists {
		var err error

		followersIDs, err = s.followsDAO.GetFollowersIDs(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
		s.cache.SetWithFields(cache.Fields{"user", userID, "followersIDs"}, followersIDs)
	}

	followers, err := s.getUsersByIDs(followersIDs, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return followers, nil
}

// TODO: this all could be done nicely in paralell
func (s *UserStorage) GetFollowees(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followeesIDs := make([]int64, 0)

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", userID, "followeesIDs"}, &followeesIDs); !exists {
		var err error

		followeesIDs, err = s.followsDAO.GetFolloweesIDs(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"user", userID, "followeesIDs"}, followeesIDs)
	}

	followees, err := s.getUsersByIDs(followeesIDs, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return followees, nil
}

func (s *UserStorage) GetUsersUsingQueryString(querystring string, requestingUserID int64) ([]*model.PublicUser, error) {
	usersIDs := make([]int64, 0)

	if exists, _ := s.cache.GetWithFields(cache.Fields{"users", "querystring", querystring}, &usersIDs); !exists {
		var err error

		usersIDs, err = s.fts.GetUsersIDs(querystring)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SetWithFields(cache.Fields{"users", "querystring", querystring}, usersIDs)
	}

	return s.getUsersByIDs(usersIDs, requestingUserID)
}

// Be careful - this is function does SIDE EFFECTS only
func (s *UserStorage) collectPublicUserData(user *model.PublicUser, requestingUserID int64) error {
	var followerCount int64
	var followeeCount int64
	var following bool

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", user.ID, "followerCount"}, &followerCount); !exists {
		var err error

		followerCount, err = s.followsDAO.GetFollowerCount(user.ID)
		if err != nil {
			return errors.UnexpectedError
		}

		s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"user", user.ID, "followerCount"}, followerCount)
	}

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", user.ID, "followeeCount"}, &followeeCount); !exists {
		var err error

		followeeCount, err = s.followsDAO.GetFolloweeCount(user.ID)
		if err != nil {
			return errors.UnexpectedError
		}

		s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"user", user.ID, "followeeCount"}, followeeCount)
	}

	if exists, _ := s.cache.GetWithFields(cache.Fields{"user", user.ID, "isFollowedBy", requestingUserID}, &following); !exists {
		var err error

		following, err = s.followsDAO.IsFollowing(requestingUserID, user.ID)
		if err != nil {
			return errors.UnexpectedError
		}

		s.cache.SetWithFieldsWithoutExpiration(cache.Fields{"user", user.ID, "isFollowedBy", requestingUserID}, following)
	}

	user.FollowerCount = followerCount
	user.FolloweeCount = followeeCount
	user.Following = following

	return nil
}

func (s *UserStorage) getUsersByIDs(usersIDs []int64, requestingUserID int64) ([]*model.PublicUser, error) {
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
		dbFollowers, err := s.userDAO.GetPublicUsersByIDs(usersIDs)
		if err != nil {
			return nil, err
		}
		users = append(users, dbFollowers...)
	}

	// fill users with missing data
	for _, user := range users {
		err := s.collectPublicUserData(user, requestingUserID)
		if err != nil {
			return nil, errors.UnexpectedError
		}
	}

	return users, nil
}
