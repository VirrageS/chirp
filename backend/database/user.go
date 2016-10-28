package database

import (
	"errors"
	"time"

	"github.com/VirrageS/chirp/backend/database/model"
)

var users = []model.User{
	{
		ID:        1,
		Name:      "george",
		Username:  "fisher",
		Email:     "corpsegrinder@cannibalcorpse.com",
		CreatedAt: time.Unix(0, 0),
	},
}

func GetUsers() ([]model.User, error) {
	return users, nil
}

func GetUser(userID int64) (model.User, error) {
	user, err := getUserWithId(userID)

	/* Emulate DB query fail? */
	if err != nil {
		return model.User{}, errors.New("")
	}

	return user, nil
}

func InsertUser(user model.User) (model.User, error) {
	if userAlreadyExists(user) {
		return model.User{}, errors.New("")
	}

	userID := insertUserToDatabase(user)
	user.ID = userID

	return user, nil
}

/* Functions that mock database queries */

func getUserWithId(userID int64) (model.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}

	return model.User{}, errors.New("")
}

func insertUserToDatabase(user model.User) int64 {
	userID := len(users) + 1
	user.ID = int64(userID)

	users = append(users, user)

	return int64(userID)
}

func userAlreadyExists(userToCheck model.User) bool {
	for _, user := range users {
		if user.Username == userToCheck.Username {
			return true
		}
	}
	return false
}
