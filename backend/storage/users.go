package storage

import (
	"time"

	"github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/async"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/model/errors"
	"github.com/VirrageS/chirp/backend/storage/cache"
	"github.com/VirrageS/chirp/backend/storage/database"
	"github.com/VirrageS/chirp/backend/storage/fulltextsearch"
)

// usersStorage is struct which implements userDataAccessor using given DAO, cache and full text search provider
type usersStorage struct {
	usersDAO   database.UsersDAO
	followsDAO database.FollowsDAO
	cache      cache.Accessor
	fts        fulltextsearch.UsersSearcher
}

// newUsersStorage constructs usersStorage that uses given usersDAO, followsDAO, Accessor and UsersSearcher
func newUsersStorage(usersDAO database.UsersDAO, followsDAO database.FollowsDAO, cache cache.Accessor, fts fulltextsearch.UsersSearcher) usersDataAccessor {
	return &usersStorage{
		usersDAO:   usersDAO,
		followsDAO: followsDAO,
		cache:      cache,
		fts:        fts,
	}
}

func (s *usersStorage) GetUserByID(userID, requestingUserID int64) (*model.PublicUser, error) {
	var (
		err  error
		user *model.PublicUser
	)

	key := cache.Key{"user", userID}
	if exists, _ := s.cache.GetSingle(key, user); !exists {
		user, err = s.usersDAO.GetPublicUserByID(userID)
		if err == errors.NoResultsError {
			return nil, err
		} else if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, user})
	}

	err = s.collectPublicUserData(user, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return user, nil
}

func (s *usersStorage) GetUserByEmail(email string) (*model.User, error) {
	// Here we don't need any additional data other than what we fetched from database.
	// Dont use cache here, since this function will be used only for authentication users and we want
	// 100% real data for this.

	user, err := s.usersDAO.GetUserByEmail(email)
	if err == errors.NoResultsError {
		return nil, err
	} else if err != nil {
		return nil, errors.UnexpectedError
	}

	return user, nil
}

func (s *usersStorage) InsertUser(newUserForm *model.NewUserForm) (*model.PublicUser, error) {
	insertedUser, err := s.usersDAO.InsertUser(newUserForm)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == database.UniqueConstraintViolationCode {
			return nil, errors.UserAlreadyExistsError
		}

		return nil, errors.UnexpectedError
	}

	s.cache.Set(cache.Entry{cache.Key{"user", insertedUser.ID}, insertedUser})

	return insertedUser, nil
}

func (s *usersStorage) UpdateUserLastLoginTime(userID int64, lastLoginTime *time.Time) error {
	err := s.usersDAO.UpdateUserLastLoginTime(userID, lastLoginTime)
	if err != nil {
		return errors.UnexpectedError
	}

	s.cache.Delete(cache.Key{"user", userID})

	return nil
}

func (s *usersStorage) FollowUser(followeeID, followerID int64) error {
	followed, err := s.followsDAO.FollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	if followed {
		s.cache.SAdd(cache.Key{"user", followeeID, "followers.ids"}, followerID)
		s.cache.SAdd(cache.Key{"user", followeeID, "followees.ids"}, followeeID)
	}

	return nil
}

func (s *usersStorage) UnfollowUser(followeeID, followerID int64) error {
	unfollowed, err := s.followsDAO.UnfollowUser(followeeID, followerID)
	if err != nil {
		return errors.UnexpectedError
	}

	if unfollowed {
		s.cache.SAdd(cache.Key{"user", followeeID, "followers.ids"}, followerID)
		s.cache.SAdd(cache.Key{"user", followeeID, "followees.ids"}, followeeID)
	}

	return nil
}

func (s *usersStorage) GetFollowers(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followersIDs := make([]int64, 0)

	key := cache.Key{"user", userID, "followers.ids"}
	if exists, _ := s.cache.GetSingle(key, &followersIDs); !exists {
		var err error

		followersIDs, err = s.followsDAO.GetFollowersIDs(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, followersIDs})
	}

	followers, err := s.getUsersByIDs(followersIDs, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return followers, nil
}

func (s *usersStorage) GetFollowees(userID, requestingUserID int64) ([]*model.PublicUser, error) {
	followeesIDs, err := s.GetFolloweesIDs(userID)
	if err != nil {
		return nil, err
	}

	followees, err := s.getUsersByIDs(followeesIDs, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return followees, nil
}

func (s *usersStorage) GetFolloweesIDs(userID int64) ([]int64, error) {
	followeesIDs := make([]int64, 0)

	key := cache.Key{"user", userID, "followees.ids"}
	if exists, _ := s.cache.SMembers(key, &followeesIDs); !exists {
		var err error

		followeesIDs, err = s.followsDAO.GetFolloweesIDs(userID)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.SAdd(key, followeesIDs)
	}

	return followeesIDs, nil
}

func (s *usersStorage) GetUsersUsingQueryString(querystring string, requestingUserID int64) ([]*model.PublicUser, error) {
	usersIDs := make([]int64, 0)

	key := cache.Key{"users", "querystring", querystring}
	if exists, _ := s.cache.GetSingle(key, &usersIDs); !exists {
		var err error

		usersIDs, err = s.fts.GetUsersIDs(querystring)
		if err != nil {
			return nil, errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, usersIDs})
	}

	return s.getUsersByIDs(usersIDs, requestingUserID)
}

func (s *usersStorage) collectPublicUsersData(users []*model.PublicUser, requestingUserID int64) error {
	pool := async.NewWorkerPool(func(task async.Task) *async.Result {
		err := s.collectPublicUserData(task.(*model.PublicUser), requestingUserID)
		return &async.Result{nil, err}
	})
	defer pool.Close()

	for _, user := range users {
		pool.PostTask(user)
	}

	for range users {
		if result := pool.GetResult(); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// Be careful - this is function does SIDE EFFECTS only
func (s *usersStorage) collectPublicUserData(user *model.PublicUser, requestingUserID int64) error {
	var (
		err           error
		followerCount int64
		followeeCount int64
		following     bool
	)

	key := cache.Key{"user", user.ID, "follower.count"}
	if exists, _ := s.cache.GetSingle(key, &followerCount); !exists {
		followerCount, err = s.followsDAO.GetFollowerCount(user.ID)
		if err != nil {
			return errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, followerCount})
	}

	key = cache.Key{"user", user.ID, "followee.count"}
	if exists, _ := s.cache.GetSingle(key, &followeeCount); !exists {
		followeeCount, err = s.followsDAO.GetFolloweeCount(user.ID)
		if err != nil {
			return errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, followeeCount})
	}

	key = cache.Key{"user", user.ID, "is.followed.by", requestingUserID}
	if exists, _ := s.cache.GetSingle(key, &following); !exists {
		following, err = s.followsDAO.IsFollowing(requestingUserID, user.ID)
		if err != nil {
			return errors.UnexpectedError
		}

		s.cache.Set(cache.Entry{key, following})
	}

	user.FollowerCount = followerCount
	user.FolloweeCount = followeeCount
	user.Following = following

	return nil
}

func (s *usersStorage) getUsersByIDs(usersIDs []int64, requestingUserID int64) ([]*model.PublicUser, error) {
	users := make([]*model.PublicUser, 0, len(usersIDs))

	pool := async.NewWorkerPool(func(task async.Task) *async.Result {
		var user *model.PublicUser
		id := task.(int64)

		key := cache.Key{"user", id}
		if exists, _ := s.cache.GetSingle(key, user); !exists {
			var err error

			user, err = s.usersDAO.GetPublicUserByID(id)
			if err != nil {
				return &async.Result{nil, err}
			}

			s.cache.Set(cache.Entry{key, user})
		}

		return &async.Result{user, nil}
	})
	defer pool.Close()

	for _, id := range usersIDs {
		pool.PostTask(id)
	}

	for range usersIDs {
		result := pool.GetResult()
		if result.Error != nil {
			return nil, errors.UnexpectedError
		}

		users = append(users, result.Value.(*model.PublicUser))
	}

	// Fill users with missing data like: followerCount, following etc...
	err := s.collectPublicUsersData(users, requestingUserID)
	if err != nil {
		return nil, errors.UnexpectedError
	}

	return users, nil
}
